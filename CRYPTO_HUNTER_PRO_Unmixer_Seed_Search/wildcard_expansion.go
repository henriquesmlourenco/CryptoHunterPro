package main

import "strings"

// expandWildcardsInWords expande todas as palavras com wildcards
// Retorna uma lista de listas, onde cada sublista contém as expansões de uma palavra
func expandWildcardsInWords(words []string, wordlist []string) [][]string {
	var expansions [][]string
	
	for _, word := range words {
		if word == "*" {
			// Asterisco puro = toda a wordlist
			expansions = append(expansions, wordlist)
		} else if strings.Contains(word, "*") {
			// Wildcard parcial = expandir
			matches := getMatchingWords(word, wordlist)
			if len(matches) == 0 {
				// Se não encontrou matches, usar a palavra original (não deveria acontecer)
				expansions = append(expansions, []string{word})
			} else {
				expansions = append(expansions, matches)
			}
		} else {
			// Palavra completa = apenas ela mesma
			expansions = append(expansions, []string{word})
		}
	}
	
	return expansions
}

// generateCartesianProduct gera o produto cartesiano de múltiplas listas
// Exemplo: [[a,b], [1,2]] -> [[a,1], [a,2], [b,1], [b,2]]
func generateCartesianProduct(lists [][]string) [][]string {
	if len(lists) == 0 {
		return [][]string{}
	}
	
	if len(lists) == 1 {
		result := make([][]string, len(lists[0]))
		for i, item := range lists[0] {
			result[i] = []string{item}
		}
		return result
	}
	
	// Recursivo: produto cartesiano do resto
	subProduct := generateCartesianProduct(lists[1:])
	var result [][]string
	
	for _, item := range lists[0] {
		for _, subList := range subProduct {
			combination := append([]string{item}, subList...)
			result = append(result, combination)
		}
	}
	
	return result
}

// countTotalCombinations calcula quantas combinações serão geradas
func countTotalCombinations(expansions [][]string) int {
	if len(expansions) == 0 {
		return 0
	}
	
	total := 1
	for _, expansion := range expansions {
		total *= len(expansion)
	}
	
	return total
}
