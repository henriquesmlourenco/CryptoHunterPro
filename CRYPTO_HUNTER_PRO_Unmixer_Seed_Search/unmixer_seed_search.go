package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/tyler-smith/go-bip39"
	"github.com/xuri/excelize/v2"
)

// ============================================================================
// WORDLISTS EMBUTIDAS
// ============================================================================

//go:embed wordlists/english.txt
var englishWordlist string

//go:embed wordlists/spanish.txt
var spanishWordlist string

//go:embed wordlists/french.txt
var frenchWordlist string

//go:embed wordlists/italian.txt
var italianWordlist string

//go:embed wordlists/portuguese.txt
var portugueseWordlist string

//go:embed wordlists/japanese.txt
var japaneseWordlist string

//go:embed wordlists/korean.txt
var koreanWordlist string

//go:embed wordlists/chinese_simplified.txt
var chineseSimplifiedWordlist string

//go:embed wordlists/chinese_traditional.txt
var chineseTraditionalWordlist string


// ============================================================================
// ESTRUTURAS PARA MODO AVANÇADO
// ============================================================================

type WordPosition struct {
Position int
Word     string
Status   string // "OK", "NOK", "?"
}

type SeedConfig struct {
WordCount     int
Positions     []WordPosition
FixedWords    map[int]string // posição -> palavra (para OK)
FloatingWords []string        // palavras que flutuam (NOK e ?)
FloatingPos   []int           // posições disponíveis para flutuantes
MissingCount  int
InputMode     string // "simple", "advanced_partial", "advanced_complete"
}

// ============================================================================
// VARIÁVEIS GLOBAIS
// ============================================================================

var uiLanguage = "pt"
var validationType = "bip39" // "bip39", "electrum", "none"

// Limites de segurança e configurações
const (
	MAX_PERMUTATIONS_TO_GENERATE = 100000000 // 100 milhões
	BATCH_SIZE                   = 100000    // Processar em lotes de 100k
	ROWS_PER_FILE                = 500000    // 500.000 seeds por arquivo
)


// ============================================================================
// FUNÇÕES DE WILDCARD MATCHING
// ============================================================================

// matchWildcard verifica se uma palavra da wordlist corresponde ao padrão com wildcard
// Suporta: lugga*, *gage, in*ury, etc.
func matchWildcard(pattern, word string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == word
	}
	
	parts := strings.Split(pattern, "*")
	
	// Se tem apenas um *, pode estar no início, meio ou fim
	if len(parts) == 2 {
		prefix := parts[0]
		suffix := parts[1]
		
		// Verificar se a palavra tem o prefixo e sufixo corretos
		if len(prefix)+len(suffix) > len(word) {
			return false
		}
		
		return strings.HasPrefix(word, prefix) && strings.HasSuffix(word, suffix)
	}
	
	// Múltiplos wildcards (mais complexo)
	pos := 0
	for i, part := range parts {
		if part == "" {
			continue
		}
		
		idx := strings.Index(word[pos:], part)
		if idx == -1 {
			return false
		}
		
		// Se não é o primeiro part, pode ter qualquer coisa antes
		if i == 0 && idx != 0 {
			return false
		}
		
		pos += idx + len(part)
	}
	
	// Se o último part não é vazio, verificar se termina corretamente
	if parts[len(parts)-1] != "" && pos != len(word) {
		return false
	}
	
	return true
}

// getMatchingWords retorna todas as palavras da wordlist que correspondem ao padrão
func getMatchingWords(pattern string, wordlist []string) []string {
	var matches []string
	
	// Se não tem wildcard, retornar apenas a palavra se ela existir
	if !strings.Contains(pattern, "*") {
		for _, word := range wordlist {
			if word == pattern {
				return []string{pattern}
			}
		}
		return matches
	}
	
	// Se é apenas *, retornar toda a wordlist
	if pattern == "*" {
		return wordlist
	}
	
	// Buscar matches
	for _, word := range wordlist {
		if matchWildcard(pattern, word) {
			matches = append(matches, word)
		}
	}
	
	return matches
}

// ============================================================================
// FUNÇÕES DE INTERFACE
// ============================================================================

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func showWelcomeHeader() {
	fmt.Println()
	fmt.Println("  ██████╗██████╗ ██╗   ██╗██████╗ ████████╗ ██████╗ ")
	fmt.Println(" ██╔════╝██╔══██╗╚██╗ ██╔╝██╔══██╗╚══██╔══╝██╔═══██╗")
	fmt.Println(" ██║     ██████╔╝ ╚████╔╝ ██████╔╝   ██║   ██║   ██║")
	fmt.Println(" ██║     ██╔══██╗  ╚██╔╝  ██╔═══╝    ██║   ██║   ██║")
	fmt.Println(" ╚██████╗██║  ██║   ██║   ██║        ██║   ╚██████╔╝")
	fmt.Println("  ╚═════╝╚═╝  ╚═╝   ╚═╝   ╚═╝        ╚═╝    ╚═════╝ ")
	fmt.Println()
	fmt.Println(" ██╗  ██╗██╗   ██╗███╗   ██╗████████╗███████╗██████╗ ")
	fmt.Println(" ██║  ██║██║   ██║████╗  ██║╚══██╔══╝██╔════╝██╔══██╗")
	fmt.Println(" ███████║██║   ██║██╔██╗ ██║   ██║   █████╗  ██████╔╝")
	fmt.Println(" ██╔══██║██║   ██║██║╚██╗██║   ██║   ██╔══╝  ██╔══██╗")
	fmt.Println(" ██║  ██║╚██████╔╝██║ ╚████║   ██║   ███████╗██║  ██║")
	fmt.Println(" ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝")
	fmt.Println()
	fmt.Println("              ██████╗ ██████╗  ██████╗ ")
	fmt.Println("              ██╔══██╗██╔══██╗██╔═══██╗")
	fmt.Println("              ██████╔╝██████╔╝██║   ██║")
	fmt.Println("              ██╔═══╝ ██╔══██╗██║   ██║")
	fmt.Println("              ██║     ██║  ██║╚██████╔╝")
	fmt.Println("              ╚═╝     ╚═╝  ╚═╝ ╚═════╝ ")
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("           CRYPTO HUNTER PRO - UNMIXER SEED")
	fmt.Println("================================================================")
	fmt.Println()
	fmt.Println("  CRIADOR / CREATOR:                    CRIADOR / CREATOR:")
	fmt.Println("  Henrique Lourenco                     Alexandre Senra")
	fmt.Println("  linkedin.com/in/henriquelourenco      linkedin.com/in/alexandresenra")
	fmt.Println("  instagram.com/henrique.web3           instagram.com/alexandresenra_")
	fmt.Println()
	fmt.Println("  DOE / DONATE:")
	fmt.Println("  BTC: bc1qpq0cgvyxczetumdf87345zzk0zr0xz96ampmhs")
	fmt.Println("  ETH: henriquelourenco.eth")
	fmt.Println("  PIX: henriquesamuel@yahoo.com.br")
	fmt.Println()
	fmt.Println("  EN: Help us keep this project alive! Donate any amount.")
	fmt.Println("  EN: Free software, made with dedication. Support the creators!")
	fmt.Println()
	fmt.Println("  PT: Ajude-nos a manter este projeto vivo! Doe qualquer valor.")
	fmt.Println("  PT: Software livre, feito com dedicacao. Apoie os criadores!")
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println()

}

func showSimpleHeader() {
	if uiLanguage == "pt" {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("       CRYPTO HUNTER PRO - MÓDULO UNMIXER SEED")
		fmt.Println("================================================================")
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("       CRYPTO HUNTER PRO - UNMIXER SEED MODULE")
		fmt.Println("================================================================")
		fmt.Println()
	}
}

func chooseUILanguage() {
	for {
		fmt.Println("Choose language / Escolha o idioma:")
		fmt.Println("1. English")
		fmt.Println("2. Português")
		fmt.Println()

		choice := getUserInput("> ")
		if choice == "1" {
			uiLanguage = "en"
			return
		} else if choice == "2" {
			uiLanguage = "pt"
			return
		}
		fmt.Println()
		fmt.Println("  [X] Invalid option! Choose 1 or 2.")
		fmt.Println("  [X] Opcao invalida! Escolha 1 ou 2.")
		fmt.Println()
	}
}

