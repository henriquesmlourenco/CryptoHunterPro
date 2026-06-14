package main

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/tyler-smith/go-bip39"
)

// generateAdvancedWithWildcards gera permutações no modo avançado expandindo wildcards primeiro
func (pg *PermutationGenerator) generateAdvancedWithWildcards() int {
	if uiLanguage == "pt" {
		fmt.Println("\n🔄 Expandindo wildcards (Modo Avançado Parcial)...")
	} else {
		fmt.Println("\n🔄 Expanding wildcards (Advanced Partial Mode)...")
	}
	
	// Expandir wildcards em cada posição
	var expandedPositions [][]WordPosition
	
	for _, pos := range pg.config.Positions {
		if strings.Contains(pos.Word, "*") {
			// Expandir wildcard
			matches := getMatchingWords(pos.Word, pg.wordlist)
			
			if len(matches) == 0 {
				// Não deveria acontecer (já foi validado), mas por segurança
				matches = []string{pos.Word}
			}
			
			// Criar uma WordPosition para cada match
			var expansions []WordPosition
			for _, match := range matches {
				expansions = append(expansions, WordPosition{
					Position: pos.Position,
					Word:     match,
					Status:   pos.Status,
				})
			}
			expandedPositions = append(expandedPositions, expansions)
		} else {
			// Palavra completa, apenas uma opção
			expandedPositions = append(expandedPositions, []WordPosition{pos})
		}
	}
	
	// Calcular total de combinações
	totalCombinations := 1
	for _, expansions := range expandedPositions {
		totalCombinations *= len(expansions)
	}
	
	if uiLanguage == "pt" {
		fmt.Printf("   Total de combinações a testar: %d\n", totalCombinations)
		fmt.Println("\n🔄 Gerando e testando combinações...")
	} else {
		fmt.Printf("   Total combinations to test: %d\n", totalCombinations)
		fmt.Println("\n🔄 Generating and testing combinations...")
	}
	
	// Gerar produto cartesiano das posições expandidas
	positionCombinations := pg.generatePositionCartesian(expandedPositions, 0, []WordPosition{})
	
	// Para cada combinação de posições, aplicar lógica de permutação
	for idx, positions := range positionCombinations {
		if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
			break
		}
		
		// Mostrar progresso
		if idx > 0 && idx%100 == 0 {
			if uiLanguage == "pt" {
				fmt.Printf("\r   Progresso: %d/%d combinações | Válidas: %d   ",
					idx, len(positionCombinations), pg.found)
			} else {
				fmt.Printf("\r   Progress: %d/%d combinations | Valid: %d   ",
					idx, len(positionCombinations), pg.found)
			}
		}
		
		// Criar config temporária com as palavras expandidas
		tempConfig := &SeedConfig{
			WordCount:     pg.config.WordCount,
			Positions:     positions,
			FixedWords:    make(map[int]string),
			FloatingWords: []string{},
			FloatingPos:   []int{},
			MissingCount:  0,
			InputMode:     pg.config.InputMode,
		}
		
		// Separar fixas e flutuantes
		for _, pos := range positions {
			if pos.Status == "OK" {
				tempConfig.FixedWords[pos.Position] = pos.Word
			} else {
				tempConfig.FloatingWords = append(tempConfig.FloatingWords, pos.Word)
				tempConfig.FloatingPos = append(tempConfig.FloatingPos, pos.Position)
			}
		}
		
		// Permutar palavras flutuantes
		pg.permuteAdvancedCombination(tempConfig)
	}
	
	fmt.Println()
	
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

// generatePositionCartesian gera produto cartesiano de posições (recursivo)
func (pg *PermutationGenerator) generatePositionCartesian(expansions [][]WordPosition, index int, current []WordPosition) [][]WordPosition {
	if index == len(expansions) {
		// Criar cópia do current
		result := make([]WordPosition, len(current))
		copy(result, current)
		return [][]WordPosition{result}
	}
	
	var results [][]WordPosition
	for _, pos := range expansions[index] {
		newCurrent := append(current, pos)
		subResults := pg.generatePositionCartesian(expansions, index+1, newCurrent)
		results = append(results, subResults...)
	}
	
	return results
}

// permuteAdvancedCombination permuta palavras flutuantes de uma combinação específica
func (pg *PermutationGenerator) permuteAdvancedCombination(config *SeedConfig) {
	if len(config.FloatingWords) == 0 {
		// Todas fixas, apenas validar
		testWords := make([]string, config.WordCount)
		for pos, word := range config.FixedWords {
			testWords[pos-1] = word
		}
		
		mnemonic := strings.Join(testWords, " ")
		
		// Configurar wordlist correta
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
		
		if bip39.IsMnemonicValid(mnemonic) {
			pg.addSeed(mnemonic)
		}
	} else {
		// Permutar flutuantes
		pg.permuteAdvancedRecursive(config.FloatingWords, config.FloatingPos, config.FixedWords, config.WordCount, 0)
	}
}

// permuteAdvancedRecursive permuta palavras flutuantes recursivamente
func (pg *PermutationGenerator) permuteAdvancedRecursive(floatingWords []string, floatingPos []int, fixedMap map[int]string, wordCount int, start int) {
	if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
		return
	}
	
	if start == len(floatingWords) {
		pg.checked++
		
		// Construir seed completa
		testWords := make([]string, wordCount)
		for pos, word := range fixedMap {
			testWords[pos-1] = word
		}
		for i, pos := range floatingPos {
			testWords[pos-1] = floatingWords[i]
		}
		
		mnemonic := strings.Join(testWords, " ")
		
		// Configurar wordlist correta
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
		
		if bip39.IsMnemonicValid(mnemonic) {
			pg.addSeed(mnemonic)
		}
		return
	}
	
	for i := start; i < len(floatingWords); i++ {
		if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
			return
		}
		
		floatingWords[start], floatingWords[i] = floatingWords[i], floatingWords[start]
		pg.permuteAdvancedRecursive(floatingWords, floatingPos, fixedMap, wordCount, start+1)
		floatingWords[start], floatingWords[i] = floatingWords[i], floatingWords[start]
	}
}
