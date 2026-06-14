package main

import (
	"fmt"
	"time"
)

// ============================================================================
// MOTOR DE ESCANEAMENTO - MOSTRA CADA ENDERECO EM TEMPO REAL
// ============================================================================

func runScan(config *ScanConfig) []ScanResult {
	var allResults []ScanResult

	startTime := time.Now()

	totalSeeds := len(config.Seeds)
	totalChecked := 0
	totalWithBalance := 0
	totalWithHistory := 0

	// Mapa para evitar consultas duplicadas (mesma chave privada na mesma rede)
	// Isso acontece no BCH onde CashAddr e Legacy geram enderecos diferentes
	// mas a mesma carteira na blockchain
	checkedKeys := make(map[string]ScanResult)

	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("ESCANEAMENTO EM ANDAMENTO", "SCANNING IN PROGRESS"))
	fmt.Println("================================================================")
	fmt.Println()

	for seedIdx, seed := range config.Seeds {
		// Mostrar seed COMPLETA (sem truncar!)
		fmt.Printf(t(
			"\n  [%d/%d] Escaneando seed: %s\n",
			"\n  [%d/%d] Scanning seed: %s\n"),
			seedIdx+1, totalSeeds, seed)

		for _, ng := range config.Networks {
			if !ng.Enabled {
				continue
			}

			for _, dp := range ng.Derivations {
				if !dp.Enabled {
					continue
				}

				networkName := ng.NameEN
				if uiLanguage == "pt" {
					networkName = ng.NamePT
				}

				fmt.Printf(t(
					"\n    -> %s [%s] indices %d-%d\n",
					"\n    -> %s [%s] indices %d-%d\n"),
					networkName, dp.Path, config.StartIndex, config.EndIndex)

				derivationBalance := 0
				derivationHistory := 0
				derivationEmpty := 0

				for idx := config.StartIndex; idx <= config.EndIndex; idx++ {
					address, privKey, err := DeriveAddress(seed, config.Passphrase, dp, idx)
					if err != nil {
						fmt.Printf("       [%d] %s -> %s\n", idx, "ERRO DERIVACAO", err.Error())
						continue
					}

					totalChecked++

					result := ScanResult{
						SeedPhrase:     seed,
						Network:        ng.ID,
						DerivationPath: dp.Path,
						Index:          idx,
						Address:        address,
						PrivateKey:     privKey,
					}

				// Verificar se esta chave privada ja foi consultada nesta rede
				// (evita duplicacao BCH CashAddr/Legacy)
				checkKey := ng.ID + ":" + privKey
				if prevResult, exists := checkedKeys[checkKey]; exists {
					// Ja consultado - copiar resultado com endereco diferente
					result.NativeBalance = prevResult.NativeBalance
					result.NativeSymbol = prevResult.NativeSymbol
					result.HasBalance = prevResult.HasBalance
					result.HasHistory = prevResult.HasHistory
					result.TxCount = prevResult.TxCount
					result.LastTxDate = prevResult.LastTxDate
					result.Tokens = prevResult.Tokens
					if prevResult.HasBalance {
						fmt.Printf("       [%d] %s -> %s (mesmo saldo do formato anterior)\n",
							idx, address, t("SALDO: "+prevResult.NativeBalance+" "+prevResult.NativeSymbol+" [duplicado - mesmo endereco]",
								"BALANCE: "+prevResult.NativeBalance+" "+prevResult.NativeSymbol+" [duplicate - same address]"))
					} else if prevResult.HasHistory {
						fmt.Printf("       [%d] %s -> %s\n",
							idx, address, t("HISTORICO [duplicado - mesmo endereco]", "HISTORY [duplicate - same address]"))
					} else {
						fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
						derivationEmpty++
					}
					// NAO adicionar ao resultado para evitar duplicacao no Excel
					continue
				}

				switch ng.ID {
				case "evm":
						for _, evmNet := range ng.EVMNetworks {
							evmResult, err := CheckEVMBalance(address, evmNet)
							if err != nil {
								fmt.Printf("       [%d] %s (%s) -> ERRO: %s\n", idx, address, evmNet, err.Error())
								continue
							}

							if evmResult.HasBalance || evmResult.HasHistory {
								result.NativeBalance = evmResult.NativeBalance
								result.NativeSymbol = evmResult.NativeSymbol
								result.HasBalance = evmResult.HasBalance
								result.HasHistory = evmResult.HasHistory
								result.Tokens = evmResult.Tokens
								result.TxCount = evmResult.TxCount
								result.LastTxDate = evmResult.LastTxDate
								result.Network = evmNet

								allResults = append(allResults, result)
								if evmResult.HasBalance {
									totalWithBalance++
									derivationBalance++
									// Mostrar com destaque
									fmt.Printf("       [%d] %s (%s) -> SALDO: %s %s",
										idx, address, evmNet, evmResult.NativeBalance, evmResult.NativeSymbol)
									for _, tok := range evmResult.Tokens {
										fmt.Printf(" | %s: %s", tok.Symbol, tok.Balance)
									}
									fmt.Println()
								} else if evmResult.HasHistory {
									totalWithHistory++
									derivationHistory++
									fmt.Printf("       [%d] %s (%s) -> %s\n",
										idx, address, evmNet, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
								}
							} else {
								fmt.Printf("       [%d] %s (%s) -> %s\n",
									idx, address, evmNet, t("vazio", "empty"))
								derivationEmpty++
							}

							time.Sleep(250 * time.Millisecond)
						}

					case "btc":
						btcResult, err := CheckBTCBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if btcResult.HasBalance {
							result.NativeBalance = btcResult.Balance
							result.NativeSymbol = "BTC"
							result.HasBalance = true
							result.HasHistory = btcResult.HasHistory
							result.TxCount = btcResult.TxCount
							result.LastTxDate = btcResult.LastTxDate
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s BTC\n", idx, address, btcResult.Balance)
						} else if btcResult.HasHistory {
							result.NativeBalance = btcResult.Balance
							result.NativeSymbol = "BTC"
							result.HasHistory = true
							result.TxCount = btcResult.TxCount
							result.LastTxDate = btcResult.LastTxDate
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

				case "bch":
					bchResult, err := CheckBCHBalance(address)
					if err != nil {
						fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
					} else if bchResult.HasBalance {
						result.NativeBalance = bchResult.Balance
						result.NativeSymbol = "BCH"
						result.HasBalance = true
						result.HasHistory = bchResult.HasHistory
						result.TxCount = bchResult.TxCount
						result.LastTxDate = bchResult.LastTxDate
						allResults = append(allResults, result)
						checkedKeys[checkKey] = result
						totalWithBalance++
						derivationBalance++
						fmt.Printf("       [%d] %s -> SALDO: %s BCH\n", idx, address, bchResult.Balance)
					} else if bchResult.HasHistory {
						result.NativeBalance = bchResult.Balance
						result.NativeSymbol = "BCH"
						result.HasHistory = true
						result.TxCount = bchResult.TxCount
						result.LastTxDate = bchResult.LastTxDate
						allResults = append(allResults, result)
						checkedKeys[checkKey] = result
						totalWithHistory++
						derivationHistory++
						fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
					} else {
						fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
						derivationEmpty++
						checkedKeys[checkKey] = result
					}
					time.Sleep(500 * time.Millisecond)

					case "trx":
						trxResult, err := CheckTRXBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if trxResult.HasBalance {
							result.NativeBalance = trxResult.Balance
							result.NativeSymbol = "TRX"
							result.HasBalance = true
							result.HasHistory = trxResult.HasHistory
							result.Tokens = trxResult.Tokens
							result.TxCount = trxResult.TxCount
							result.LastTxDate = trxResult.LastTxDate
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s TRX", idx, address, trxResult.Balance)
							for _, tok := range trxResult.Tokens {
								fmt.Printf(" | %s: %s", tok.Symbol, tok.Balance)
							}
							fmt.Println()
						} else if trxResult.HasHistory {
							result.NativeBalance = trxResult.Balance
							result.NativeSymbol = "TRX"
							result.HasHistory = true
							result.TxCount = trxResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "sol":
						solResult, err := CheckSOLBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if solResult.HasBalance {
							result.NativeBalance = solResult.Balance
							result.NativeSymbol = "SOL"
							result.HasBalance = true
							result.HasHistory = solResult.HasHistory
							result.Tokens = solResult.Tokens
							result.TxCount = solResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s SOL", idx, address, solResult.Balance)
							for _, tok := range solResult.Tokens {
								fmt.Printf(" | %s: %s", tok.Symbol, tok.Balance)
							}
							fmt.Println()
						} else if solResult.HasHistory {
							result.NativeBalance = solResult.Balance
							result.NativeSymbol = "SOL"
							result.HasHistory = true
							result.TxCount = solResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "ltc":
						ltcResult, err := CheckLTCBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if ltcResult.HasBalance {
							result.NativeBalance = ltcResult.Balance
							result.NativeSymbol = "LTC"
							result.HasBalance = true
							result.HasHistory = ltcResult.HasHistory
							result.TxCount = ltcResult.TxCount
							result.LastTxDate = ltcResult.LastTxDate
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s LTC\n", idx, address, ltcResult.Balance)
						} else if ltcResult.HasHistory {
							result.NativeBalance = ltcResult.Balance
							result.NativeSymbol = "LTC"
							result.HasHistory = true
							result.TxCount = ltcResult.TxCount
							result.LastTxDate = ltcResult.LastTxDate
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "doge":
						dogeResult, err := CheckDOGEBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if dogeResult.HasBalance {
							result.NativeBalance = dogeResult.Balance
							result.NativeSymbol = "DOGE"
							result.HasBalance = true
							result.HasHistory = dogeResult.HasHistory
							result.TxCount = dogeResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s DOGE\n", idx, address, dogeResult.Balance)
						} else if dogeResult.HasHistory {
							result.NativeBalance = dogeResult.Balance
							result.NativeSymbol = "DOGE"
							result.HasHistory = true
							result.TxCount = dogeResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "ton":
						tonResult, err := CheckTONBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if tonResult.HasBalance {
							result.NativeBalance = tonResult.Balance
							result.NativeSymbol = "TON"
							result.HasBalance = true
							result.HasHistory = tonResult.HasHistory
							result.TxCount = tonResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s TON\n", idx, address, tonResult.Balance)
						} else if tonResult.HasHistory {
							result.NativeBalance = tonResult.Balance
							result.NativeSymbol = "TON"
							result.HasHistory = true
							result.TxCount = tonResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "zec":
						zecResult, err := CheckZECBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if zecResult.HasBalance {
							result.NativeBalance = zecResult.Balance
							result.NativeSymbol = "ZEC"
							result.HasBalance = true
							result.HasHistory = zecResult.HasHistory
							result.TxCount = zecResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s ZEC\n", idx, address, zecResult.Balance)
						} else if zecResult.HasHistory {
							result.NativeBalance = zecResult.Balance
							result.NativeSymbol = "ZEC"
							result.HasHistory = true
							result.TxCount = zecResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "xrp":
						xrpResult, err := CheckXRPBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if xrpResult.HasBalance {
							result.NativeBalance = xrpResult.Balance
							result.NativeSymbol = "XRP"
							result.HasBalance = true
							result.HasHistory = xrpResult.HasHistory
							result.TxCount = xrpResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s XRP\n", idx, address, xrpResult.Balance)
						} else if xrpResult.HasHistory {
							result.NativeBalance = xrpResult.Balance
							result.NativeSymbol = "XRP"
							result.HasHistory = true
							result.TxCount = xrpResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "xlm":
						xlmResult, err := CheckXLMBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if xlmResult.HasBalance {
							result.NativeBalance = xlmResult.Balance
							result.NativeSymbol = "XLM"
							result.HasBalance = true
							result.HasHistory = xlmResult.HasHistory
							result.Tokens = xlmResult.Tokens
							result.TxCount = xlmResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s XLM", idx, address, xlmResult.Balance)
							for _, tok := range xlmResult.Tokens {
								fmt.Printf(" | %s: %s", tok.Symbol, tok.Balance)
							}
							fmt.Println()
						} else if xlmResult.HasHistory {
							result.NativeBalance = xlmResult.Balance
							result.NativeSymbol = "XLM"
							result.HasHistory = true
							result.TxCount = xlmResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "algo":
						algoResult, err := CheckALGOBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if algoResult.HasBalance {
							result.NativeBalance = algoResult.Balance
							result.NativeSymbol = "ALGO"
							result.HasBalance = true
							result.HasHistory = algoResult.HasHistory
							result.Tokens = algoResult.Tokens
							result.TxCount = algoResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s ALGO", idx, address, algoResult.Balance)
							for _, tok := range algoResult.Tokens {
								fmt.Printf(" | %s: %s", tok.Symbol, tok.Balance)
							}
							fmt.Println()
						} else if algoResult.HasHistory {
							result.NativeBalance = algoResult.Balance
							result.NativeSymbol = "ALGO"
							result.HasHistory = true
							result.TxCount = algoResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "sui":
						suiResult, err := CheckSUIBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if suiResult.HasBalance {
							result.NativeBalance = suiResult.Balance
							result.NativeSymbol = "SUI"
							result.HasBalance = true
							result.HasHistory = suiResult.HasHistory
							result.Tokens = suiResult.Tokens
							result.TxCount = suiResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s SUI", idx, address, suiResult.Balance)
							for _, tok := range suiResult.Tokens {
								fmt.Printf(" | %s: %s", tok.Symbol, tok.Balance)
							}
							fmt.Println()
						} else if suiResult.HasHistory {
							result.NativeBalance = suiResult.Balance
							result.NativeSymbol = "SUI"
							result.HasHistory = true
							result.TxCount = suiResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)

					case "near":
						nearResult, err := CheckNEARBalance(address)
						if err != nil {
							fmt.Printf("       [%d] %s -> ERRO: %s\n", idx, address, err.Error())
						} else if nearResult.HasBalance {
							result.NativeBalance = nearResult.Balance
							result.NativeSymbol = "NEAR"
							result.HasBalance = true
							result.HasHistory = nearResult.HasHistory
							result.TxCount = nearResult.TxCount
							allResults = append(allResults, result)
							totalWithBalance++
							derivationBalance++
							fmt.Printf("       [%d] %s -> SALDO: %s NEAR\n", idx, address, nearResult.Balance)
						} else if nearResult.HasHistory {
							result.NativeBalance = nearResult.Balance
							result.NativeSymbol = "NEAR"
							result.HasHistory = true
							result.TxCount = nearResult.TxCount
							allResults = append(allResults, result)
							totalWithHistory++
							derivationHistory++
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("HISTORICO (sem saldo atual)", "HISTORY (no current balance)"))
						} else {
							fmt.Printf("       [%d] %s -> %s\n", idx, address, t("vazio", "empty"))
							derivationEmpty++
						}
						time.Sleep(500 * time.Millisecond)
					}
					}

					// Resumo da derivacao
				fmt.Printf(t(
					"       Resultado: %d com saldo, %d com historico, %d vazios\n",
					"       Result: %d with balance, %d with history, %d empty\n"),
					derivationBalance, derivationHistory, derivationEmpty)
			}
		}
	}

	elapsed := time.Since(startTime)

	// Resumo final
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("RESULTADO DO ESCANEAMENTO", "SCAN RESULTS"))
	fmt.Println("================================================================")
	fmt.Println()
	fmt.Printf("  %s: %s\n", t("Tempo total", "Total time"), elapsed.Round(time.Second))
	fmt.Printf("  %s: %d\n", t("Enderecos verificados", "Addresses checked"), totalChecked)
	fmt.Printf("  %s: %d\n", t("Enderecos com saldo", "Addresses with balance"), totalWithBalance)
	fmt.Printf("  %s: %d\n", t("Enderecos com historico", "Addresses with history"), totalWithHistory)
	fmt.Printf("  %s: %d\n", t("Total de resultados", "Total results"), len(allResults))

	return allResults
}