func showUnmixerSeedSearchInfo() {
	clearScreen()
	showSimpleHeader()

	if uiLanguage == "pt" {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║         CRYPTO HUNTER PRO - MÓDULO UNMIXER SEED                      ║")
		fmt.Println("║              RECUPERADOR DE SEED PHRASES EMBARALHADAS                ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("O QUE ESTA FERRAMENTA FAZ:")
		fmt.Println()
		fmt.Println("Esta ferramenta foi desenvolvida para ajudar na recuperação de seed")
		fmt.Println("phrases que estão com as palavras EMBARALHADAS ou com PALAVRAS FALTANTES.")
		fmt.Println()
		fmt.Println(">>> FUNCIONALIDADES:")
		fmt.Println()
		fmt.Println("   Aceita seeds de 12, 15, 18, 21 ou 24 palavras")
		fmt.Println("   Suporta palavras em qualquer ordem (embaralhadas)")
		fmt.Println("   Permite palavras faltantes sem limite (use * para marcar)")
		fmt.Println("   Calcula QUANTAS permutações serão geradas ANTES de processar")
		fmt.Println("   Mostra quantos arquivos serão criados")
		fmt.Println("   Gera TODAS as combinações válidas possíveis")
		fmt.Println("   Exporta resultados em múltiplos arquivos Excel (.xlsx)")
		fmt.Println("   Divisão automática em arquivos de até 500.000 seeds")
		fmt.Println("   Nomenclatura sequencial profissional")
		fmt.Println("   Suporta TODOS os idiomas BIP39 (9 idiomas)")
		fmt.Println()
		fmt.Println("*** IDIOMAS SUPORTADOS:")
		fmt.Println()
		fmt.Println("   - English (Inglês)")
		fmt.Println("   - Español (Espanhol)")
		fmt.Println("   - Français (Francês)")
		fmt.Println("   - Italiano")
		fmt.Println("   - Português")
		fmt.Println("   - Japonês (Japanese)")
		fmt.Println("   - Coreano (Korean)")
		fmt.Println("   - Chinês Simplificado (Chinese Simplified)")
		fmt.Println("   - Chinês Tradicional (Chinese Traditional)")
		fmt.Println()
		fmt.Println("*** COMPATIBILIDADE - SEEDS BIP39 UNIVERSAIS:")
		fmt.Println()
		fmt.Println("Seeds BIP39 funcionam em TODAS as redes que suportam BIP39:")
		fmt.Println()
		fmt.Println("   REDE                  DERIVAÇÃO (PATH)         ENDEREÇO")
		fmt.Println("   -------------------- ----------------------- -------------")
		fmt.Println("   BTC (Legacy)          m/44'/0'/0'/0/x         1...")
		fmt.Println("   BTC (SegWit)          m/49'/0'/0'/0/x         3...")
		fmt.Println("   BTC (Native SegWit)   m/84'/0'/0'/0/x         bc1q...")
		fmt.Println("   BTC (Taproot)         m/86'/0'/0'/0/x         bc1p...")
		fmt.Println("   BCH                   m/44'/145'/0'/0/x       q... / bitcoincash:q...")
		fmt.Println("   ETH e redes EVM       m/44'/60'/0'/0/x        0x...")
		fmt.Println("   (BSC, Polygon, Avalanche, Arbitrum, Optimism, Base,")
		fmt.Println("    Fantom, Cronos, zkSync, Linea, Scroll, Mantle)")
		fmt.Println("   Litecoin (LTC)        m/44'/2'/0'/0/x         L...")
		fmt.Println("   Dogecoin (DOGE)       m/44'/3'/0'/0/x         D...")
		fmt.Println("   Tron (TRX)            m/44'/195'/0'/0/x       T...")
		fmt.Println("   Solana (SOL)          m/44'/501'/0'/0'        Base58")
		fmt.Println("   Cardano (ADA)         m/44'/1815'/0'/0/x      addr1...")
		fmt.Println("   XRP (Ripple)          m/44'/144'/0'/0/x       r...")
		fmt.Println("   Stellar (XLM)         m/44'/148'/0'/0/x       G...")
		fmt.Println("   Cosmos (ATOM)         m/44'/118'/0'/0/x       cosmos1...")
		fmt.Println("   Polkadot (DOT)        m/44'/354'/0'/0/x       1...")
		fmt.Println()
		fmt.Println("   NOTA IMPORTANTE: Este módulo gera apenas as seed phrases base.")
		fmt.Println("   A seed é universal e idêntica para todas as redes blockchain listadas")
		fmt.Println("   neste programa. O derivation path (caminho de derivação) pode ser")
		fmt.Println("   aplicado posteriormente pelo módulo principal do Crypto Hunter Pro -")
		fmt.Println("   Crypto Intelligence Engine (CIE), onde você poderá importar os arquivos")
		fmt.Println("   Excel gerados e selecionar os paths específicos para cada rede durante")
		fmt.Println("   a investigação forense.")
		fmt.Println()
		fmt.Println("═══════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Pressione Enter para continuar...")
	} else {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║         CRYPTO HUNTER PRO - UNMIXER SEED MODULE                      ║")
		fmt.Println("║              SHUFFLED SEED PHRASE RECOVERY TOOL                      ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("WHAT THIS TOOL DOES:")
		fmt.Println()
		fmt.Println("This tool was designed to help recover seed phrases that have")
		fmt.Println("SHUFFLED words or MISSING words.")
		fmt.Println()
		fmt.Println(">>> FEATURES:")
		fmt.Println()
		fmt.Println("   Accepts seeds with 12, 15, 18, 21, or 24 words")
		fmt.Println("   Supports words in any order (shuffled)")
		fmt.Println("   Allows missing words with no limit (use * to mark)")
		fmt.Println("   Calculates HOW MANY permutations will be generated BEFORE processing")
		fmt.Println("   Shows how many files will be created")
		fmt.Println("   Generates ALL possible valid combinations")
		fmt.Println("   Exports results to multiple Excel files (.xlsx)")
		fmt.Println("   Automatically splits into files up to 500,000 seeds")
		fmt.Println("   Sequential professional naming")
		fmt.Println("   Supports ALL BIP39 languages (9 languages)")
		fmt.Println()
		fmt.Println("*** SUPPORTED LANGUAGES:")
		fmt.Println()
		fmt.Println("   - English")
		fmt.Println("   - Español (Spanish)")
		fmt.Println("   - Français (French)")
		fmt.Println("   - Italiano (Italian)")
		fmt.Println("   - Português (Portuguese)")
		fmt.Println("   - Japanese")
		fmt.Println("   - Korean")
		fmt.Println("   - Chinese Simplified")
		fmt.Println("   - Chinese Traditional")
		fmt.Println()
		fmt.Println("*** COMPATIBILITY - UNIVERSAL BIP39 SEEDS:")
		fmt.Println()
		fmt.Println("BIP39 seeds work on ALL networks that support BIP39:")
		fmt.Println()
		fmt.Println("   NETWORK               DERIVATION PATH          ADDRESS")
		fmt.Println("   -------------------- ----------------------- -------------")
		fmt.Println("   BTC (Legacy)          m/44'/0'/0'/0/x         1...")
		fmt.Println("   BTC (SegWit)          m/49'/0'/0'/0/x         3...")
		fmt.Println("   BTC (Native SegWit)   m/84'/0'/0'/0/x         bc1q...")
		fmt.Println("   BTC (Taproot)         m/86'/0'/0'/0/x         bc1p...")
		fmt.Println("   BCH                   m/44'/145'/0'/0/x       q... / bitcoincash:q...")
		fmt.Println("   ETH and EVM networks  m/44'/60'/0'/0/x        0x...")
		fmt.Println("   (BSC, Polygon, Avalanche, Arbitrum, Optimism, Base,")
		fmt.Println("    Fantom, Cronos, zkSync, Linea, Scroll, Mantle)")
		fmt.Println("   Litecoin (LTC)        m/44'/2'/0'/0/x         L...")
		fmt.Println("   Dogecoin (DOGE)       m/44'/3'/0'/0/x         D...")
		fmt.Println("   Tron (TRX)            m/44'/195'/0'/0/x       T...")
		fmt.Println("   Solana (SOL)          m/44'/501'/0'/0'        Base58")
		fmt.Println("   Cardano (ADA)         m/44'/1815'/0'/0/x      addr1...")
		fmt.Println("   XRP (Ripple)          m/44'/144'/0'/0/x       r...")
		fmt.Println("   Stellar (XLM)         m/44'/148'/0'/0/x       G...")
		fmt.Println("   Cosmos (ATOM)         m/44'/118'/0'/0/x       cosmos1...")
		fmt.Println("   Polkadot (DOT)        m/44'/354'/0'/0/x       1...")
		fmt.Println()
		fmt.Println("   IMPORTANT NOTE: This module generates only the base seed phrases.")
		fmt.Println("   The seed is universal and identical for all blockchain networks listed")
		fmt.Println("   in this program. The derivation path can be applied later by the main")
		fmt.Println("   module of Crypto Hunter Pro - Crypto Intelligence Engine (CIE), where")
		fmt.Println("   you can import the generated Excel files and select specific paths for")
		fmt.Println("   each network during the forensic investigation.")
		fmt.Println()
		fmt.Println("═══════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Press Enter to continue...")
	}

	getUserInput("")
}

// Scanner global para leitura de input (evita problemas com múltiplos scanners)
var globalScanner *bufio.Scanner

func getGlobalScanner() *bufio.Scanner {
	if globalScanner == nil {
		globalScanner = bufio.NewScanner(os.Stdin)
		globalScanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	}
	return globalScanner
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := getGlobalScanner()
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// ============================================================================
// FUNÇÕES DE PROCESSAMENTO DE SEED
// ============================================================================

func chooseSeedLanguage() string {
	clearScreen()
	showSimpleHeader()

	if uiLanguage == "pt" {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    IDIOMA DA SEED PHRASE (BIP39)                    ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("Escolha o idioma das palavras da sua seed phrase:")
		fmt.Println()
		fmt.Println("1. English (Inglês)")
		fmt.Println("2. Español (Espanhol)")
		fmt.Println("3. Français (Francês)")
		fmt.Println("4. Italiano")
		fmt.Println("5. Português")
		fmt.Println("6. Japonês (Japanese)")
		fmt.Println("7. Coreano (Korean)")
		fmt.Println("8. Chinês Simplificado (Chinese Simplified)")
		fmt.Println("9. Chinês Tradicional (Chinese Traditional)")
		fmt.Println()
	} else {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    SEED PHRASE LANGUAGE (BIP39)                     ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("Choose your seed phrase words language:")
		fmt.Println()
		fmt.Println("1. English")
		fmt.Println("2. Español (Spanish)")
		fmt.Println("3. Français (French)")
		fmt.Println("4. Italiano (Italian)")
		fmt.Println("5. Português (Portuguese)")
		fmt.Println("6. Japanese")
		fmt.Println("7. Korean")
		fmt.Println("8. Chinese Simplified")
		fmt.Println("9. Chinese Traditional")
		fmt.Println()
	}

	for {
		choice := getUserInput("> ")

		languageMap := map[string]string{
			"1": "english",
			"2": "spanish",
			"3": "french",
			"4": "italian",
			"5": "portuguese",
			"6": "japanese",
			"7": "korean",
			"8": "chinese_simplified",
			"9": "chinese_traditional",
		}

		if lang, ok := languageMap[choice]; ok {
			return lang
		}
		if uiLanguage == "pt" {
			fmt.Println("\n  [X] Opcao invalida! Escolha de 1 a 9.")
		} else {
			fmt.Println("\n  [X] Invalid option! Choose 1 to 9.")
		}
		fmt.Println()
	}
}

func chooseSeedWordCount() int {
	clearScreen()
	showSimpleHeader()

	if uiLanguage == "pt" {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                  QUANTIDADE DE PALAVRAS DA SEED                      ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("Quantas palavras tem a sua seed phrase?")
		fmt.Println()
		fmt.Println("1. 12 palavras")
		fmt.Println("2. 15 palavras")
		fmt.Println("3. 18 palavras")
		fmt.Println("4. 21 palavras")
		fmt.Println("5. 24 palavras")
		fmt.Println()
	} else {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    SEED PHRASE WORD COUNT                            ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("How many words does your seed phrase have?")
		fmt.Println()
		fmt.Println("1. 12 words")
		fmt.Println("2. 15 words")
		fmt.Println("3. 18 words")
		fmt.Println("4. 21 words")
		fmt.Println("5. 24 words")
		fmt.Println()
	}

	for {
		choice := getUserInput("> ")

		wordCountMap := map[string]int{
			"1": 12,
			"2": 15,
			"3": 18,
			"4": 21,
			"5": 24,
		}

		if count, ok := wordCountMap[choice]; ok {
			return count
		}
		if uiLanguage == "pt" {
			fmt.Println("\n  [X] Opcao invalida! Escolha de 1 a 5.")
		} else {
			fmt.Println("\n  [X] Invalid option! Choose 1 to 5.")
		}
		fmt.Println()
	}
}

// Adicionar após a função chooseSeedWordCount e antes de getSeedInput

func chooseInputMode() string {
clearScreen()
showSimpleHeader()

if uiLanguage == "pt" {
fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
fmt.Println("║                    MODO DE INPUT DA SEED                             ║")
fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
fmt.Println()
fmt.Println("Escolha o modo de input:")
fmt.Println()
			fmt.Println("Modo 1 - Ordem Conhecida (com wildcards)")
			fmt.Println("   Use quando: SOUBER a ordem correta das palavras")
			fmt.Println("   Suporta wildcards: lugga*, in*ury, bo*, *")
			fmt.Println("   (use * quando não souber uma parte da palavra ou não souber a palavra inteira)")
			fmt.Println("   Cole todas as palavras NA ORDEM CORRETA")
			fmt.Println("   NÃO faz permutações (ordem fixa)")
fmt.Println()
			fmt.Println("Modo 2 - Ordem Desconhecida + Wildcards")
		fmt.Println("   Use quando: NÃO souber a ordem E tiver palavras INCOMPLETAS")
		fmt.Println("   Aceita wildcards: de*gn, bo*, in*ury, *")
		fmt.Println("   (use * quando não souber uma parte da palavra ou não souber a palavra inteira)")
		fmt.Println("   Digite palavra por palavra com ORDENAÇÃO")
			fmt.Println("   FAZ permutações (ordem desconhecida)")
fmt.Println()
			fmt.Println("Modo 3 - Ordem Desconhecida + SEM Wildcards")
	fmt.Println("   Use quando: NÃO souber a ordem mas tiver TODAS palavras COMPLETAS")
	fmt.Println("   NÃO aceita wildcards (apenas palavras completas BIP39)")
		fmt.Println("   Digite palavra por palavra com ORDENAÇÃO")
			fmt.Println("   FAZ permutações (ordem desconhecida)")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("O que são WILDCARDS (*)?")
	fmt.Println()
	fmt.Println("   O asterisco (*) substitui letras que você NÃO lembra da palavra.")
	fmt.Println()
	fmt.Println("   Exemplos de uso:")
	fmt.Println()
	fmt.Println("   - Palavra TOTALMENTE desconhecida:")
	fmt.Println("     Digite apenas: *")
	fmt.Println()
	fmt.Println("   - Sabe o INÍCIO da palavra:")
	fmt.Println("     book   -> bo*")
	fmt.Println("     design -> de*")
	fmt.Println("     injury -> in*")
	fmt.Println()
	fmt.Println("   - Sabe o MEIO da palavra:")
	fmt.Println("     design -> de*gn")
	fmt.Println("     injury -> in*ury")
	fmt.Println()
	fmt.Println("   - Sabe o FINAL da palavra:")
	fmt.Println("     luggage -> *age")
	fmt.Println("     ability -> *ity")
	fmt.Println()
	fmt.Println("   IMPORTANTE: O * pode substituir 1 ou mais letras ou a palavra")
	fmt.Println("               inteira caso você não saiba!")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("Observação sobre ORDENAÇÃO:")
	fmt.Println("   (quanto mais OK, MENOS possibilidades de seeds serão geradas,")
	fmt.Println("    facilitando para você fazer a busca em menos possibilidades)")
	fmt.Println()
	fmt.Println("   - OK  = Palavra FIXA nesta posição")
	fmt.Println("           -> Reduz MUITO a quantidade de possibilidades de seeds geradas")
	fmt.Println()
	fmt.Println("   - NOK = Palavra NÃO está nesta posição")
	fmt.Println("           -> Exclui esta posição (reduz um pouco as possibilidades")
	fmt.Println("              de seeds geradas)")
	fmt.Println()
	fmt.Println("   - ?   = Não sei se está nesta posição")
	fmt.Println("           -> Palavra pode estar em qualquer lugar (não reduz as")
	fmt.Println("              possibilidades de seeds geradas)")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()
			fmt.Println("DICA: Qual modo escolher? (Maior detalhamento, escrito acima)")
		fmt.Println("")
		fmt.Println("   - Modo 1: Não sei uma ou mais palavras e/ou tenho palavras")
		fmt.Println("             incompletas (ex: de*gn, bo*), mas SEI a ordem")
		fmt.Println("             correta de TODAS as palavras")
		fmt.Println("")
		fmt.Println("   - Modo 2: Não sei uma ou mais palavras e/ou tenho palavras")
		fmt.Println("             incompletas (ex: de*gn, bo*), e NÃO SEI a ordem")
		fmt.Println("             de duas ou mais palavras")
		fmt.Println("")
		fmt.Println("   - Modo 3: Tenho TODAS as palavras completas, mas NÃO SEI")
			fmt.Println("             a ordem de duas ou mais palavras")
		fmt.Println("")
		fmt.Println("   - Modo 4: DESCRAMBLER - Tenho TODAS as palavras completas,")
		fmt.Println("             NÃO SEI a ordem, e quero testar TODAS as")
		fmt.Println("             permutações automaticamente (sem definir OK/NOK/?)")
fmt.Println()
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("Qual Modo você deseja prosseguir: 1, 2, 3 ou 4:")
	fmt.Println()
} else {
fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
fmt.Println("║                    SEED INPUT MODE                                   ║")
fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
fmt.Println()
fmt.Println("Choose input mode:")
fmt.Println()
		fmt.Println("Mode 1. Simple Mode - Known Order (with wildcards)")
	fmt.Println("   Use when: You KNOW the correct word order")
	fmt.Println("   Supports wildcards: lugga*, in*ury, bo*, *")
	fmt.Println("   Paste all words IN CORRECT ORDER")
	fmt.Println("   NO permutations (fixed order)")
fmt.Println()
		fmt.Println("Mode 2. Advanced Partial - Unknown Order + Wildcards")
	fmt.Println("   Use when: You DON'T know order AND have INCOMPLETE words")
	fmt.Println("   Accepts wildcards: de*gn, bo*, in*ury, *")
	fmt.Println("   Enter word by word with ORDER")
	fmt.Println("   DOES permutations (unknown order)")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("Note about ORDER:")
		fmt.Println("   (more OK = LESS seed possibilities generated, easier search!):")
		fmt.Println("   - OK  = Word FIXED in this position")
		fmt.Println("           -> Reduces A LOT the amount of seed possibilities generated")
		fmt.Println("")
		fmt.Println("   - NOK = Word is NOT in this position")
		fmt.Println("           -> Excludes this position (reduces a bit the seed possibilities)")
		fmt.Println("")
		fmt.Println("   - ?   = Don't know if it's in this position")
		fmt.Println("           -> Word can be anywhere (doesn't reduce seed possibilities)")
	fmt.Println()
		fmt.Println("Mode 3. Advanced Complete - Unknown Order + NO Wildcards")
	fmt.Println("   Use when: You DON'T know order but have ALL COMPLETE words")
	fmt.Println("   NO wildcards (only complete BIP39 words)")
		fmt.Println("   Enter word by word with ORDER")
		fmt.Println("   DOES permutations (unknown order)")
fmt.Println()
fmt.Println()
	fmt.Println("Mode 4. Descrambler - Unknown Order (test ALL permutations)")
	fmt.Println("   Use when: You have ALL COMPLETE words but DON'T know order")
	fmt.Println("   NO wildcards, NO OK/NOK/? - tests ALL permutations automatically")
	fmt.Println("   Fastest to configure, but tests every possible combination")
fmt.Println()
fmt.Println("════════════════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println()
	fmt.Println("Which Mode do you want to proceed: 1, 2, 3 or 4:")
	fmt.Println()
}


for {
choice := getUserInput("> ")

if choice == "1" {
return "simple"
} else if choice == "2" {
return "advanced_partial"
} else if choice == "3" {
return "advanced_complete"
} else if choice == "4" {
return "descrambler"
}
if uiLanguage == "pt" {
fmt.Println("\n  [X] Opcao invalida! Escolha 1, 2, 3 ou 4.")
} else {
fmt.Println("\n  [X] Invalid option! Choose 1, 2, 3 or 4.")
}
fmt.Println()
}
}

func getSeedInputAdvancedPartial(wordCount int, seedLanguage string) *SeedConfig {
clearScreen()
	showSimpleHeader()

config := &SeedConfig{
WordCount:   wordCount,
Positions:   make([]WordPosition, wordCount),
FixedWords:  make(map[int]string),
FloatingWords: []string{},
FloatingPos: []int{},
MissingCount: 0,
InputMode:   "advanced_partial",
}

if uiLanguage == "pt" {
fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
fmt.Println("║                INPUT AVANÇADO - PALAVRA POR PALAVRA                  ║")
fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
fmt.Println()
fmt.Println("INSTRUÇÕES:")
fmt.Println()
	fmt.Println("Para cada posição, digite:")
	fmt.Println("  - A palavra COMPLETA ou com WILDCARD (ex: de*gn, bo*)")
	fmt.Println("  - Use * para palavra FALTANTE")
	fmt.Println("  - ORDENAÇÃO: OK, NOK ou ? (sem espaço)")
fmt.Println()
fmt.Println("ORDENAÇÃO:")
fmt.Println("  - OK  = Tenho CERTEZA que esta palavra está nesta posição")
fmt.Println("  - NOK = Tenho CERTEZA que esta palavra NÃO está nesta posição")
fmt.Println("  - ?   = NÃO SEI se está nesta posição")
fmt.Println()
	fmt.Println("EXEMPLOS:")
	fmt.Println("  Posição 1: abandon OK   (certeza que 'abandon' é a 1ª palavra)")
	fmt.Println("  Posição 2: ability NOK  (certeza que 'ability' NÃO é a 2ª)")
	fmt.Println("  Posição 3: de*gn ?      (wildcard: design, etc.)")
	fmt.Println("  Posição 4: bo* ?        (wildcard: board, boat, body, etc.)")
	fmt.Println("  Posição 5: * ?          (palavra faltante)")
fmt.Println()
fmt.Println("════════════════════════════════════════════════════════════════════════")
fmt.Println()
fmt.Println("QUANTO MAIS 'OK', MENOS PERMUTAÇÕES!")
fmt.Println()
fmt.Println("Exemplos de redução:")
fmt.Println("  - 12 palavras, todas NOK ou ? = ~30 milhões de seeds")
fmt.Println("  - 12 palavras, 3 com OK       = ~362 mil seeds (99.9% menos!)")
fmt.Println("  - 12 palavras, 6 com OK       = ~720 seeds (99.999% menos!)")
fmt.Println()
fmt.Println("════════════════════════════════════════════════════════════════════════")
fmt.Println()
} else {
fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
fmt.Println("║                ADVANCED INPUT - WORD BY WORD                         ║")
fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
fmt.Println()
fmt.Println("INSTRUCTIONS:")
fmt.Println()
	fmt.Println("For each position, enter:")
	fmt.Println("  - The COMPLETE word or with WILDCARD (ex: de*gn, bo*)")
	fmt.Println("  - Use * for MISSING word")
	fmt.Println("  - ORDER: OK, NOK or ? (no space)")
fmt.Println()
fmt.Println("ORDER:")
fmt.Println("  - OK  = I'm SURE this word is in this position")
fmt.Println("  - NOK = I'm SURE this word is NOT in this position")
fmt.Println("  - ?   = I DON'T KNOW if it's in this position")
fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  Position 1: abandon OK   (sure 'abandon' is 1st word)")
	fmt.Println("  Position 2: ability NOK  (sure 'ability' is NOT 2nd)")
	fmt.Println("  Position 3: de*gn ?      (wildcard: design, etc.)")
	fmt.Println("  Position 4: bo* ?        (wildcard: board, boat, body, etc.)")
	fmt.Println("  Position 5: * ?          (missing word)")
fmt.Println()
fmt.Println("════════════════════════════════════════════════════════════════════════")
fmt.Println()
fmt.Println("MORE 'OK' = LESS PERMUTATIONS!")
fmt.Println()
fmt.Println("Reduction examples:")
fmt.Println("  - 12 words, all NOK or ? = ~30 million seeds")
fmt.Println("  - 12 words, 3 with OK    = ~362 thousand seeds (99.9% less!)")
fmt.Println("  - 12 words, 6 with OK    = ~720 seeds (99.999% less!)")
fmt.Println()
fmt.Println("════════════════════════════════════════════════════════════════════════")
fmt.Println()
}

if uiLanguage == "pt" {
	fmt.Println("Pressione Enter para começar...")
} else {
	fmt.Println("Press Enter to start...")
}
getUserInput("")

// Input palavra por palavra
wordlist := getWordList(seedLanguage)

	for i := 0; i < wordCount; i++ {
	if uiLanguage == "pt" {
	fmt.Printf("\nPosição %d de %d:\n", i+1, wordCount)
	fmt.Print("Digite: palavra ORDENAÇÃO\n")
	fmt.Print("Exemplo: abandon OK\n")
fmt.Printf("> ")
		} else {
		fmt.Printf("\nPosition %d of %d:\n", i+1, wordCount)
		fmt.Print("Enter: word ORDER\n")
		fmt.Print("Example: abandon OK\n")
		fmt.Printf("> ")
		}
		
		input := getUserInput("")
		
		parts := strings.Fields(input)
		if len(parts) != 2 {
		if uiLanguage == "pt" {
		fmt.Println("Formato inválido! ORDENAÇÃO é OBRIGATÓRIA: palavra OK/NOK/?")
		} else {
		fmt.Println("Invalid format! ORDER is REQUIRED: word OK/NOK/?")
		}
		i--
		continue
		}
		
		word := strings.ToLower(parts[0])
		status := strings.ToUpper(parts[1])
	
	// Validar status
	if status != "OK" && status != "NOK" && status != "?" {
	if uiLanguage == "pt" {
	fmt.Println("ORDENAÇÃO inválida! Use: OK, NOK ou ?")
	} else {
	fmt.Println("Invalid ORDER! Use: OK, NOK or ?")
	}
	i--
	continue
	}

	// Validar palavra (se não for *)
if word != "*" {
	// Se tem wildcard, validar se pelo menos 1 palavra corresponde
	if strings.Contains(word, "*") {
		matches := getMatchingWords(word, wordlist)
		if len(matches) == 0 {
			if uiLanguage == "pt" {
				fmt.Printf("O padrão '%s' não corresponde a nenhuma palavra BIP39!\n", word)
				fmt.Println("   Exemplo: de*gn corresponde a 'design'")
			} else {
				fmt.Printf("The pattern '%s' doesn't match any BIP39 word!\n", word)
				fmt.Println("   Example: de*gn matches 'design'")
			}
			i--
			continue
		}
		// Wildcard válido, mostrar quantas palavras correspondem
		if uiLanguage == "pt" {
			fmt.Printf("Padrão '%s' corresponde a %d palavra(s)\n", word, len(matches))
		} else {
			fmt.Printf("Pattern '%s' matches %d word(s)\n", word, len(matches))
		}
	} else {
		// Palavra completa, validar normalmente
		found := false
		for _, w := range wordlist {
			if w == word {
				found = true
				break
			}
		}
		if !found {
			if uiLanguage == "pt" {
				fmt.Printf("A palavra '%s' não está na lista BIP39!\n", word)
			} else {
				fmt.Printf("The word '%s' is not in the BIP39 wordlist!\n", word)
			}
			i--
			continue
		}
	}
} else {
config.MissingCount++
}

// Armazenar
config.Positions[i] = WordPosition{
Position: i + 1,
Word:     word,
Status:   status,
}

		// Organizar por tipo
		if status == "OK" {
		config.FixedWords[i+1] = word
		} else {
		if word != "*" {
		config.FloatingWords = append(config.FloatingWords, word)
		}
		config.FloatingPos = append(config.FloatingPos, i+1)
		}
	}
	
	// Mostrar resumo completo
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	if uiLanguage == "pt" {
		fmt.Println("\n>>> RESUMO DAS PALAVRAS DIGITADAS:\n")
	} else {
		fmt.Println("\n>>> SUMMARY OF ENTERED WORDS:\n")
	}
	
	for i, pos := range config.Positions {
		var statusIcon string
		switch pos.Status {
			case "OK":
				statusIcon = "[OK] "
			case "NOK":
				statusIcon = "[NOK]"
			case "?":
				statusIcon = "[?]  "
		}
		if uiLanguage == "pt" {
			fmt.Printf("   %s Posição %2d: %-15s [%s]\n", statusIcon, i+1, pos.Word, pos.Status)
		} else {
			fmt.Printf("   %s Position %2d: %-15s [%s]\n", statusIcon, i+1, pos.Word, pos.Status)
		}
	}
	
	fmt.Println()
	if uiLanguage == "pt" {
		fmt.Printf("   Palavras FIXAS (OK): %d\n", len(config.FixedWords))
		fmt.Printf("   Palavras FLUTUANTES (NOK/?): %d\n", len(config.FloatingWords))
		if config.MissingCount > 0 {
			fmt.Printf("   Palavras FALTANTES (*): %d\n", config.MissingCount)
		}
	} else {
		fmt.Printf("   FIXED words (OK): %d\n", len(config.FixedWords))
		fmt.Printf("   FLOATING words (NOK/?): %d\n", len(config.FloatingWords))
		if config.MissingCount > 0 {
			fmt.Printf("   MISSING words (*): %d\n", config.MissingCount)
		}
	}
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()
	
	if uiLanguage == "pt" {
		fmt.Print("Pressione Enter para continuar...")
	} else {
		fmt.Print("Press Enter to continue...")
	}
	getUserInput("")
	
	return config
}

func getSeedInputAdvancedComplete(wordCount int, seedLanguage string) *SeedConfig {
clearScreen()
	showSimpleHeader()

config := &SeedConfig{
WordCount:   wordCount,
Positions:   make([]WordPosition, wordCount),
FixedWords:  make(map[int]string),
FloatingWords: []string{},
FloatingPos: []int{},
MissingCount: 0,
InputMode:   "advanced_complete",
}

if uiLanguage == "pt" {
fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
fmt.Println("║          INPUT AVANÇADO - PALAVRAS COMPLETAS (SEM WILDCARDS)         ║")
fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
fmt.Println()
fmt.Println("INSTRUÇÕES:")
fmt.Println()
fmt.Println("Para cada posição, digite:")
fmt.Println("  - A palavra COMPLETA da BIP39 (sem wildcards)")
fmt.Println("  - ORDENAÇÃO: OK, NOK ou ? (sem espaço)")
fmt.Println()
fmt.Println("ORDENAÇÃO:")
fmt.Println("  - OK  = Tenho CERTEZA que esta palavra está nesta posição")
fmt.Println("  - NOK = Tenho CERTEZA que esta palavra NÃO está nesta posição")
fmt.Println("  - ?   = NÃO SEI se está nesta posição")
fmt.Println()
fmt.Println("ATENÇÃO: Wildcards (*) NÃO são permitidos neste modo!")
fmt.Println("   Use Modo 2 se precisar de wildcards.")
fmt.Println()
fmt.Println("════════════════════════════════════════════════════════════════════════")
fmt.Println()
} else {
fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
fmt.Println("║          ADVANCED INPUT - COMPLETE WORDS (NO WILDCARDS)              ║")
fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
fmt.Println()
fmt.Println("INSTRUCTIONS:")
fmt.Println()
fmt.Println("For each position, enter:")
fmt.Println("  - The COMPLETE BIP39 word (no wildcards)")
fmt.Println("  - ORDER: OK, NOK or ? (no space)")
fmt.Println()
fmt.Println("ORDER:")
fmt.Println("  - OK  = I'm SURE this word is in this position")
fmt.Println("  - NOK = I'm SURE this word is NOT in this position")
fmt.Println("  - ?   = I DON'T KNOW if it's in this position")
fmt.Println()
fmt.Println("WARNING: Wildcards (*) are NOT allowed in this mode!")
fmt.Println("   Use Mode 2 if you need wildcards.")
fmt.Println()
fmt.Println("════════════════════════════════════════════════════════════════════════")
fmt.Println()
}

if uiLanguage == "pt" {
	fmt.Println("Pressione Enter para começar...")
} else {
	fmt.Println("Press Enter to start...")
}
getUserInput("")

// Input palavra por palavra
wordlist := getWordList(seedLanguage)

	for i := 0; i < wordCount; i++ {
	if uiLanguage == "pt" {
	fmt.Printf("\nPosição %d de %d:\n", i+1, wordCount)
	fmt.Print("Digite: palavra ORDENAÇÃO\n")
	fmt.Print("Exemplo: abandon OK\n")
	fmt.Printf("> ")
	} else {
	fmt.Printf("\nPosition %d of %d:\n", i+1, wordCount)
	fmt.Print("Enter: word ORDER\n")
	fmt.Print("Example: abandon OK\n")
	fmt.Printf("> ")
	}
	
	input := getUserInput("")
	
	parts := strings.Fields(input)
	if len(parts) != 2 {
	if uiLanguage == "pt" {
	fmt.Println("Formato inválido! ORDENAÇÃO é OBRIGATÓRIA: palavra OK/NOK/?")
	} else {
	fmt.Println("Invalid format! ORDER is REQUIRED: word OK/NOK/?")
	}
	i--
	continue
	}
	
	word := strings.ToLower(parts[0])
	status := strings.ToUpper(parts[1])
	
	// Validar status
	if status != "OK" && status != "NOK" && status != "?" {
	if uiLanguage == "pt" {
	fmt.Println("ORDENAÇÃO inválida! Use: OK, NOK ou ?")
	} else {
	fmt.Println("Invalid ORDER! Use: OK, NOK or ?")
	}
	i--
	continue
	}

	// VALIDAÇÃO: NÃO permitir wildcards neste modo
	if strings.Contains(word, "*") {
	if uiLanguage == "pt" {
	fmt.Println("Wildcards (*) NÃO são permitidos neste modo!")
	fmt.Println("   Use Modo 2 (Palavras Parciais) se precisar de wildcards.")
	} else {
	fmt.Println("Wildcards (*) are NOT allowed in this mode!")
	fmt.Println("   Use Mode 2 (Partial Words) if you need wildcards.")
	}
	i--
	continue
	}

	// Validar palavra na BIP39
found := false
for _, w := range wordlist {
if w == word {
found = true
break
}
}
		if !found {
		if uiLanguage == "pt" {
		fmt.Printf("A palavra '%s' não está na lista BIP39!\n", word)
		} else {
		fmt.Printf("The word '%s' is not in the BIP39 wordlist!\n", word)
		}
		i--
		continue
		}

// Armazenar
config.Positions[i] = WordPosition{
Position: i + 1,
Word:     word,
Status:   status,
}

			// Organizar por tipo
			if status == "OK" {
			config.FixedWords[i+1] = word
			} else {
			config.FloatingWords = append(config.FloatingWords, word)
			config.FloatingPos = append(config.FloatingPos, i+1)
			}
	}
	
	// Mostrar resumo completo
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	if uiLanguage == "pt" {
		fmt.Println("\n>>> RESUMO DAS PALAVRAS DIGITADAS:\n")
	} else {
		fmt.Println("\n>>> SUMMARY OF ENTERED WORDS:\n")
	}
	
	for i, pos := range config.Positions {
		var statusIcon string
		switch pos.Status {
			case "OK":
				statusIcon = "[OK] "
			case "NOK":
				statusIcon = "[NOK]"
			case "?":
				statusIcon = "[?]  "
			}
			if uiLanguage == "pt" {
			fmt.Printf("   %s Posição %2d: %-15s [%s]\n", statusIcon, i+1, pos.Word, pos.Status)
		} else {
			fmt.Printf("   %s Position %2d: %-15s [%s]\n", statusIcon, i+1, pos.Word, pos.Status)
		}
		}
	
	fmt.Println()
	if uiLanguage == "pt" {
		fmt.Printf("   Palavras FIXAS (OK): %d\n", len(config.FixedWords))
		fmt.Printf("   Palavras FLUTUANTES (NOK/?): %d\n", len(config.FloatingWords))
	} else {
		fmt.Printf("   FIXED words (OK): %d\n", len(config.FixedWords))
		fmt.Printf("   FLOATING words (NOK/?): %d\n", len(config.FloatingWords))
	}
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()
	
	if uiLanguage == "pt" {
		fmt.Print("Pressione Enter para continuar...")
	} else {
		fmt.Print("Press Enter to continue...")
	}
	getUserInput("")
	
	return config
}

func getSeedInput(wordCount int, seedLanguage string) []string {
	clearScreen()
	showSimpleHeader()

	if uiLanguage == "pt" {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    DIGITE OU COLE SUA SEED PHRASE                    ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Printf("Digite ou cole as %d palavras da sua seed (separadas por espaço):\n", wordCount)
		fmt.Println()
			fmt.Println("IMPORTANTE:")
			fmt.Println("   - As palavras DEVEM estar na ORDEM CORRETA")
			fmt.Println("   - Use wildcards para palavras parciais: lugga*, in*ury, injur*")
			fmt.Println("   - Use * (asterisco) para palavras FALTANTES (sem limite)")
			fmt.Println("   - Exemplo: abandon lugga* in*ury * word5 word6 ...")
		fmt.Println()
		fmt.Println("Digite abaixo:")
		fmt.Println()
	} else {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    ENTER OR PASTE YOUR SEED PHRASE                  ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Printf("Enter or paste your %d seed words (space separated):\n", wordCount)
		fmt.Println()
			fmt.Println("IMPORTANT:")
			fmt.Println("   - Words MUST be in CORRECT ORDER")
			fmt.Println("   - Use wildcards for partial words: lugga*, in*ury, injur*")
			fmt.Println("   - Use * (asterisk) for MISSING words (no limit)")
			fmt.Println("   - Example: abandon lugga* in*ury * word5 word6 ...")
		fmt.Println()
		fmt.Println("Type below:")
		fmt.Println()
	}

	input := getUserInput("> ")
	words := strings.Fields(input)

	// Validar quantidade de palavras
	if len(words) != wordCount {
		if uiLanguage == "pt" {
			fmt.Printf("\nERRO: Você digitou %d palavras, mas deveria ter %d palavras!\n", len(words), wordCount)
			fmt.Println("\nPressione Enter para tentar novamente...")
		} else {
			fmt.Printf("\nERROR: You entered %d words, but should have %d words!\n", len(words), wordCount)
			fmt.Println("\nPress Enter to try again...")
		}
		getUserInput("")
		return getSeedInput(wordCount, seedLanguage)
	}

	// Contar asteriscos
	asteriskCount := 0
	for _, word := range words {
		if word == "*" {
			asteriskCount++
		}
	}

	if asteriskCount > 2 {
		if uiLanguage == "pt" {
			fmt.Printf("\n[!]  AVISO: Você usou %d palavras totalmente desconhecidas (*).\n", asteriskCount)
			fmt.Println("   Isso gerará MUITAS combinações e pode levar MUITO tempo!")
			fmt.Println("   Quanto mais *, mais tempo será necessário.")
			fmt.Println()
		} else {
			fmt.Printf("\n[!]  WARNING: You used %d completely unknown words (*).\n", asteriskCount)
			fmt.Println("   This will generate MANY combinations and may take a VERY long time!")
			fmt.Println("   The more *, the more time will be needed.")
			fmt.Println()
		}
	}

	// Validar palavras (exceto asteriscos e wildcards) contra a wordlist BIP39
	wordlist := getWordList(seedLanguage)
	for _, word := range words {
		if word != "*" {
			// Se tem wildcard, validar se pelo menos 1 palavra corresponde
			if strings.Contains(word, "*") {
				matches := getMatchingWords(word, wordlist)
				if len(matches) == 0 {
					if uiLanguage == "pt" {
						fmt.Printf("\nERRO: O padrão '%s' não corresponde a nenhuma palavra BIP39!\n", word)
						fmt.Println("\nVerifique o padrão e tente novamente.")
						fmt.Println("Pressione Enter para continuar...")
					} else {
						fmt.Printf("\nERROR: The pattern '%s' doesn't match any BIP39 word!\n", word)
						fmt.Println("\nCheck the pattern and try again.")
						fmt.Println("Press Enter to continue...")
					}
					getUserInput("")
					return getSeedInput(wordCount, seedLanguage)
				}
			} else {
				// Palavra completa, validar normalmente
				if !contains(wordlist, word) {
					// Buscar sugestões de palavras similares
					suggestions := findSimilarWords(word, wordlist, 5)
					if uiLanguage == "pt" {
						fmt.Printf("\nERRO: A palavra '%s' não está na lista BIP39 de %s!\n", word, seedLanguage)
						if len(suggestions) > 0 {
							fmt.Println("\n   Você quis dizer?")
							for _, s := range suggestions {
								fmt.Printf("   -> %s\n", s)
							}
						}
						fmt.Println("\nVerifique a ortografia e tente novamente.")
						fmt.Println("Pressione Enter para continuar...")
					} else {
						fmt.Printf("\nERROR: The word '%s' is not in the BIP39 %s wordlist!\n", word, seedLanguage)
						if len(suggestions) > 0 {
							fmt.Println("\n   Did you mean?")
							for _, s := range suggestions {
								fmt.Printf("   -> %s\n", s)
							}
						}
						fmt.Println("\nCheck spelling and try again.")
						fmt.Println("Press Enter to continue...")
					}
					getUserInput("")
					return getSeedInput(wordCount, seedLanguage)
				}
			}
		}
	}

	// Mostrar correspondências de wildcards antes de processar
	hasWildcards := false
	for _, word := range words {
		if strings.Contains(word, "*") {
			hasWildcards = true
			break
		}
	}

	if hasWildcards {
		if uiLanguage == "pt" {
			fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
			fmt.Println("║              CORRESPONDÊNCIAS DE WILDCARDS                          ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		} else {
			fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
			fmt.Println("║                WILDCARD MATCHES                                      ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		}

		for _, word := range words {
			if strings.Contains(word, "*") {
				if word == "*" {
					if uiLanguage == "pt" {
						fmt.Printf("\n   %s -> TODAS as %d palavras da wordlist\n", word, len(wordlist))
					} else {
						fmt.Printf("\n   %s -> ALL %d words from wordlist\n", word, len(wordlist))
					}
				} else {
					matches := getMatchingWords(word, wordlist)
					if uiLanguage == "pt" {
						fmt.Printf("\n   %s -> %d palavra(s) encontrada(s):", word, len(matches))
					} else {
						fmt.Printf("\n   %s -> %d word(s) found:", word, len(matches))
					}
					if len(matches) <= 20 {
						fmt.Printf(" %s", strings.Join(matches, ", "))
					} else {
						// Mostrar apenas as 20 primeiras
						fmt.Printf(" %s", strings.Join(matches[:20], ", "))
						if uiLanguage == "pt" {
							fmt.Printf(" ... e mais %d", len(matches)-20)
						} else {
							fmt.Printf(" ... and %d more", len(matches)-20)
						}
					}
					fmt.Println()
				}
			}
		}
		fmt.Println()
	}

	return words
}

func getWordList(language string) []string {
	var wordlistStr string

	switch language {
	case "english":
		wordlistStr = englishWordlist
	case "spanish":
		wordlistStr = spanishWordlist
	case "french":
		wordlistStr = frenchWordlist
	case "italian":
		wordlistStr = italianWordlist
	case "portuguese":
		wordlistStr = portugueseWordlist
	case "japanese":
		wordlistStr = japaneseWordlist
	case "korean":
		wordlistStr = koreanWordlist
	case "chinese_simplified":
		wordlistStr = chineseSimplifiedWordlist
	case "chinese_traditional":
		wordlistStr = chineseTraditionalWordlist
	default:
		wordlistStr = englishWordlist
	}

	// Dividir por linhas e remover vazias
	lines := strings.Split(wordlistStr, "\n")
	wordlist := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			wordlist = append(wordlist, line)
		}
	}

	return wordlist
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ============================================================================
// FUNÇÕES DE CÁLCULO DE PERMUTAÇÕES
// ============================================================================

func factorial(n int) *big.Int {
	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

func calculateTotalPermutations(wordCount int, asteriskCount int) *big.Int {
	// Permutações de n palavras = n!
	totalPerms := factorial(wordCount)

	// Se tem asteriscos, multiplica por 2048^asteriskCount
	if asteriskCount > 0 {
		wordlistSize := big.NewInt(2048)
		multiplier := new(big.Int).Exp(wordlistSize, big.NewInt(int64(asteriskCount)), nil)
		totalPerms.Mul(totalPerms, multiplier)
	}

	return totalPerms
}

func estimateValidSeeds(totalPermutations *big.Int, wordCount int) *big.Int {
	// Aproximadamente 1 em cada 256 permutações é válida (checksum de 8 bits para 12 palavras)
	// Para 12 palavras: 1/256
	// Para 15 palavras: 1/32
	// Para 18 palavras: 1/64
	// Para 21 palavras: 1/128
	// Para 24 palavras: 1/256

	var divisor int64
	switch wordCount {
	case 12:
		divisor = 16 // Aproximadamente 1 em 16 (4 bits de checksum)
	case 15:
		divisor = 32 // Aproximadamente 1 em 32 (5 bits de checksum)
	case 18:
		divisor = 64 // Aproximadamente 1 em 64 (6 bits de checksum)
	case 21:
		divisor = 128 // Aproximadamente 1 em 128 (7 bits de checksum)
	case 24:
		divisor = 256 // Aproximadamente 1 em 256 (8 bits de checksum)
	default:
		divisor = 16
	}

	estimated := new(big.Int).Div(totalPermutations, big.NewInt(divisor))
	return estimated
}

func estimateFiles(estimatedValidSeeds *big.Int) int64 {
	rowsPerFile := big.NewInt(int64(ROWS_PER_FILE))
	files := new(big.Int).Div(estimatedValidSeeds, rowsPerFile)
	
	// Se tem resto, adiciona mais um arquivo
	remainder := new(big.Int).Mod(estimatedValidSeeds, rowsPerFile)
	if remainder.Cmp(big.NewInt(0)) > 0 {
		files.Add(files, big.NewInt(1))
	}

	// Se for muito grande, retorna -1 para indicar "muitos"
	if files.Cmp(big.NewInt(1000000)) > 0 {
		return -1
	}

	return files.Int64()
}

func formatBigNumber(n *big.Int) string {
	str := n.String()
	
	// Adicionar separadores de milhar
	var result []rune
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, digit)
	}
	
	return string(result)
}

func estimateTime(totalPermutations *big.Int) string {
	// Taxa média: ~200.000 permutações/segundo
	rate := 200000.0
	
	permsFloat := new(big.Float).SetInt(totalPermutations)
	seconds := new(big.Float).Quo(permsFloat, big.NewFloat(rate))
	
	secondsFloat, _ := seconds.Float64()
	
	if secondsFloat < 60 {
		return fmt.Sprintf("%.0f segundos", secondsFloat)
	} else if secondsFloat < 3600 {
		return fmt.Sprintf("%.1f minutos", secondsFloat/60)
	} else if secondsFloat < 86400 {
		return fmt.Sprintf("%.1f horas", secondsFloat/3600)
	} else if secondsFloat < 2592000 {
		return fmt.Sprintf("%.1f dias", secondsFloat/86400)
	} else if secondsFloat < 31536000 {
		return fmt.Sprintf("%.1f meses", secondsFloat/2592000)
	} else {
		return fmt.Sprintf("%.1f anos", secondsFloat/31536000)
	}
}

// ============================================================================
// FUNÇÕES DE GERAÇÃO DE PERMUTAÇÕES
// ============================================================================

type PermutationGenerator struct {
	words           []string
	wordlist        []string
	asteriskCount   int
	asteriskPos     []int
	checked         int
	found           int
	startTime       time.Time
	currentFile     *excelize.File
	currentFileNum  int
	currentRow      int
	filesCreated    []string
	dateStamp       string
	outputDir       string
	resultMap       map[string]bool
	seedLanguage    string
	config          *SeedConfig
	isAdvancedMode  bool
}

func NewPermutationGenerator(words []string, wordlist []string, seedLanguage string) *PermutationGenerator {
	dateStamp := time.Now().Format("20060102_150405")
	
	// Criar pasta com data
	var folderName string
	if uiLanguage == "pt" {
		folderName = fmt.Sprintf("Resultados_%s", dateStamp)
	} else {
		folderName = fmt.Sprintf("Results_%s", dateStamp)
	}
	
	if err := os.MkdirAll(folderName, 0755); err != nil {
		if uiLanguage == "pt" {
			fmt.Printf("\nAviso: Não foi possível criar pasta %s: %v\n", folderName, err)
			fmt.Println("   Arquivos serão salvos no diretório atual.\n")
		} else {
			fmt.Printf("\nWarning: Could not create folder %s: %v\n", folderName, err)
			fmt.Println("   Files will be saved in current directory.\n")
		}
		folderName = "."
	}
	
	pg := &PermutationGenerator{
		words:        words,
		wordlist:     wordlist,
		startTime:    time.Now(),
		currentRow:   2, // Linha 1 é o cabeçalho
		dateStamp:    dateStamp,
		outputDir:    folderName,
		resultMap:    make(map[string]bool),
		filesCreated: []string{},
		seedLanguage: seedLanguage,
	}

	// Contar asteriscos
	for i, word := range words {
		if word == "*" {
			pg.asteriskCount++
			pg.asteriskPos = append(pg.asteriskPos, i)
		}
	}

	return pg
}

func (pg *PermutationGenerator) createNewFile() error {
	// Salvar arquivo anterior se existir
	if pg.currentFile != nil {
		pg.saveCurrentFile()
	}

	// Criar novo arquivo
	pg.currentFileNum++
	pg.currentRow = 2 // Reset para linha 2 (após cabeçalho)

	f := excelize.NewFile()
	sheetName := "Valid Seeds"

	// Criar sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Definir cabeçalhos
	f.SetCellValue(sheetName, "A1", "Index")
	f.SetCellValue(sheetName, "B1", "Seed Phrase")

	// Estilo do cabeçalho
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	f.SetCellStyle(sheetName, "A1", "B1", headerStyle)

	// Ajustar largura das colunas
	f.SetColWidth(sheetName, "A", "A", 10)
	f.SetColWidth(sheetName, "B", "B", 120)

	// Definir sheet ativa
	f.SetActiveSheet(index)

	// Deletar Sheet1 padrão
	f.DeleteSheet("Sheet1")

	pg.currentFile = f

	return nil
}

func (pg *PermutationGenerator) saveCurrentFile() {
	if pg.currentFile == nil {
		return
	}

	var filename string
	var prefix string
	if uiLanguage == "pt" {
		prefix = "Combinacoes_Possiveis"
	} else {
		prefix = "Possible_Combinations"
	}
	
	var basename string
	if pg.currentFileNum == 1 && pg.found < ROWS_PER_FILE {
		// Único arquivo ou arquivo final
		basename = fmt.Sprintf("%s_FINAL_%s.xlsx", prefix, pg.dateStamp)
	} else {
		basename = fmt.Sprintf("%s_%02d_%s.xlsx", prefix, pg.currentFileNum, pg.dateStamp)
	}
	
	// Caminho completo com pasta
	filename = fmt.Sprintf("%s/%s", pg.outputDir, basename)

	if err := pg.currentFile.SaveAs(filename); err != nil{
		if uiLanguage == "pt" {
			fmt.Printf("\nERRO ao salvar arquivo %s: %v\n", filename, err)
		} else {
			fmt.Printf("\nERROR saving file %s: %v\n", filename, err)
		}
	} else {
		pg.filesCreated = append(pg.filesCreated, filename)
	}

	pg.currentFile = nil
}

func (pg *PermutationGenerator) addSeed(mnemonic string) error {
	// Verificar duplicata
	if pg.resultMap[mnemonic] {
		return nil
	}

	pg.resultMap[mnemonic] = true
	pg.found++

	// Criar novo arquivo se necessário
	if pg.currentFile == nil || pg.currentRow > ROWS_PER_FILE+1 {
		if err := pg.createNewFile(); err != nil {
			return err
		}

		if uiLanguage == "pt" {
			fmt.Printf("\n>>> Criando arquivo %d...\n", pg.currentFileNum)
		} else {
			fmt.Printf("\n>>> Creating file %d...\n", pg.currentFileNum)
		}
	}

	sheetName := "Valid Seeds"
	globalIndex := pg.found

	// Adicionar dados
	pg.currentFile.SetCellValue(sheetName, fmt.Sprintf("A%d", pg.currentRow), globalIndex)
	pg.currentFile.SetCellValue(sheetName, fmt.Sprintf("B%d", pg.currentRow), mnemonic)

	pg.currentRow++

	return nil
}

func (pg *PermutationGenerator) Generate() int {
	// Modo Avançado: permutar apenas palavras flutuantes
	if pg.isAdvancedMode && pg.config != nil {
		return pg.generateAdvanced()
	}
	
	// EXPANSÃO DE WILDCARDS (Modo Simple)
	// Verificar se há wildcards parciais (não apenas *)
	hasPartialWildcards := false
	for _, word := range pg.words {
		if word != "*" && strings.Contains(word, "*") {
			hasPartialWildcards = true
			break
		}
	}
	
	// Se tem wildcards parciais OU 3+ asteriscos puros, usar expansão
	if hasPartialWildcards || pg.asteriskCount > 2 {
		return pg.generateWithWildcardExpansion()
	}
	
		if pg.asteriskCount == 0 {
			// Sem asteriscos: MODO 1 = ordem fixa, apenas validar
			if uiLanguage == "pt" {
				fmt.Println("\nModo 1: Ordem conhecida - Validando seed...")
			} else {
				fmt.Println("\nMode 1: Known order - Validating seed...")
			}
	
			// Validar seed na ordem fixa (sem permutações)
			pg.validateSeed(pg.words)

	} else if pg.asteriskCount == 1 {
		// 1 asterisco
		if uiLanguage == "pt" {
			fmt.Println("\nGerando combinações com 1 palavra faltante...")
			fmt.Printf("   Testando 2048 palavras possíveis...\n")
		} else {
			fmt.Println("\nGenerating combinations with 1 missing word...")
			fmt.Printf("   Testing 2048 possible words...\n")
		}

		pos := pg.asteriskPos[0]
		for idx, word := range pg.wordlist {
			if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
				break
			}

			if idx%100 == 0 {
				if uiLanguage == "pt" {
					fmt.Printf("\r   Progresso: %d/2048 palavras testadas | %d seeds válidas encontradas", idx, pg.found)
				} else {
					fmt.Printf("\r   Progress: %d/2048 words tested | %d valid seeds found", idx, pg.found)
				}
			}

			tempWords := make([]string, len(pg.words))
			copy(tempWords, pg.words)
			tempWords[pos] = word

			// Modo 1: validar na ordem fixa (sem permutações)
			pg.validateSeed(tempWords)
		}
		fmt.Println()

	} else if pg.asteriskCount == 2 {
		// 2 asteriscos
		if uiLanguage == "pt" {
			fmt.Println("\nGerando combinações com 2 palavras faltantes...")
			fmt.Printf("   Testando 2048 x 2048 = 4.194.304 combinações possíveis...\n")
			fmt.Println("   ISSO PODE DEMORAR MUITO! Aguarde...")
		} else {
			fmt.Println("\nGenerating combinations with 2 missing words...")
			fmt.Printf("   Testing 2048 x 2048 = 4,194,304 possible combinations...\n")
			fmt.Println("   THIS MAY TAKE VERY LONG! Please wait...")
		}

		pos1 := pg.asteriskPos[0]
		pos2 := pg.asteriskPos[1]

		totalCombinations := len(pg.wordlist) * len(pg.wordlist)
		currentCombination := 0

		for _, word1 := range pg.wordlist {
			if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
				break
			}

			for _, word2 := range pg.wordlist {
				if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
					break
				}

				currentCombination++
				if currentCombination%10000 == 0 {
					elapsed := time.Since(pg.startTime).Seconds()
					rate := float64(currentCombination) / elapsed
					remaining := float64(totalCombinations-currentCombination) / rate

					if uiLanguage == "pt" {
						fmt.Printf("\r   Progresso: %d/%d combinações (%.1f%%) | %d válidas | Tempo restante: %.0fs   ",
							currentCombination, totalCombinations,
							float64(currentCombination)*100/float64(totalCombinations),
							pg.found, remaining)
					} else {
						fmt.Printf("\r   Progress: %d/%d combinations (%.1f%%) | %d valid | Time remaining: %.0fs   ",
							currentCombination, totalCombinations,
							float64(currentCombination)*100/float64(totalCombinations),
							pg.found, remaining)
					}
				}

				tempWords := make([]string, len(pg.words))
				copy(tempWords, pg.words)
				tempWords[pos1] = word1
				tempWords[pos2] = word2

				// Modo 1: validar na ordem fixa (sem permutações)
				pg.validateSeed(tempWords)
			}
		}
		fmt.Println()
	}

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
		
		// Formato: Combinacoes_Possiveis_60FINAL_20251120_135710.xlsx
		newBasename := fmt.Sprintf("%s_%02dFINAL_%s.xlsx", prefix, pg.currentFileNum, pg.dateStamp)
		newName := fmt.Sprintf("%s/%s", pg.outputDir, newBasename)
		
		os.Rename(lastFile, newName)
		pg.filesCreated[len(pg.filesCreated)-1] = newName
	}

	return pg.found
}

// generateWithWildcardExpansion gera combinações expandindo wildcards
func (pg *PermutationGenerator) generateWithWildcardExpansion() int {
	if uiLanguage == "pt" {
		fmt.Println("\nExpandindo wildcards...")
	} else {
		fmt.Println("\nExpanding wildcards...")
	}
	
	// Expandir cada palavra
	expansions := expandWildcardsInWords(pg.words, pg.wordlist)
	
	// Mostrar quantas combinações serão geradas
	totalCombinations := countTotalCombinations(expansions)
	
	if uiLanguage == "pt" {
		fmt.Printf("   Total de combinações a testar: %d\n", totalCombinations)
		fmt.Println("\nGerando e validando combinações...")
	} else {
		fmt.Printf("   Total combinations to test: %d\n", totalCombinations)
		fmt.Println("\nGenerating and validating combinations...")
	}
	
	// Gerar produto cartesiano
	combinations := generateCartesianProduct(expansions)
	
	// Testar cada combinação
	for idx, combination := range combinations {
		if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
			break
		}
		
		pg.checked++
		
		// Mostrar progresso a cada 1000 combinações
		if idx > 0 && idx%1000 == 0 {
			elapsed := time.Since(pg.startTime).Seconds()
			rate := float64(pg.checked) / elapsed
			remaining := float64(len(combinations)-idx) / rate
			
			if uiLanguage == "pt" {
				fmt.Printf("\r   Progresso: %d/%d (%.1f%%) | Válidas: %d | Tempo restante: %.0fs   ",
					idx, len(combinations),
					float64(idx)*100/float64(len(combinations)),
					pg.found, remaining)
			} else {
				fmt.Printf("\r   Progress: %d/%d (%.1f%%) | Valid: %d | Time remaining: %.0fs   ",
					idx, len(combinations),
					float64(idx)*100/float64(len(combinations)),
					pg.found, remaining)
			}
		}
		
		// Validar BIP39
		mnemonic := strings.Join(combination, " ")
		
		// Evitar duplicatas
		if pg.resultMap[mnemonic] {
			continue
		}
		
		if isValidSeed(mnemonic) {
			pg.resultMap[mnemonic] = true
			pg.found++
			
			// Criar novo arquivo se necessário
			if pg.currentFile == nil {
				pg.createNewFile()
			}
			
			// Adicionar ao Excel
			pg.currentFile.SetCellValue("Valid Seeds", fmt.Sprintf("A%d", pg.currentRow), pg.found)
			pg.currentFile.SetCellValue("Valid Seeds", fmt.Sprintf("B%d", pg.currentRow), mnemonic)
			pg.currentRow++
			
			// Salvar e criar novo arquivo se atingiu limite
			if pg.currentRow > ROWS_PER_FILE+1 {
				pg.saveCurrentFile()
				pg.createNewFile()
			}
		}
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

// validateSeed valida uma seed na ordem fixa (sem permutações) - Modo 1
func (pg *PermutationGenerator) validateSeed(arr []string) {
	if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
		return
	}

	pg.checked++

	mnemonic := strings.Join(arr, " ")

	// Definir wordlist correta para validação
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

	if isValidSeed(mnemonic) {
		pg.addSeed(mnemonic)
		if uiLanguage == "pt" {
			fmt.Printf("\n   Seed válida encontrada: %s\n", mnemonic)
		} else {
			fmt.Printf("\n   Valid seed found: %s\n", mnemonic)
		}
	} else {
		if uiLanguage == "pt" {
			fmt.Printf("\n   Seed inválida (checksum incorreto): %s\n", mnemonic)
		} else {
			fmt.Printf("\n   Invalid seed (incorrect checksum): %s\n", mnemonic)
		}
	}
}

func (pg *PermutationGenerator) permuteWithLimit(arr []string, start int) {
	if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
		return
	}

	if start == len(arr)-1 {
		pg.checked++

		// Mostrar progresso a cada 100k verificações
		if pg.checked%100000 == 0 {
			elapsed := time.Since(pg.startTime).Seconds()
			rate := float64(pg.checked) / elapsed

			if uiLanguage == "pt" {
				fmt.Printf("\r   Verificadas: %d permutações | Válidas: %d | Taxa: %.0f/s   ", pg.checked, pg.found, rate)
			} else {
				fmt.Printf("\r   Checked: %d permutations | Valid: %d | Rate: %.0f/s   ", pg.checked, pg.found, rate)
			}
		}

		mnemonic := strings.Join(arr, " ")

		// Definir wordlist correta para validação
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

		if isValidSeed(mnemonic) {
			pg.addSeed(mnemonic)

			// Parar se atingir o limite
			if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
				if uiLanguage == "pt" {
					fmt.Printf("\n\nLIMITE ATINGIDO: %d seeds válidas encontradas!\n", MAX_PERMUTATIONS_TO_GENERATE)
					fmt.Println("   O processamento foi interrompido por segurança.")
				} else {
					fmt.Printf("\n\nLIMIT REACHED: %d valid seeds found!\n", MAX_PERMUTATIONS_TO_GENERATE)
					fmt.Println("   Processing stopped for safety.")
				}
			}
		}
		return
	}

	for i := start; i < len(arr); i++ {
		if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
			return
		}

		arr[start], arr[i] = arr[i], arr[start]
		pg.permuteWithLimit(arr, start+1)
		arr[start], arr[i] = arr[i], arr[start]
	}
}

