package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/pbkdf2"
)

// ============================================================================
// ELECTRUM SEED SUPPORT
// Electrum usa salt "electrum" ao inves de "mnemonic" (BIP-39)
// Derivation paths:
//   Standard: m/0/x (receive), m/1/x (change)
//   Segwit:   m/0'/0/x (receive), m/0'/1/x (change)
//   2FA:      m/1'/0/x (user key no multisig)
// Busca apenas BTC e LTC
// ============================================================================

// Tipos de seed Electrum
const (
	ElectrumSeedStandard = "standard" // prefix "01" - P2PKH (1...)
	ElectrumSeedSegwit   = "segwit"   // prefix "100" - P2WPKH (bc1q...)
	ElectrumSeed2FA      = "2fa"      // prefix "101" - P2SH multisig (3...)
)

// normalizeElectrumSeed normaliza a seed conforme Electrum:
// NFKD, lowercase, remove combining characters (accents), normalize whitespace
func normalizeElectrumSeed(seed string) string {
	seed = strings.ToLower(seed)
	// Remove combining characters (accents/diacritics)
	var result []rune
	for _, r := range seed {
		if !unicode.Is(unicode.Mn, r) {
			result = append(result, r)
		}
	}
	seed = string(result)
	// Normalize whitespace
	fields := strings.Fields(seed)
	return strings.Join(fields, " ")
}

// detectElectrumSeedType detecta o tipo de seed Electrum pelo prefixo HMAC
// Retorna: "standard", "segwit", "2fa", ou "" se nao for seed Electrum valida
func detectElectrumSeedType(seed string) string {
	normalized := normalizeElectrumSeed(seed)
	h := hmac.New(sha512.New, []byte("Seed version"))
	h.Write([]byte(normalized))
	hexHash := hex.EncodeToString(h.Sum(nil))

	// Ordem importa: verificar prefixos mais longos primeiro
	if strings.HasPrefix(hexHash, "101") {
		return ElectrumSeed2FA
	}
	if strings.HasPrefix(hexHash, "100") {
		return ElectrumSeedSegwit
	}
	if strings.HasPrefix(hexHash, "01") {
		return ElectrumSeedStandard
	}
	return ""
}

// electrumMnemonicToSeed converte mnemonic Electrum para seed bytes (64 bytes)
// Usa salt "electrum" + passphrase (diferente do BIP-39 que usa "mnemonic")
func electrumMnemonicToSeed(mnemonic string, passphrase string) []byte {
	normalized := normalizeElectrumSeed(mnemonic)
	salt := "electrum" + passphrase
	return pbkdf2.Key([]byte(normalized), []byte(salt), 2048, 64, sha512.New)
}

// ElectrumDerivedAddress armazena um endereco derivado de seed Electrum
type ElectrumDerivedAddress struct {
	Address     string
	PrivateKey  string
	Network     string // "btc" ou "ltc"
	AddressType string // "legacy", "native_segwit", "2fa_user_key"
	Path        string // ex: "m/0/0", "m/0'/0/0"
	Index       int
}

