package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tyler-smith/go-bip39"
)

// ============================================================================
// MODO 4 - DESCRAMBLER (DESEMBARALHADOR)
// Inspirado no BTCRecover - Descrambling de seeds
// Todas as palavras são conhecidas, mas a ordem foi perdida.
// O programa testa TODAS as permutações possíveis e valida o checksum.
// ============================================================================

// getSeedInputDescrambler coleta as palavras do usuário no modo Descrambler
func getSeedInputDescrambler(wordCount int, seedLanguage string) []string {
	clearScreen()
	showSimpleHeader()

	if uiLanguage == "pt" {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║           MODO 4 - DESCRAMBLER (DESEMBARALHADOR)                    ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("   Este modo é para quando você TEM TODAS as palavras da seed,")
		fmt.Println("   mas NÃO SABE a ordem correta.")
		fmt.Println()
		fmt.Println("   O programa testará TODAS as permutações possíveis e validará")
		fmt.Println("   o checksum BIP39/Electrum para encontrar a ordem correta.")
		fmt.Println()
		fmt.Println("   DIFERENÇA do Modo 3:")
		fmt.Println("   - Modo 3: Você define OK/NOK/? para cada posição (reduz permutações)")
		fmt.Println("   - Modo 4: Todas as palavras são permutadas automaticamente (mais rápido")
		fmt.Println("             de configurar, mas testa TODAS as combinações)")
		fmt.Println()
		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Printf("   Digite as %d palavras da sua seed (separadas por espaço):\n", wordCount)
		fmt.Println("   (A ordem NÃO importa - o programa testará todas as ordens)")
		fmt.Println()
	} else {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║           MODE 4 - DESCRAMBLER (UNSCRAMBLER)                        ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("   This mode is for when you HAVE ALL the seed words,")
		fmt.Println("   but DON'T KNOW the correct order.")
		fmt.Println()
		fmt.Println("   The program will test ALL possible permutations and validate")
		fmt.Println("   the BIP39/Electrum checksum to find the correct order.")
		fmt.Println()
		fmt.Println("   DIFFERENCE from Mode 3:")
		fmt.Println("   - Mode 3: You define OK/NOK/? for each position (reduces permutations)")
		fmt.Println("   - Mode 4: All words are permuted automatically (faster to configure,")
		fmt.Println("             but tests ALL combinations)")
		fmt.Println()
		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Printf("   Enter the %d words of your seed (separated by space):\n", wordCount)
		fmt.Println("   (Order does NOT matter - the program will test all orders)")
		fmt.Println()
	}

	wordlist := getWordList(seedLanguage)

	for {
		input := getUserInput("  > ")
		input = strings.TrimSpace(input)

		if input == "" {
			if uiLanguage == "pt" {
				fmt.Println("  [X] Entrada vazia! Digite as palavras.")
			} else {
				fmt.Println("  [X] Empty input! Enter the words.")
			}
			fmt.Println()
			continue
		}

		words := strings.Fields(input)

		if len(words) != wordCount {
			if uiLanguage == "pt" {
				fmt.Printf("  [X] Você digitou %d palavras, mas selecionou %d. Tente novamente.\n", len(words), wordCount)
			} else {
				fmt.Printf("  [X] You entered %d words, but selected %d. Try again.\n", len(words), wordCount)
			}
			fmt.Println()
			continue
		}

		// Validar que todas as palavras estão na wordlist
		allValid := true
		for i, word := range words {
			found := false
			for _, w := range wordlist {
				if w == word {
					found = true
					break
				}
			}
			if !found {
				allValid = false
				similar := findSimilarWords(word, wordlist, 5)
				if uiLanguage == "pt" {
					fmt.Printf("  [X] Palavra %d '%s' NÃO está na wordlist BIP39!\n", i+1, word)
					if len(similar) > 0 {
						fmt.Printf("      Você quis dizer: %s ?\n", strings.Join(similar, ", "))
					}
				} else {
					fmt.Printf("  [X] Word %d '%s' is NOT in the BIP39 wordlist!\n", i+1, word)
					if len(similar) > 0 {
						fmt.Printf("      Did you mean: %s ?\n", strings.Join(similar, ", "))
					}
				}
			}
		}

		if !allValid {
			fmt.Println()
			if uiLanguage == "pt" {
				fmt.Println("  Corrija as palavras e tente novamente.")
			} else {
				fmt.Println("  Fix the words and try again.")
			}
			fmt.Println()
			continue
		}

		// Verificar duplicatas
		wordMap := make(map[string]int)
		for _, w := range words {
			wordMap[w]++
		}
		hasDuplicates := false
		for w, count := range wordMap {
			if count > 1 {
				hasDuplicates = true
				if uiLanguage == "pt" {
					fmt.Printf("  [!] AVISO: Palavra '%s' aparece %d vezes.\n", w, count)
				} else {
					fmt.Printf("  [!] WARNING: Word '%s' appears %d times.\n", w, count)
				}
			}
		}
		if hasDuplicates {
			if uiLanguage == "pt" {
				fmt.Println("  Palavras duplicadas são permitidas (podem ser válidas).")
				fmt.Println("  Deseja continuar? (s/n)")
			} else {
				fmt.Println("  Duplicate words are allowed (may be valid).")
				fmt.Println("  Do you want to continue? (y/n)")
			}
			confirm := getUserInput("  > ")
			if confirm != "s" && confirm != "S" && confirm != "y" && confirm != "Y" {
				fmt.Println()
				continue
			}
		}

		// Mostrar resumo
		fmt.Println()
		if uiLanguage == "pt" {
			fmt.Println("  ┌─────────────────────────────────────────────────────────┐")
			fmt.Println("  │ PALAVRAS INSERIDAS:                                      │")
			fmt.Println("  └─────────────────────────────────────────────────────────┘")
		} else {
			fmt.Println("  ┌─────────────────────────────────────────────────────────┐")
			fmt.Println("  │ WORDS ENTERED:                                           │")
			fmt.Println("  └─────────────────────────────────────────────────────────┘")
		}
		for i, w := range words {
			fmt.Printf("   %2d. %s\n", i+1, w)
		}
		fmt.Println()

		return words
	}
}