// ============================================================================
// FUNÇÃO PRINCIPAL
// ============================================================================

func main() {
	checkLicense()
	clearScreen()
	showWelcomeHeader()
	chooseUILanguage()
	showUnmixerSeedSearchInfo()

	// Escolher idioma da seed
	seedLanguage := chooseSeedLanguage()

	// Escolher quantidade de palavras
	wordCount := chooseSeedWordCount()

	// Escolher tipo de validação
	chooseValidationType()

	// Escolher modo de input
	inputMode := chooseInputMode()

	// Obter input da seed
	var words []string
	var config *SeedConfig
	
	if inputMode == "advanced_partial" {
		config = getSeedInputAdvancedPartial(wordCount, seedLanguage)
		// Extrair palavras do config para compatibilidade
		for _, pos := range config.Positions {
			words = append(words, pos.Word)
		}
	} else if inputMode == "advanced_complete" {
		config = getSeedInputAdvancedComplete(wordCount, seedLanguage)
		// Extrair palavras do config para compatibilidade
		for _, pos := range config.Positions {
			words = append(words, pos.Word)
		}
	} else if inputMode == "descrambler" {
		words = getSeedInputDescrambler(wordCount, seedLanguage)
	} else {
		words = getSeedInput(wordCount, seedLanguage)
	}

	// Contar asteriscos e wildcards
	asteriskCount := 0
	wildcardCount := 0
	for _, word := range words {
		if word == "*" {
			asteriskCount++
		}
		if strings.Contains(word, "*") {
			wildcardCount++
		}
	}

	// CALCULAR E MOSTRAR ESTIMATIVAS
	clearScreen()
	showSimpleHeader()

	if uiLanguage == "pt" {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    CÁLCULO DE PERMUTAÇÕES                            ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
	} else {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    PERMUTATION CALCULATION                           ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
	}

	// Mostrar resumo da configuração
	var floatingCount int
	if inputMode == "advanced_partial" || inputMode == "advanced_complete" && config != nil {
		// Modo Avançado: contar palavras fixas (OK) e flutuantes (NOK/?)
		fixedCount := 0
		var fixedWords []string
		var floatingWords []string
		
		for _, pos := range config.Positions {
			if pos.Status == "OK" {
				fixedCount++
				fixedWords = append(fixedWords, fmt.Sprintf("Pos %d: %s", pos.Position, pos.Word))
			} else {
				floatingWords = append(floatingWords, pos.Word)
			}
		}
		floatingCount = wordCount - fixedCount
		
		if uiLanguage == "pt" {
			fmt.Printf(">>> RESUMO DA CONFIGURAÇÃO:\n\n")
			fmt.Printf("   Total de palavras: %d\n", wordCount)
			validationLabel := "BIP39 (Padrão)"
			if validationType == "electrum" {
				validationLabel = "Electrum / Electron Cash"
			} else if validationType == "none" {
				validationLabel = "Sem Validação (Força Bruta)"
			}
			fmt.Printf("   Tipo de validação: %s\n", validationLabel)
			fmt.Printf("   Palavras FIXAS (OK): %d\n", fixedCount)
			if fixedCount > 0 {
				for _, fw := range fixedWords {
					fmt.Printf("      - %s\n", fw)
				}
			}
			fmt.Printf("   Palavras FLUTUANTES (NOK/?): %d\n", floatingCount)
			if floatingCount > 0 {
				fmt.Printf("      Palavras: %s\n", strings.Join(floatingWords, ", "))
			}
			fmt.Println()
			fmt.Println(">>> CALCULANDO ESTIMATIVAS...")
			fmt.Println()
		} else {
			fmt.Printf(">>> CONFIGURATION SUMMARY:\n\n")
			fmt.Printf("   Total words: %d\n", wordCount)
			validationLabelEN := "BIP39 (Standard)"
			if validationType == "electrum" {
				validationLabelEN = "Electrum / Electron Cash"
			} else if validationType == "none" {
				validationLabelEN = "No Validation (Brute Force)"
			}
			fmt.Printf("   Validation type: %s\n", validationLabelEN)
			fmt.Printf("   FIXED words (OK): %d\n", fixedCount)
			if fixedCount > 0 {
				for _, fw := range fixedWords {
					fmt.Printf("      - %s\n", fw)
				}
			}
			fmt.Printf("   FLOATING words (NOK/?): %d\n", floatingCount)
			if floatingCount > 0 {
				fmt.Printf("      Words: %s\n", strings.Join(floatingWords, ", "))
			}
			fmt.Println()
			fmt.Println(">>> CALCULATING ESTIMATES...")
			fmt.Println()
		}
	} else {
		// Modo Simples (Modo 1 - ordem fixa)
		floatingCount = wordCount
		if uiLanguage == "pt" {
			fmt.Printf(">>> RESUMO DA CONFIGURAÇÃO:\n\n")
			fmt.Printf("   Total de palavras: %d\n", wordCount)
			validationLabel2 := "BIP39 (Padrão)"
			if validationType == "electrum" {
				validationLabel2 = "Electrum / Electron Cash"
			} else if validationType == "none" {
				validationLabel2 = "Sem Validação (Força Bruta)"
			}
			fmt.Printf("   Tipo de validação: %s\n", validationLabel2)
			fmt.Printf("   Palavras completas: %d\n", wordCount - wildcardCount)
			if wildcardCount > 0 {
				fmt.Printf("   Palavras incompletas/wildcards: %d\n", wildcardCount)
			}
			fmt.Printf("   Seed: %s\n", strings.Join(words, " "))
			fmt.Println()
			fmt.Println(">>> CALCULANDO ESTIMATIVAS...")
			fmt.Println()
		} else {
			fmt.Printf(">>> CONFIGURATION SUMMARY:\n\n")
			fmt.Printf("   Total words: %d\n", wordCount)
			validationLabelEN2 := "BIP39 (Standard)"
			if validationType == "electrum" {
				validationLabelEN2 = "Electrum / Electron Cash"
			} else if validationType == "none" {
				validationLabelEN2 = "No Validation (Brute Force)"
			}
			fmt.Printf("   Validation type: %s\n", validationLabelEN2)
			fmt.Printf("   Complete words: %d\n", wordCount - wildcardCount)
			if wildcardCount > 0 {
				fmt.Printf("   Incomplete words/wildcards: %d\n", wildcardCount)
			}
			fmt.Printf("   Seed: %s\n", strings.Join(words, " "))
			fmt.Println()
			fmt.Println(">>> CALCULATING ESTIMATES...")
			fmt.Println()
		}
	}

	// Calcular permutações/combinações
	var totalPerms *big.Int
	if inputMode == "advanced_partial" || inputMode == "advanced_complete" && config != nil {
		// Modo Avançado: calcular apenas permutações das palavras flutuantes
		totalPerms = calculateTotalPermutations(floatingCount, asteriskCount)
	} else {
		// Modo 1 (Simples): calcular APENAS combinações de wildcards (SEM permutações)
		wordlist := getWordList(seedLanguage)
		totalPerms = calculateMode1Combinations(words, wordlist)
	}
	estimatedValid := estimateValidSeeds(totalPerms, wordCount)
	estimatedFiles := estimateFiles(estimatedValid)
	estimatedTimeStr := estimateTime(totalPerms)

	if uiLanguage == "pt" {
		fmt.Printf("Total de permutações possíveis: %s\n", formatBigNumber(totalPerms))
		fmt.Printf("Seeds válidas estimadas: %s\n", formatBigNumber(estimatedValid))
		
		if estimatedFiles == -1 {
			fmt.Printf(">>> Arquivos estimados: MUITOS (milhares+)\n")
		} else if estimatedFiles == 0 {
			fmt.Printf(">>> Arquivos estimados: 1\n")
		} else {
			fmt.Printf(">>> Arquivos estimados: %d\n", estimatedFiles)
		}
		
		fmt.Printf(">>>  Tempo estimado: %s\n", estimatedTimeStr)
		fmt.Println()
		
		// Verificar se é viável
		if totalPerms.Cmp(big.NewInt(1000000000000)) > 0 { // > 1 trilhão
			fmt.Println("AVISO CRÍTICO:")
			fmt.Println("   Esta quantidade de permutações é IMPRATICÁVEL!")
			fmt.Println("   Pode levar ANOS para processar.")
			fmt.Println("   Recomendamos:")
			fmt.Println("   - Usar seeds de 12 palavras")
			fmt.Println("   - Evitar palavras faltantes")
			fmt.Println("   - Ter certeza das palavras antes de processar")
			fmt.Println()
		} else if totalPerms.Cmp(big.NewInt(10000000000)) > 0 { // > 10 bilhões
			fmt.Println("AVISO:")
			fmt.Println("   Esta quantidade de permutações é MUITO GRANDE!")
			fmt.Println("   Pode levar DIAS ou SEMANAS para processar.")
			fmt.Println()
		}

		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Deseja continuar com o processamento? (S/N)")
	} else {
		fmt.Printf("Total possible permutations: %s\n", formatBigNumber(totalPerms))
		fmt.Printf("Estimated valid seeds: %s\n", formatBigNumber(estimatedValid))
		
		if estimatedFiles == -1 {
			fmt.Printf(">>> Estimated files: MANY (thousands+)\n")
		} else if estimatedFiles == 0 {
			fmt.Printf(">>> Estimated files: 1\n")
		} else {
			fmt.Printf(">>> Estimated files: %d\n", estimatedFiles)
		}
		
		fmt.Printf(">>>  Estimated time: %s\n", estimatedTimeStr)
		fmt.Println()
		
		// Verificar se é viável
		if totalPerms.Cmp(big.NewInt(1000000000000)) > 0 { // > 1 trillion
			fmt.Println("CRITICAL WARNING:")
			fmt.Println("   This amount of permutations is IMPRACTICAL!")
			fmt.Println("   May take YEARS to process.")
			fmt.Println("   We recommend:")
			fmt.Println("   - Use 12-word seeds")
			fmt.Println("   - Avoid missing words")
			fmt.Println("   - Be sure of the words before processing")
			fmt.Println()
		} else if totalPerms.Cmp(big.NewInt(10000000000)) > 0 { // > 10 billion
			fmt.Println("WARNING:")
			fmt.Println("   This amount of permutations is VERY LARGE!")
			fmt.Println("   May take DAYS or WEEKS to process.")
			fmt.Println()
		}

		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Do you want to continue with processing? (Y/N)")
	}

	for {
		confirm := getUserInput("> ")
		confirmUpper := strings.ToUpper(confirm)
		if confirmUpper == "S" || confirmUpper == "Y" {
			break
		}
		if confirmUpper == "N" {
			if uiLanguage == "pt" {
				fmt.Println("\nOperação cancelada pelo usuário.")
			} else {
				fmt.Println("\nOperation cancelled by user.")
			}
			return
		}
		if uiLanguage == "pt" {
			fmt.Println("\n  [X] Opcao invalida! Digite S (Sim) ou N (Nao).")
		} else {
			fmt.Println("\n  [X] Invalid option! Type Y (Yes) or N (No).")
		}
		fmt.Println()
	}

	// Processar
	if uiLanguage == "pt" {
		fmt.Println("\nIniciando processamento...")
		fmt.Println("   Isso pode demorar. Aguarde...")
		fmt.Println()
	} else {
		fmt.Println("\nStarting processing...")
		fmt.Println("   This may take a while. Please wait...")
		fmt.Println()
	}

	wordlist := getWordList(seedLanguage)
	startTime := time.Now()

	pg := NewPermutationGenerator(words, wordlist, seedLanguage)
	var totalFound int
	if inputMode == "descrambler" {
		totalFound = pg.generateDescrambler()
	} else {
		if inputMode == "advanced_partial" || inputMode == "advanced_complete" && config != nil {
			pg.config = config
			pg.isAdvancedMode = true
		}
		totalFound = pg.Generate()
	}

	elapsed := time.Since(startTime)

	fmt.Println()

	if totalFound == 0 {
		if uiLanguage == "pt" {
			fmt.Println("\nNENHUMA SEED VÁLIDA ENCONTRADA!")
			fmt.Println("   Possíveis causas:")
			fmt.Println("   - As palavras podem estar incorretas")
			fmt.Println("   - Pode haver um erro de digitação")
			fmt.Println("   - A seed pode não ser BIP39 válida")
			fmt.Println()
			fmt.Println("   Verifique as palavras e tente novamente.")
		} else {
			fmt.Println("\nNO VALID SEEDS FOUND!")
			fmt.Println("   Possible causes:")
			fmt.Println("   - The words may be incorrect")
			fmt.Println("   - There may be a typo")
			fmt.Println("   - The seed may not be BIP39 valid")
			fmt.Println()
			fmt.Println("   Check the words and try again.")
		}
		fmt.Println()
		if uiLanguage == "pt" {
			fmt.Println("Pressione Enter para sair...")
		} else {
			fmt.Println("Press Enter to exit...")
		}
		getUserInput("")
		return
	}

	// Sucesso
	clearScreen()
	showSimpleHeader()

	if uiLanguage == "pt" {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    PROCESSAMENTO CONCLUÍDO!                       ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Printf("   Seeds válidas encontradas: %d\n", totalFound)
		fmt.Printf("   >>>  Tempo de processamento: %.2f segundos (%.2f minutos)\n", elapsed.Seconds(), elapsed.Minutes())
		fmt.Printf("   >>> Arquivos criados: %d\n", len(pg.filesCreated))
		fmt.Println()
		fmt.Println("   Arquivos gerados:")
		for i, filename := range pg.filesCreated {
			fmt.Printf("   %d. %s\n", i+1, filename)
		}
		fmt.Println()
		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("PRÓXIMOS PASSOS:")
		fmt.Println()
		fmt.Println("   1. Abra os arquivos Excel gerados")
		fmt.Println("   2. Copie os arquivos Excel para a pasta 'Importar_Seeds' do CIE")
		fmt.Println("      (Crypto Hunter Pro - CIE Crypto Intelligence Engine)")
		fmt.Println("   3. No CIE, selecione 'Importar Excel do Unmixer Seed'")
		fmt.Println("   4. O CIE testara automaticamente o saldo de cada seed")
		fmt.Println("      em multiplas blockchains (BTC, ETH, SOL, TRX, etc.)")
		fmt.Println()
		fmt.Println("   OU, se preferir testar manualmente:")
		fmt.Println("   5. Teste cada seed phrase em sua carteira")
		fmt.Println("   6. Verifique qual delas restaura seus fundos")
		fmt.Println()
		fmt.Println("SEGURANÇA:")
		fmt.Println()
		fmt.Println("   - NUNCA compartilhe suas seed phrases com ninguém")
		fmt.Println("   - Mantenha os arquivos Excel em local seguro")
		fmt.Println("   - Delete os arquivos após encontrar a seed correta")
		fmt.Println()
		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Pressione Enter para sair...")
	} else {
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    PROCESSING COMPLETED!                          ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Printf("   Valid seeds found: %d\n", totalFound)
		fmt.Printf("   >>>  Processing time: %.2f seconds (%.2f minutes)\n", elapsed.Seconds(), elapsed.Minutes())
		fmt.Printf("   >>> Files created: %d\n", len(pg.filesCreated))
		fmt.Println()
		fmt.Println("   Generated files:")
		for i, filename := range pg.filesCreated {
			fmt.Printf("   %d. %s\n", i+1, filename)
		}
		fmt.Println()
		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("NEXT STEPS:")
		fmt.Println()
		fmt.Println("   1. Open the generated Excel files")
		fmt.Println("   2. Copy the Excel files to the 'Import_Seeds' folder of CIE")
		fmt.Println("      (Crypto Hunter Pro - CIE Crypto Intelligence Engine)")
		fmt.Println("   3. In CIE, select 'Import Unmixer Seed Excel'")
		fmt.Println("   4. CIE will automatically test the balance of each seed")
		fmt.Println("      across multiple blockchains (BTC, ETH, SOL, TRX, etc.)")
		fmt.Println()
		fmt.Println("   OR, if you prefer to test manually:")
		fmt.Println("   5. Test each seed phrase in your wallet")
		fmt.Println("   6. Check which one restores your funds")
		fmt.Println()
		fmt.Println("SECURITY:")
		fmt.Println()
		fmt.Println("   - NEVER share your seed phrases with anyone")
		fmt.Println("   - Keep the Excel files in a secure location")
		fmt.Println("   - Delete the files after finding the correct seed")
		fmt.Println()
		fmt.Println("════════════════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Press Enter to exit...")
	}

	getUserInput("")
}

