package main

import (
	"math/big"
	"strings"
)

// calculateMode1Combinations calcula combinações para Modo 1 (ordem fixa, apenas wildcards)
func calculateMode1Combinations(words []string, wordlist []string) *big.Int {
	total := big.NewInt(1)
	
	for _, word := range words {
		if strings.Contains(word, "*") {
			// Contar quantas palavras correspondem a este wildcard
			matches := getMatchingWords(word, wordlist)
			matchCount := big.NewInt(int64(len(matches)))
			total.Mul(total, matchCount)
		}
		// Se não tem wildcard, é palavra fixa (multiplica por 1, não altera total)
	}
	
	return total
}
