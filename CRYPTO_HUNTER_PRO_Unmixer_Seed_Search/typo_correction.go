package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// CORREÇÃO AUTOMÁTICA DE TYPOS
// Inspirado no BTCRecover - Typo Maps
// Quando o usuário digita uma palavra que não está na wordlist BIP39,
// o programa automaticamente sugere correções baseadas em:
// 1. Distância de Levenshtein (já existente)
// 2. Teclas adjacentes no teclado (typo map)
// 3. Letras trocadas (transposição)
// 4. Letras duplicadas acidentalmente
// ============================================================================

// keyboardAdjacentMap mapeia cada tecla para suas teclas adjacentes no teclado QWERTY
var keyboardAdjacentMap = map[byte]string{
	'q': "wa", 'w': "qeas", 'e': "wrds", 'r': "etdf", 't': "ryfg",
	'y': "tugh", 'u': "yijh", 'i': "uojk", 'o': "iplk", 'p': "ol",
	'a': "qwsz", 's': "wedxza", 'd': "erfcxs", 'f': "rtgvcd", 'g': "tyhbvf",
	'h': "yujnbg", 'j': "uikmnh", 'k': "iolmj", 'l': "opk",
	'z': "asx", 'x': "zsdc", 'c': "xdfv", 'v': "cfgb", 'b': "vghn",
	'n': "bhjm", 'm': "njk",
}

// generateTypoVariants gera variantes de uma palavra com erros de digitação comuns
func generateTypoVariants(word string) []string {
	variants := make(map[string]bool)

	// 1. Substituição por tecla adjacente
	for i := 0; i < len(word); i++ {
		if adjacent, ok := keyboardAdjacentMap[word[i]]; ok {
			for j := 0; j < len(adjacent); j++ {
				variant := word[:i] + string(adjacent[j]) + word[i+1:]
				variants[variant] = true
			}
		}
	}

	// 2. Transposição de letras adjacentes
	for i := 0; i < len(word)-1; i++ {
		variant := word[:i] + string(word[i+1]) + string(word[i]) + word[i+2:]
		variants[variant] = true
	}

	// 3. Remoção de letra duplicada
	for i := 0; i < len(word)-1; i++ {
		if word[i] == word[i+1] {
			variant := word[:i] + word[i+1:]
			variants[variant] = true
		}
	}

	// 4. Adição de letra duplicada
	for i := 0; i < len(word); i++ {
		variant := word[:i] + string(word[i]) + word[i:]
		variants[variant] = true
	}

	// 5. Remoção de uma letra (deleção)
	for i := 0; i < len(word); i++ {
		variant := word[:i] + word[i+1:]
		variants[variant] = true
	}

	result := make([]string, 0, len(variants))
	for v := range variants {
		result = append(result, v)
	}
	return result
}

// suggestTypoCorrections sugere correções para uma palavra inválida
// usando typo maps + Levenshtein + prefixo
// Retorna as melhores sugestões ordenadas por relevância
func suggestTypoCorrections(input string, wordlist []string, maxResults int) []string {
	if maxResults <= 0 {
		maxResults = 8
	}

	type suggestion struct {
		word   string
		score  int
		method string
	}

	seen := make(map[string]bool)
	var suggestions []suggestion

	// Método 1: Levenshtein (distância de edição <= 2 = muito provável)
	for _, word := range wordlist {
		dist := levenshteinDistance(input, word)
		if dist <= 2 && dist > 0 {
			if !seen[word] {
				seen[word] = true
				suggestions = append(suggestions, suggestion{word, dist, "levenshtein"})
			}
		}
	}

	// Método 2: Typo map (tecla adjacente)
	typoVariants := generateTypoVariants(input)
	wordSet := make(map[string]bool)
	for _, w := range wordlist {
		wordSet[w] = true
	}
	for _, variant := range typoVariants {
		if wordSet[variant] && !seen[variant] {
			seen[variant] = true
			suggestions = append(suggestions, suggestion{variant, 1, "typo_map"})
		}
	}

	// Método 3: Prefixo (primeiras 3-4 letras iguais)
	if len(input) >= 3 {
		prefix := input[:3]
		for _, word := range wordlist {
			if strings.HasPrefix(word, prefix) && !seen[word] {
				seen[word] = true
				dist := levenshteinDistance(input, word)
				suggestions = append(suggestions, suggestion{word, dist + 2, "prefix"})
			}
		}
	}

	// Ordenar por score (menor = melhor)
	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].score < suggestions[i].score {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	// Limitar resultados
	var result []string
	for i := 0; i < len(suggestions) && i < maxResults; i++ {
		result = append(result, suggestions[i].word)
	}

	return result
}

// validateWordWithTypoCorrection valida uma palavra e oferece correção automática
// Retorna: palavra corrigida (ou original se válida), e se foi aceita
func validateWordWithTypoCorrection(input string, wordlist []string) (string, bool) {
	// Verificar se a palavra já é válida
	for _, w := range wordlist {
		if w == input {
			return input, true
		}
	}

	// Palavra inválida - sugerir correções
	corrections := suggestTypoCorrections(input, wordlist, 8)

	if len(corrections) == 0 {
		if uiLanguage == "pt" {
			fmt.Printf("  [X] Palavra '%s' NAO esta na wordlist BIP39!\n", input)
			fmt.Println("      Nenhuma sugestao encontrada.")
		} else {
			fmt.Printf("  [X] Word '%s' is NOT in the BIP39 wordlist!\n", input)
			fmt.Println("      No suggestions found.")
		}
		return input, false
	}

	if uiLanguage == "pt" {
		fmt.Printf("\n  [!] Palavra '%s' NAO esta na wordlist BIP39!\n", input)
		fmt.Println("      Voce quis dizer:")
		fmt.Println()
	} else {
		fmt.Printf("\n  [!] Word '%s' is NOT in the BIP39 wordlist!\n", input)
		fmt.Println("      Did you mean:")
		fmt.Println()
	}

	for i, c := range corrections {
		fmt.Printf("      [%d] %s\n", i+1, c)
	}
	fmt.Println()

	if uiLanguage == "pt" {
		fmt.Printf("      [0] Manter '%s' como esta (ignorar)\n", input)
		fmt.Println()
		fmt.Println("      Digite o numero da correcao ou 0 para ignorar:")
	} else {
		fmt.Printf("      [0] Keep '%s' as is (ignore)\n", input)
		fmt.Println()
		fmt.Println("      Enter the correction number or 0 to ignore:")
	}

	choice := getUserInput("      > ")
	choice = strings.TrimSpace(choice)

	if choice == "0" || choice == "" {
		return input, false
	}

	// Converter escolha para índice
	idx := 0
	for _, c := range choice {
		if c >= '0' && c <= '9' {
			idx = idx*10 + int(c-'0')
		}
	}

	if idx >= 1 && idx <= len(corrections) {
		corrected := corrections[idx-1]
		if uiLanguage == "pt" {
			fmt.Printf("      [OK] Corrigido: '%s' -> '%s'\n", input, corrected)
		} else {
			fmt.Printf("      [OK] Corrected: '%s' -> '%s'\n", input, corrected)
		}
		return corrected, true
	}

	return input, false
}