// DeriveElectrumAddresses gera enderecos BTC e LTC a partir de uma seed Electrum
func DeriveElectrumAddresses(seed string, passphrase string, seedType string, startIdx int, endIdx int) ([]ElectrumDerivedAddress, error) {
	var results []ElectrumDerivedAddress

	// Gerar master seed com salt "electrum"
	masterSeed := electrumMnemonicToSeed(seed, passphrase)

	// Gerar master key BIP32 (mesma funcao usada para BIP-39, so muda o seed)
	masterKey := SeedToMasterKey(masterSeed)

	switch seedType {
	case ElectrumSeedStandard:
		// Standard: m/0/x para receive (P2PKH)
		// Derivar m/0 primeiro (non-hardened)
		childM0, err := masterKey.DeriveChild(0)
		if err != nil {
			return nil, fmt.Errorf("erro ao derivar m/0: %v", err)
		}

		for idx := startIdx; idx <= endIdx; idx++ {
			// m/0/x
			childKey, err := childM0.DeriveChild(uint32(idx))
			if err != nil {
				continue
			}
			// BTC Legacy (1...)
			btcAddr, _ := deriveBTCLegacy(childKey.Key)
			results = append(results, ElectrumDerivedAddress{
				Address:     btcAddr,
				PrivateKey:  hex.EncodeToString(childKey.Key),
				Network:     "btc",
				AddressType: "legacy",
				Path:        fmt.Sprintf("m/0/%d", idx),
				Index:       idx,
			})
			// LTC Legacy (L...)
			ltcAddr, _ := deriveLTCLegacy(childKey.Key)
			results = append(results, ElectrumDerivedAddress{
				Address:     ltcAddr,
				PrivateKey:  hex.EncodeToString(childKey.Key),
				Network:     "ltc",
				AddressType: "legacy",
				Path:        fmt.Sprintf("m/0/%d", idx),
				Index:       idx,
			})
		}

	case ElectrumSeedSegwit:
		// Segwit: m/0'/0/x para receive (P2WPKH)
		// Derivar m/0' (hardened)
		childM0H, err := masterKey.DeriveChild(0 + 0x80000000)
		if err != nil {
			return nil, fmt.Errorf("erro ao derivar m/0': %v", err)
		}
		// Derivar m/0'/0 (non-hardened)
		childM0H0, err := childM0H.DeriveChild(0)
		if err != nil {
			return nil, fmt.Errorf("erro ao derivar m/0'/0: %v", err)
		}

		for idx := startIdx; idx <= endIdx; idx++ {
			// m/0'/0/x
			childKey, err := childM0H0.DeriveChild(uint32(idx))
			if err != nil {
				continue
			}
			// BTC Native Segwit (bc1q...)
			btcAddr, _ := deriveBTCNativeSegWit(childKey.Key)
			results = append(results, ElectrumDerivedAddress{
				Address:     btcAddr,
				PrivateKey:  hex.EncodeToString(childKey.Key),
				Network:     "btc",
				AddressType: "native_segwit",
				Path:        fmt.Sprintf("m/0'/0/%d", idx),
				Index:       idx,
			})
			// LTC Native Segwit (ltc1q...)
			ltcAddr, _ := deriveLTCNativeSegWit(childKey.Key)
			results = append(results, ElectrumDerivedAddress{
				Address:     ltcAddr,
				PrivateKey:  hex.EncodeToString(childKey.Key),
				Network:     "ltc",
				AddressType: "native_segwit",
				Path:        fmt.Sprintf("m/0'/0/%d", idx),
				Index:       idx,
			})
		}

	case ElectrumSeed2FA:
		// 2FA: m/1'/0/x - chave do usuario (1 de 3 no multisig)
		// Nota: sem o xpub do TrustedCoin, nao e possivel gerar o endereco P2SH correto
		// Geramos o endereco P2PKH da chave do usuario para referencia/busca
		childM1H, err := masterKey.DeriveChild(1 + 0x80000000)
		if err != nil {
			return nil, fmt.Errorf("erro ao derivar m/1': %v", err)
		}
		childM1H0, err := childM1H.DeriveChild(0)
		if err != nil {
			return nil, fmt.Errorf("erro ao derivar m/1'/0: %v", err)
		}

		for idx := startIdx; idx <= endIdx; idx++ {
			childKey, err := childM1H0.DeriveChild(uint32(idx))
			if err != nil {
				continue
			}
			// BTC - chave do usuario (P2PKH para referencia)
			btcAddr, _ := deriveBTCLegacy(childKey.Key)
			results = append(results, ElectrumDerivedAddress{
				Address:     btcAddr,
				PrivateKey:  hex.EncodeToString(childKey.Key),
				Network:     "btc",
				AddressType: "2fa_user_key",
				Path:        fmt.Sprintf("m/1'/0/%d", idx),
				Index:       idx,
			})
		}

		// Tambem tentar path standard m/0/x como fallback
		childM0, err := masterKey.DeriveChild(0)
		if err == nil {
			for idx := startIdx; idx <= endIdx; idx++ {
				childKey, err := childM0.DeriveChild(uint32(idx))
				if err != nil {
					continue
				}
				btcAddr, _ := deriveBTCLegacy(childKey.Key)
				results = append(results, ElectrumDerivedAddress{
					Address:     btcAddr,
					PrivateKey:  hex.EncodeToString(childKey.Key),
					Network:     "btc",
					AddressType: "legacy",
					Path:        fmt.Sprintf("m/0/%d", idx),
					Index:       idx,
				})
			}
		}
	}

	return results, nil
}

