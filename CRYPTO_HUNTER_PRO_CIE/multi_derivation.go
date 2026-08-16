package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// MULTI DERIVATION PATH - TESTAR TODOS OS PATHS AUTOMATICAMENTE
// Inspirado no BTCRecover - Múltiplos Derivation Paths
// Permite ao usuário ativar TODOS os derivation paths de uma rede com um clique
// ============================================================================

// askTestAllDerivations pergunta ao usuário se deseja testar todos os paths
// para redes que possuem múltiplas derivações (BTC, LTC, BCH)
// Esta função é chamada ANTES de selectDerivations() e pode pré-habilitar tudo
func askTestAllDerivations(networks []NetworkGroup) []NetworkGroup {
	// Verificar se há alguma rede com múltiplas derivações habilitada
	hasMultiDerivation := false
	for _, ng := range networks {
		if ng.Enabled && len(ng.Derivations) > 1 {
			hasMultiDerivation = true
			break
		}
	}

	if !hasMultiDerivation {
		return networks
	}

	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("MODO DE DERIVACAO", "DERIVATION MODE"))
	fmt.Println("================================================================")
	fmt.Println()

	if uiLanguage == "pt" {
		fmt.Println("  Voce selecionou redes que possuem MULTIPLOS derivation paths.")
		fmt.Println("  (BTC, LTC, BCH possuem diferentes formatos de endereco)")
		fmt.Println()
		fmt.Println("  Deseja testar TODOS os derivation paths automaticamente?")
		fmt.Println("  (Isso aumenta a chance de encontrar fundos em qualquer formato)")
		fmt.Println()
		fmt.Println("  [1] Sim - Testar TODOS os paths de cada rede (recomendado)")
		fmt.Println("      (BTC: Legacy + SegWit + Native SegWit + Taproot)")
		fmt.Println("      (LTC: Legacy + SegWit + Native SegWit)")
		fmt.Println("      (BCH: CashAddr + Legacy)")
		fmt.Println()
		fmt.Println("  [2] Nao - Escolher manualmente quais paths testar")
		fmt.Println()
	} else {
		fmt.Println("  You selected networks that have MULTIPLE derivation paths.")
		fmt.Println("  (BTC, LTC, BCH have different address formats)")
		fmt.Println()
		fmt.Println("  Do you want to test ALL derivation paths automatically?")
		fmt.Println("  (This increases the chance of finding funds in any format)")
		fmt.Println()
		fmt.Println("  [1] Yes - Test ALL paths for each network (recommended)")
		fmt.Println("      (BTC: Legacy + SegWit + Native SegWit + Taproot)")
		fmt.Println("      (LTC: Legacy + SegWit + Native SegWit)")
		fmt.Println("      (BCH: CashAddr + Legacy)")
		fmt.Println()
		fmt.Println("  [2] No - Manually choose which paths to test")
		fmt.Println()
	}

	choice := getUserInput("  > ")
	choice = strings.TrimSpace(choice)

	if choice == "1" {
		// Habilitar TODAS as derivações de todas as redes habilitadas
		for i, ng := range networks {
			if ng.Enabled && len(ng.Derivations) > 1 {
				for j := range networks[i].Derivations {
					networks[i].Derivations[j].Enabled = true
				}
			}
		}

		fmt.Println()
		if uiLanguage == "pt" {
			fmt.Println("  [OK] Todos os derivation paths foram habilitados!")
			fmt.Println()
			fmt.Println("  Paths habilitados:")
		} else {
			fmt.Println("  [OK] All derivation paths have been enabled!")
			fmt.Println()
			fmt.Println("  Enabled paths:")
		}

		for _, ng := range networks {
			if ng.Enabled && len(ng.Derivations) > 1 {
				name := ng.NameEN
				if uiLanguage == "pt" {
					name = ng.NamePT
				}
				fmt.Printf("    %s:\n", name)
				for _, dp := range ng.Derivations {
					dpName := dp.NameEN
					if uiLanguage == "pt" {
						dpName = dp.NamePT
					}
					fmt.Printf("      - %s [%s]\n", dpName, dp.Path)
				}
			}
		}
		fmt.Println()

		// Marcar que não precisa perguntar derivações individuais
		// Retornar com flag especial (derivações já habilitadas)
		return networks
	}

	// Se escolheu "2" ou qualquer outra coisa, retorna sem alterar
	// e o fluxo normal de selectDerivations() será chamado
	return networks
}

// allDerivationsSelected verifica se todas as derivações de redes multi-path
// já foram habilitadas (para pular selectDerivations)
func allDerivationsSelected(networks []NetworkGroup) bool {
	for _, ng := range networks {
		if !ng.Enabled {
			continue
		}
		if len(ng.Derivations) <= 1 {
			continue
		}
		// Verificar se todas as derivações estão habilitadas
		allEnabled := true
		for _, dp := range ng.Derivations {
			if !dp.Enabled {
				allEnabled = false
				break
			}
		}
		if !allEnabled {
			return false
		}
	}
	return true
}