// generateDescrambler executa a permutação completa no modo Descrambler
func (pg *PermutationGenerator) generateDescrambler() int {
	if uiLanguage == "pt" {
		fmt.Println("\nModo 4 - Descrambler: Testando TODAS as permutações...")
		fmt.Printf("   Total de palavras: %d\n", len(pg.words))
		fmt.Printf("   Permutações a testar: %d! (fatorial)\n", len(pg.words))
		fmt.Println("   Validando checksum para cada permutação...")
		fmt.Println()
	} else {
		fmt.Println("\nMode 4 - Descrambler: Testing ALL permutations...")
		fmt.Printf("   Total words: %d\n", len(pg.words))
		fmt.Printf("   Permutations to test: %d! (factorial)\n", len(pg.words))
		fmt.Println("   Validating checksum for each permutation...")
		fmt.Println()
	}

	// Configurar wordlist BIP39
	switch pg.seedLanguage {
	case "english":
		bip39.SetWordList(strings.Split(englishWordlist, "\n"))
	case "spanish":
		bip39.SetWordList(strings.Split(spanishWordlist, "\n"))
	case "french":
		bip39.SetWordList(strings.Split(frenchWordlist, "\n"))
	case "italian":
		bip39.SetWordList(strings.Split(italianWordlist, "\n"))
	case "portuguese":
		bip39.SetWordList(strings.Split(portugueseWordlist, "\n"))
	case "japanese":
		bip39.SetWordList(strings.Split(japaneseWordlist, "\n"))
	case "korean":
		bip39.SetWordList(strings.Split(koreanWordlist, "\n"))
	case "chinese_simplified":
		bip39.SetWordList(strings.Split(chineseSimplifiedWordlist, "\n"))
	case "chinese_traditional":
		bip39.SetWordList(strings.Split(chineseTraditionalWordlist, "\n"))
	}

	// Permutar todas as palavras
	pg.permuteDescrambler(pg.words, 0)

	// Salvar último arquivo
	pg.saveCurrentFile()

	// Renomear último arquivo para incluir FINAL
	if len(pg.filesCreated) > 0 {
		lastFile := pg.filesCreated[len(pg.filesCreated)-1]
		var prefix string
		if uiLanguage == "pt" {
			prefix = "Combinacoes_Possiveis"
		} else {
			prefix = "Possible_Combinations"
		}
		newBasename := fmt.Sprintf("%s_%02dFINAL_%s.xlsx", prefix, pg.currentFileNum, pg.dateStamp)
		newName := fmt.Sprintf("%s/%s", pg.outputDir, newBasename)
		os.Rename(lastFile, newName)
		pg.filesCreated[len(pg.filesCreated)-1] = newName
	}

	return pg.found
}

// permuteDescrambler gera todas as permutações recursivamente (Heap's algorithm)
func (pg *PermutationGenerator) permuteDescrambler(arr []string, start int) {
	if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
		return
	}

	if start == len(arr) {
		pg.checked++

		// Progresso a cada 100.000 permutações
		if pg.checked%100000 == 0 {
			elapsed := time.Since(pg.startTime).Seconds()
			rate := float64(pg.checked) / elapsed
			if uiLanguage == "pt" {
				fmt.Printf("\r   Verificadas: %d permutações | Válidas: %d | Taxa: %.0f/s   ", pg.checked, pg.found, rate)
			} else {
				fmt.Printf("\r   Checked: %d permutations | Valid: %d | Rate: %.0f/s   ", pg.checked, pg.found, rate)
			}
		}

		// Montar mnemônico e validar
		mnemonic := strings.Join(arr, " ")
		if isValidSeed(mnemonic) {
			pg.addSeed(mnemonic)
		}
		return
	}

	for i := start; i < len(arr); i++ {
		arr[start], arr[i] = arr[i], arr[start]
		pg.permuteDescrambler(arr, start+1)
		arr[start], arr[i] = arr[i], arr[start]
	}
}