// runElectrumScan executa o scan para seeds Electrum (apenas BTC e LTC)
func runElectrumScan(config *ScanConfig) []ScanResult {
	var allResults []ScanResult

	totalWithBalance := 0
	totalWithHistory := 0

	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("ESCANEAMENTO ELECTRUM EM ANDAMENTO", "ELECTRUM SCANNING IN PROGRESS"))
	fmt.Println("================================================================")
	fmt.Println()

	for seedIdx, seed := range config.Seeds {
		fmt.Printf(t(
			"\n  [%d/%d] Escaneando seed Electrum: %s\n",
			"\n  [%d/%d] Scanning Electrum seed: %s\n"),
			seedIdx+1, len(config.Seeds), seed)

		// Detectar tipo da seed
		seedType := detectElectrumSeedType(seed)
		if seedType == "" {
			fmt.Println(t(
				"    [!] AVISO: Esta seed NAO foi reconhecida como seed Electrum valida.",
				"    [!] WARNING: This seed was NOT recognized as a valid Electrum seed."))
			fmt.Println(t(
				"    [!] Tentando como Standard (m/0/x) mesmo assim...",
				"    [!] Trying as Standard (m/0/x) anyway..."))
			seedType = ElectrumSeedStandard
		} else {
			seedTypeName := map[string]string{
				ElectrumSeedStandard: t("Standard (P2PKH - enderecos 1...)", "Standard (P2PKH - addresses 1...)"),
				ElectrumSeedSegwit:   t("Segwit (P2WPKH - enderecos bc1q...)", "Segwit (P2WPKH - addresses bc1q...)"),
				ElectrumSeed2FA:      t("2FA (Multisig - chave do usuario)", "2FA (Multisig - user key)"),
			}
			fmt.Printf(t(
				"    Tipo detectado: %s\n",
				"    Detected type: %s\n"), seedTypeName[seedType])
		}

		// Derivar enderecos
		addresses, err := DeriveElectrumAddresses(seed, config.Passphrase, seedType, config.StartIndex, config.EndIndex)
		if err != nil {
			fmt.Printf("    ERRO: %s\n", err.Error())
			continue
		}

		// Agrupar por rede para exibicao
		fmt.Printf(t(
			"\n    -> BTC + LTC [Electrum %s] indices %d-%d\n",
			"\n    -> BTC + LTC [Electrum %s] indices %d-%d\n"),
			seedType, config.StartIndex, config.EndIndex)

		// Consultar saldo de cada endereco
		for _, addr := range addresses {
			result := ScanResult{
				SeedPhrase:     seed,
				Network:        addr.Network,
				DerivationPath: addr.Path,
				Index:          addr.Index,
				Address:        addr.Address,
				PrivateKey:     addr.PrivateKey,
			}

			switch addr.Network {
			case "btc":
				btcResult, err := CheckBTCBalance(addr.Address)
				if err != nil {
					fmt.Printf("       [%d] %s (BTC %s) -> ERRO: %s\n", addr.Index, addr.Address, addr.AddressType, err.Error())
				} else if btcResult.HasBalance {
					result.NativeBalance = btcResult.Balance
					result.NativeSymbol = "BTC"
					result.HasBalance = true
					result.HasHistory = btcResult.HasHistory
					result.TxCount = btcResult.TxCount
					allResults = append(allResults, result)
					totalWithBalance++
					fmt.Printf("       [%d] %s (BTC %s) -> SALDO: %s BTC\n", addr.Index, addr.Address, addr.AddressType, btcResult.Balance)
				} else if btcResult.HasHistory {
					result.NativeBalance = btcResult.Balance
					result.NativeSymbol = "BTC"
					result.HasHistory = true
					result.TxCount = btcResult.TxCount
					allResults = append(allResults, result)
					totalWithHistory++
					fmt.Printf("       [%d] %s (BTC %s) -> %s\n", addr.Index, addr.Address, addr.AddressType, t("HISTORICO", "HISTORY"))
				} else {
					fmt.Printf("       [%d] %s (BTC %s) -> %s\n", addr.Index, addr.Address, addr.AddressType, t("vazio", "empty"))
				}

			case "ltc":
				ltcResult, err := CheckLTCBalance(addr.Address)
				if err != nil {
					fmt.Printf("       [%d] %s (LTC %s) -> ERRO: %s\n", addr.Index, addr.Address, addr.AddressType, err.Error())
				} else if ltcResult.HasBalance {
					result.NativeBalance = ltcResult.Balance
					result.NativeSymbol = "LTC"
					result.HasBalance = true
					result.HasHistory = ltcResult.HasHistory
					result.TxCount = ltcResult.TxCount
					allResults = append(allResults, result)
					totalWithBalance++
					fmt.Printf("       [%d] %s (LTC %s) -> SALDO: %s LTC\n", addr.Index, addr.Address, addr.AddressType, ltcResult.Balance)
				} else if ltcResult.HasHistory {
					result.NativeBalance = ltcResult.Balance
					result.NativeSymbol = "LTC"
					result.HasHistory = true
					result.TxCount = ltcResult.TxCount
					allResults = append(allResults, result)
					totalWithHistory++
					fmt.Printf("       [%d] %s (LTC %s) -> %s\n", addr.Index, addr.Address, addr.AddressType, t("HISTORICO", "HISTORY"))
				} else {
					fmt.Printf("       [%d] %s (LTC %s) -> %s\n", addr.Index, addr.Address, addr.AddressType, t("vazio", "empty"))
				}
			}
		}
	}

	// Resumo final
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Printf(t(
		"  RESULTADO: %d com saldo, %d com historico\n",
		"  RESULT: %d with balance, %d with history\n"),
		totalWithBalance, totalWithHistory)
	fmt.Println("================================================================")

	return allResults
}
