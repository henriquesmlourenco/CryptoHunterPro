package main

import (
	"fmt"
	"os"
)

func main() {
	// 0. Verificar licenca (NTP)
	if !checkLicense() {
		return
	}

	initScanner()

	// 1. Escolher idioma
	chooseUILanguage()

	// 2. Mostrar informações do programa
	showProgramInfo()

	// 3. Configurar APIs (antes do menu principal)
	showAPIConfigMenu()

	// Loop principal
	for {
		choice := showMainMenu()

		switch choice {
		case 1:
			// Digitar seed manualmente
			runManualSeedFlow()

		case 2:
			// Importar Excel do Unmixer Seed
			runImportFlow()

		case 0:
			fmt.Println()
			fmt.Println(t(
				"  Obrigado por usar o Crypto Hunter Pro!",
				"  Thank you for using Crypto Hunter Pro!"))
			fmt.Println()
			os.Exit(0)

		default:
			fmt.Println(t(
				"  [!] Opcao invalida.",
				"  [!] Invalid option."))
		}
	}
}

func runManualSeedFlow() {
	// Validação
	skipChecksum := chooseSeedValidation()

	// Input da seed
	seed := getSeedManual(skipChecksum)
	if seed == "" {
		return
	}

	// Passphrase
	passphrase := askPassphrase()

	// Seleção de redes
	networks := selectNetworks()

	// Multi-derivation: perguntar se quer testar todos os paths
	networks = askTestAllDerivations(networks)

	// Seleção de derivações (só pergunta se não habilitou tudo antes)
	if !allDerivationsSelected(networks) {
		networks = selectDerivations(networks)
	}

	// Range de índices
	startIdx, endIdx := selectIndexRange()

	// Configurar scan
	config := &ScanConfig{
		Seeds:      []string{seed},
		SeedSource: t("Manual", "Manual"),
		Networks:   networks,
		StartIndex: startIdx,
		EndIndex:   endIdx,
		Passphrase: passphrase,
	}

	// Resumo e confirmação
	if !showScanSummary(config) {
		fmt.Println(t("  Escaneamento cancelado.", "  Scan cancelled."))
		return
	}

	// Executar scan
	results := runScan(config)

	// Gerar relatório Excel
	if len(results) > 0 {
		fmt.Println()
		excelPath := generateExcelReport(results)
		if excelPath != "" {
			fmt.Printf(t(
				"\n  [OK] Relatorio Excel gerado: %s\n",
				"\n  [OK] Excel report generated: %s\n"), excelPath)
		}
	} else {
		fmt.Println()
		fmt.Println(t(
			"  Nenhum saldo ou historico encontrado.",
			"  No balance or history found."))
	}

	fmt.Println()
	getUserInput(t("  Pressione Enter para voltar ao menu...", "  Press Enter to return to menu..."))
}

func runImportFlow() {
	// Importar seeds do Excel
	seeds := importUnmixerSeedExcel()
	if len(seeds) == 0 {
		fmt.Println(t(
			"  [!] Nenhuma seed importada. Voltando ao menu.",
			"  [!] No seeds imported. Returning to menu."))
		return
	}

	// Validação (opcional para seeds importadas)
	fmt.Println()
	fmt.Println(t(
		"  Deseja validar o checksum BIP39 das seeds importadas?",
		"  Do you want to validate the BIP39 checksum of imported seeds?"))
	fmt.Println(t(
		"  [1] Sim - Remover seeds com checksum invalido",
		"  [1] Yes - Remove seeds with invalid checksum"))
	fmt.Println(t(
		"  [2] Nao - Manter todas as seeds",
		"  [2] No - Keep all seeds"))
	fmt.Println()

	choice := getUserInput("  > ")
	if choice == "1" {
		var validSeeds []string
		for _, s := range seeds {
			if ValidateSeedPhrase(s) {
				validSeeds = append(validSeeds, s)
			}
		}
		fmt.Printf(t(
			"  Seeds validas: %d de %d\n",
			"  Valid seeds: %d of %d\n"), len(validSeeds), len(seeds))
		seeds = validSeeds
		if len(seeds) == 0 {
			fmt.Println(t(
				"  [!] Nenhuma seed valida encontrada.",
				"  [!] No valid seeds found."))
			return
		}
	}

	// Passphrase
	passphrase := askPassphrase()

	// Seleção de redes
	networks := selectNetworks()

	// Multi-derivation: perguntar se quer testar todos os paths
	networks = askTestAllDerivations(networks)

	// Seleção de derivações (só pergunta se não habilitou tudo antes)
	if !allDerivationsSelected(networks) {
		networks = selectDerivations(networks)
	}

	// Range de índices
	startIdx, endIdx := selectIndexRange()

	// Configurar scan
	config := &ScanConfig{
		Seeds:      seeds,
		SeedSource: t("Excel Unmixer Seed", "Unmixer Seed Excel"),
		Networks:   networks,
		StartIndex: startIdx,
		EndIndex:   endIdx,
		Passphrase: passphrase,
	}

	// Resumo e confirmação
	if !showScanSummary(config) {
		fmt.Println(t("  Escaneamento cancelado.", "  Scan cancelled."))
		return
	}

	// Executar scan
	results := runScan(config)

	// Gerar relatório Excel
	if len(results) > 0 {
		fmt.Println()
		excelPath := generateExcelReport(results)
		if excelPath != "" {
			fmt.Printf(t(
				"\n  [OK] Relatorio Excel gerado: %s\n",
				"\n  [OK] Excel report generated: %s\n"), excelPath)
		}
	} else {
		fmt.Println()
		fmt.Println(t(
			"  Nenhum saldo ou historico encontrado.",
			"  No balance or history found."))
	}

	fmt.Println()
	getUserInput(t("  Pressione Enter para voltar ao menu...", "  Press Enter to return to menu..."))
}
