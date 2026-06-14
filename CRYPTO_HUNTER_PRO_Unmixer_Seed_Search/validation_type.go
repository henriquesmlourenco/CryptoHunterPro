package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

// chooseValidationType exibe menu para escolher o tipo de validação de seed
func chooseValidationType() {
	for {
		if uiLanguage == "pt" {
			fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
			fmt.Println("║               TIPO DE VALIDAÇÃO DA SEED                              ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
			fmt.Println()
			fmt.Println("   Selecione o tipo de validação para as seeds geradas:")
			fmt.Println()
			fmt.Println("   [1] BIP39 (Padrão)")
			fmt.Println("       - Validação por checksum (padrão da indústria)")
			fmt.Println("       - Compatível com: Ledger, Trezor, MetaMask, Trust Wallet,")
			fmt.Println("         Exodus, Coinbase Wallet, e a maioria das carteiras")
			fmt.Println("       - Funciona em TODAS as redes: BTC, ETH, BSC, Polygon,")
			fmt.Println("         BCH, LTC, DOGE, SOL, ADA, TRX, XRP, e centenas mais")
			fmt.Println("       - Gera MENOS resultados (mais rápido e preciso)")
			fmt.Println()
			fmt.Println("   [2] HMAC-SHA512 (Electrum/Electron Cash nativo)")
			fmt.Println("       - Validação com HMAC-SHA512 (prefixo de versão)")
			fmt.Println("       - Para seeds criadas NATIVAMENTE dentro do Electrum (BTC)")
			fmt.Println("         ou Electron Cash (BCH)")
			fmt.Println("       - Usa as MESMAS palavras do BIP39, mas validação diferente")
			fmt.Println("       - IMPORTANTE: Se você criou a seed no Electrum/Electron Cash")
			fmt.Println("         mas ela é BIP39 (importada), use a opção [1]")
			fmt.Println("       - NÃO se aplica a ETH, BSC, Polygon ou outras redes EVM")
			fmt.Println()
			fmt.Println("   [3] Sem Validação (Força Bruta)")
			fmt.Println("       - Ignora checksum - gera TODAS as combinações possíveis")
			fmt.Println("       - Equivale ao \"Forçar interpretação BIP39\" do Electron Cash")
			fmt.Println("       - Útil quando não sabe qual carteira criou a seed")
			fmt.Println("       - AVISO: Gera ~16x MAIS resultados que a opção 1!")
			fmt.Println("       - Tempo de processamento MUITO maior!")
			fmt.Println()
			fmt.Print("   Escolha (1/2/3): ")
		} else {
			fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
			fmt.Println("║                 SEED VALIDATION TYPE                                 ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
			fmt.Println()
			fmt.Println("   Select the validation type for generated seeds:")
			fmt.Println()
			fmt.Println("   [1] BIP39 (Standard)")
			fmt.Println("       - Checksum validation (industry standard)")
			fmt.Println("       - Compatible with: Ledger, Trezor, MetaMask, Trust Wallet,")
			fmt.Println("         Exodus, Coinbase Wallet, and most wallets")
			fmt.Println("       - Works on ALL networks: BTC, ETH, BSC, Polygon,")
			fmt.Println("         BCH, LTC, DOGE, SOL, ADA, TRX, XRP, and hundreds more")
			fmt.Println("       - Generates FEWER results (faster and more precise)")
			fmt.Println()
			fmt.Println("   [2] HMAC-SHA512 (Native Electrum/Electron Cash)")
			fmt.Println("       - HMAC-SHA512 validation (version prefix)")
			fmt.Println("       - For seeds created NATIVELY inside Electrum (BTC)")
			fmt.Println("         or Electron Cash (BCH)")
			fmt.Println("       - Uses the SAME words as BIP39, but different validation")
			fmt.Println("       - IMPORTANT: If you created the seed in Electrum/Electron Cash")
			fmt.Println("         but it's BIP39 (imported), use option [1]")
			fmt.Println("       - Does NOT apply to ETH, BSC, Polygon or other EVM networks")
			fmt.Println()
			fmt.Println("   [3] No Validation (Brute Force)")
			fmt.Println("       - Ignores checksum - generates ALL possible combinations")
			fmt.Println("       - Equivalent to \"Force BIP39 interpretation\" in Electron Cash")
			fmt.Println("       - Useful when you don't know which wallet created the seed")
			fmt.Println("       - WARNING: Generates ~16x MORE results than option 1!")
			fmt.Println("       - Processing time MUCH longer!")
			fmt.Println()
			fmt.Print("   Choose (1/2/3): ")
		}

		choice := getUserInput("")
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			validationType = "bip39"
			if uiLanguage == "pt" {
				fmt.Println("\n   >> Validação BIP39 selecionada (Padrão)")
			} else {
				fmt.Println("\n   >> BIP39 validation selected (Standard)")
			}
			return
		case "2":
			validationType = "electrum"
			if uiLanguage == "pt" {
				fmt.Println("\n   >> Validação HMAC-SHA512 (Electrum/Electron Cash nativo) selecionada")
			} else {
				fmt.Println("\n   >> HMAC-SHA512 validation (Native Electrum/Electron Cash) selected")
			}
			return
		case "3":
			validationType = "none"
			if uiLanguage == "pt" {
				fmt.Println("\n   >> Sem Validação selecionada (Força Bruta)")
				fmt.Println("   AVISO: Isso gerará MUITOS mais resultados!")
			} else {
				fmt.Println("\n   >> No Validation selected (Brute Force)")
				fmt.Println("   WARNING: This will generate MANY more results!")
			}
			return
		default:
			if uiLanguage == "pt" {
				fmt.Println("\n   [X] Opção inválida! Escolha 1, 2 ou 3.")
			} else {
				fmt.Println("\n   [X] Invalid option! Choose 1, 2 or 3.")
			}
		}
	}
}

// isValidSeed valida uma seed de acordo com o tipo de validação selecionado
func isValidSeed(mnemonic string) bool {
	switch validationType {
	case "bip39":
		return bip39.IsMnemonicValid(mnemonic)
	case "electrum":
		return isValidElectrum(mnemonic)
	case "none":
		return true // Sem validação - aceita tudo
	default:
		return bip39.IsMnemonicValid(mnemonic)
	}
}

// isValidElectrum valida usando HMAC-SHA512 com prefixos de versão do Electrum
func isValidElectrum(mnemonic string) bool {
	// Electrum usa HMAC-SHA512("Seed version", mnemonic) e verifica o prefixo hex
	// Prefixos conhecidos do Electrum:
	// "01" = Standard (P2PKH)
	// "100" = SegWit (P2WPKH-P2SH)
	// "4b44" = 2FA

	hmacHash := hmac.New(sha512.New, []byte("Seed version"))
	hmacHash.Write([]byte(mnemonic))
	hexResult := hex.EncodeToString(hmacHash.Sum(nil))

	// Verificar prefixos válidos do Electrum
	electrumPrefixes := []string{
		"01",   // Standard (P2PKH) - Electrum BTC e Electron Cash BCH
		"100",  // SegWit (P2WPKH-P2SH)
		"4b44", // 2FA
	}

	for _, prefix := range electrumPrefixes {
		if strings.HasPrefix(hexResult, prefix) {
			return true
		}
	}

	return false
}