func (pg *PermutationGenerator) generateAdvanced() int {
// EXPANSÃO DE WILDCARDS (Modo Advanced Partial)
// Expandir wildcards nas palavras antes de permutar
if pg.config.InputMode == "advanced_partial" {
	// Verificar se há wildcards
	hasWildcards := false
	for _, pos := range pg.config.Positions {
		if strings.Contains(pos.Word, "*") {
			hasWildcards = true
			break
		}
	}
	
	if hasWildcards {
		return pg.generateAdvancedWithWildcards()
	}
}

if uiLanguage == "pt" {
fmt.Println("\nGerando permutações (Modo Avançado)...")
fmt.Println("   Apenas palavras flutuantes serão permutadas.")
} else {
fmt.Println("\nGenerating permutations (Advanced Mode)...")
fmt.Println("   Only floating words will be permuted.")
}

// Separar palavras fixas e flutuantes
var floatingWords []string
var floatingPositions []int
fixedMap := make(map[int]string)

for _, pos := range pg.config.Positions {
if pos.Status == "OK" {
fixedMap[pos.Position-1] = pos.Word
} else {
floatingWords = append(floatingWords, pos.Word)
floatingPositions = append(floatingPositions, pos.Position-1)
}
}

	if len(floatingWords) == 0 {
		testWords := make([]string, len(pg.config.Positions))
		for i, pos := range pg.config.Positions {
			testWords[i] = pos.Word
		}
		// Validar seed completa
		mnemonic := strings.Join(testWords, " ")
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
		if isValidSeed(mnemonic) {
			pg.addSeed(mnemonic)
		}
} else {
pg.permuteAdvanced(floatingWords, floatingPositions, fixedMap, 0)
}

pg.saveCurrentFile()

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

func (pg *PermutationGenerator) permuteAdvanced(floatingWords []string, floatingPositions []int, fixedMap map[int]string, start int) {
if pg.found >= MAX_PERMUTATIONS_TO_GENERATE {
return
}

	if start == len(floatingWords) {
		pg.checked++
		if pg.checked%100000 == 0 {
			elapsed := time.Since(pg.startTime).Seconds()
			rate := float64(pg.checked) / elapsed
			if uiLanguage == "pt" {
				fmt.Printf("\r   Verificadas: %d permutações | Válidas: %d | Taxa: %.0f/s   ", pg.checked, pg.found, rate)
			} else {
				fmt.Printf("\r   Checked: %d permutations | Valid: %d | Rate: %.0f/s   ", pg.checked, pg.found, rate)
			}
		}
		
		testWords := make([]string, len(pg.config.Positions))
		for pos, word := range fixedMap {
			testWords[pos] = word
		}
		for i, pos := range floatingPositions {
			testWords[pos] = floatingWords[i]
		}
		
		mnemonic := strings.Join(testWords, " ")
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
		
		if isValidSeed(mnemonic) {
			pg.addSeed(mnemonic)
		}
		return
	}

for i := start; i < len(floatingWords); i++ {
floatingWords[start], floatingWords[i] = floatingWords[i], floatingWords[start]
pg.permuteAdvanced(floatingWords, floatingPositions, fixedMap, start+1)
floatingWords[start], floatingWords[i] = floatingWords[i], floatingWords[start]
}
}

// findSimilarWords encontra palavras similares na wordlist usando distância de Levenshtein simplificada
func findSimilarWords(input string, wordlist []string, maxResults int) []string {
	type wordScore struct {
		word  string
		score int
	}

	var results []wordScore

	for _, word := range wordlist {
		score := levenshteinDistance(input, word)
		// Considerar apenas palavras com distância <= 3
		if score <= 3 {
			results = append(results, wordScore{word, score})
		}
	}

	// Ordenar por score (menor = mais similar)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score < results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Também adicionar palavras que começam com as mesmas letras
	prefix := input
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}
	for _, word := range wordlist {
		if strings.HasPrefix(word, prefix) {
			found := false
			for _, r := range results {
				if r.word == word {
					found = true
					break
				}
			}
			if !found {
				results = append(results, wordScore{word, 4})
			}
		}
	}

	// Limitar resultados
	var output []string
	for i := 0; i < len(results) && i < maxResults; i++ {
		output = append(output, results[i].word)
	}

	return output
}

// levenshteinDistance calcula a distância de edição entre duas strings
func levenshteinDistance(a, b string) int {
	la := len(a)
	lb := len(b)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Criar matriz
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
	}

	for i := 0; i <= la; i++ {
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			d[i][j] = min3(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}

	return d[la][lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
