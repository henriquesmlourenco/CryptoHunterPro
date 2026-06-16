package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ============================================================================
// FUNCOES DE INPUT
// ============================================================================

func initScanner() {
	globalScanner = bufio.NewScanner(os.Stdin)
	globalScanner.Buffer(make([]byte, 1024*1024), 1024*1024)
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	if globalScanner.Scan() {
		return strings.TrimSpace(globalScanner.Text())
	}
	return ""
}

// ============================================================================
// TEXTOS BILINGUES
// ============================================================================

func t(pt, en string) string {
	if uiLanguage == "pt" {
		return pt
	}
	return en
}

// ============================================================================
// LOGO ASCII - CRYPTO HUNTER PRO
// ============================================================================

func showLogo() {
	fmt.Println()
	fmt.Println("    ================================================================")
	fmt.Println("    |                                                              |")
	fmt.Println("    |     ██████╗██████╗ ██╗   ██╗██████╗ ████████╗ ██████╗        |")
	fmt.Println("    |    ██╔════╝██╔══██╗╚██╗ ██╔╝██╔══██╗╚══██╔══╝██╔═══██╗       |")
	fmt.Println("    |    ██║     ██████╔╝ ╚████╔╝ ██████╔╝   ██║   ██║   ██║       |")
	fmt.Println("    |    ██║     ██╔══██╗  ╚██╔╝  ██╔═══╝    ██║   ██║   ██║       |")
	fmt.Println("    |    ╚██████╗██║  ██║   ██║   ██║        ██║   ╚██████╔╝       |")
	fmt.Println("    |     ╚═════╝╚═╝  ╚═╝   ╚═╝   ╚═╝        ╚═╝    ╚═════╝        |")
	fmt.Println("    |                                                              |")
	fmt.Println("    |    ██╗  ██╗██╗   ██╗███╗   ██╗████████╗███████╗██████╗       |")
	fmt.Println("    |    ██║  ██║██║   ██║████╗  ██║╚══██╔══╝██╔════╝██╔══██╗      |")
	fmt.Println("    |    ███████║██║   ██║██╔██╗ ██║   ██║   █████╗  ██████╔╝      |")
	fmt.Println("    |    ██╔══██║██║   ██║██║╚██╗██║   ██║   ██╔══╝  ██╔══██╗      |")
	fmt.Println("    |    ██║  ██║╚██████╔╝██║ ╚████║   ██║   ███████╗██║  ██║      |")
	fmt.Println("    |    ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝      |")
	fmt.Println("    |                                                              |")
	fmt.Println("    |                  ██████╗ ██████╗  ██████╗                    |")
	fmt.Println("    |                  ██╔══██╗██╔══██╗██╔═══██╗                   |")
	fmt.Println("    |                  ██████╔╝██████╔╝██║   ██║                   |")
	fmt.Println("    |                  ██╔═══╝ ██╔══██╗██║   ██║                   |")
	fmt.Println("    |                  ██║     ██║  ██║╚██████╔╝                   |")
	fmt.Println("    |                  ╚═╝     ╚═╝  ╚═╝ ╚═════╝                    |")
	fmt.Println("    |                                                              |")
        fmt.Println("    |          Crypto Intelligence Engine (CIE)                    |")
        fmt.Println("    |                                                              |")
        fmt.Println("    ================================================================")
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
    fmt.Println("  ================================================================")
}

// ============================================================================
// TELA INICIAL
// ============================================================================

func showWelcomeHeader() {
	showLogo()
}

func chooseUILanguage() {
	showWelcomeHeader()
	for {
		fmt.Println("Choose language / Escolha o idioma:")
		fmt.Println("1. English")
		fmt.Println("2. Portugues")
		fmt.Println()

		choice := getUserInput("> ")
		switch choice {
		case "1":
			uiLanguage = "en"
			return
		case "2":
			uiLanguage = "pt"
			return
		default:
			fmt.Println()
			fmt.Println("  [!] Invalid option. Please enter 1 or 2.")
			fmt.Println("  [!] Opcao invalida. Digite 1 ou 2.")
			fmt.Println()
		}
	}
}

// ============================================================================
// INFORMACOES DO PROGRAMA
// ============================================================================

func showProgramInfo() {
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("CRYPTO HUNTER PRO - CIE", "CRYPTO HUNTER PRO - CIE"))
	fmt.Println("  " + t("Crypto Intelligence Engine", "Crypto Intelligence Engine"))
	fmt.Println("================================================================")
	fmt.Println()

	if uiLanguage == "pt" {
		fmt.Println("  Este modulo realiza a busca de saldos e historico de transacoes")
		fmt.Println("  em multiplas redes blockchain a partir de seed phrases BIP39.")
		fmt.Println()
		fmt.Println("  FUNCIONALIDADES:")
		fmt.Println("  - Busca de saldo nativo e tokens em 18 redes EVM")
		fmt.Println("    (ETH, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Sonic,")
		fmt.Println("     Mantle, Flare, Linea, Scroll, Gnosis, zkSync, Blast, Cronos,")
		fmt.Println("     Celo, Berachain)")
		fmt.Println("  - Descoberta AUTOMATICA de tokens ERC-20 via Blockscout")
		fmt.Println("    (encontra qualquer token, inclusive memecoins, sem configuracao)")
		fmt.Println("  - Busca de saldo em BTC (Legacy, SegWit, Native SegWit, Taproot)")
		fmt.Println("  - Busca de saldo em BCH (CashAddr e Legacy)")
		fmt.Println("  - Busca de saldo e tokens TRC-20 em Tron (TRX)")
		fmt.Println("  - Busca de saldo e tokens SPL em Solana (SOL)")
		fmt.Println("  - Busca de saldo em XRP (Ripple)")
		fmt.Println("  - Busca de saldo e tokens em Stellar (XLM)")
		fmt.Println("  - Busca de saldo e ASA tokens em Algorand (ALGO)")
		fmt.Println("  - Busca de saldo e tokens em Sui (SUI)")
		fmt.Println("  - Busca de saldo em Near (NEAR)")
		fmt.Println("  - Busca de saldo em TON (Toncoin)")
		fmt.Println("  - Busca de saldo em Zcash (ZEC - transparente)")
		fmt.Println("  - Busca de saldo em Litecoin (LTC) e Dogecoin (DOGE)")
		fmt.Println("  - Importacao de arquivos Excel do modulo Unmixer Seed")
		fmt.Println("  - Selecao de derivation paths especificos por rede")
		fmt.Println("  - Configuracao de chaves de API para maior velocidade")
		fmt.Println("  - Relatorio Excel com 3 abas (Resumo, Com Saldo, Com Historico)")
		fmt.Println()
		fmt.Println("  NOTA IMPORTANTE: Este modulo aceita apenas seed phrases completas.")
		fmt.Println("  Para recuperar seeds com palavras faltantes ou embaralhadas, use")
		fmt.Println("  o modulo Crypto Hunter Pro - Unmixer Seed.")
		fmt.Println()
		fmt.Println("  ============================================================")
		fmt.Println("  REDES SUPORTADAS E DERIVACOES (28 blockchains):")
		fmt.Println("  ============================================================")
		fmt.Println()
		fmt.Println("  REDE                    DERIVACAO (PATH)         ENDERECO")
		fmt.Println("  ----------------------  -----------------------  --------")
		fmt.Println("  BTC (Legacy)            m/44'/0'/0'/0/x          1...")
		fmt.Println("  BTC (SegWit)            m/49'/0'/0'/0/x          3...")
		fmt.Println("  BTC (Native SegWit)     m/84'/0'/0'/0/x          bc1q...")
		fmt.Println("  BTC (Taproot)           m/86'/0'/0'/0/x          bc1p...")
		fmt.Println("  BCH (CashAddr)          m/44'/145'/0'/0/x        bitcoincash:q...")
		fmt.Println("  BCH (Legacy)            m/44'/145'/0'/0/x        1...")
		fmt.Println("    (mesmo endereco, formatos diferentes - contabilizado 1x)")
		fmt.Println("  18 redes EVM            m/44'/60'/0'/0/x         0x...")
		fmt.Println("  LTC (Legacy)            m/44'/2'/0'/0/x          L...")
		fmt.Println("  LTC (SegWit)            m/49'/2'/0'/0/x          3... / M...")
		fmt.Println("  LTC (Native SegWit)     m/84'/2'/0'/0/x          ltc1q...")
		fmt.Println("  Dogecoin (DOGE)         m/44'/3'/0'/0/x          D...")
		fmt.Println("  Tron (TRX)              m/44'/195'/0'/0/x        T...")
		fmt.Println("  Solana (SOL)            m/44'/501'/0'/0'         Base58...")
		fmt.Println("  XRP (Ripple)            m/44'/144'/0'/0/x        r...")
		fmt.Println("  Stellar (XLM)           m/44'/148'/0'            G...")
		fmt.Println("  Algorand (ALGO)         m/44'/283'/0'/0/x        Base32...")
		fmt.Println("  Sui (SUI)               m/44'/784'/0'/0'/0'      0x...")
		fmt.Println("  Near (NEAR)             m/44'/397'/0'            hex...")
		fmt.Println("  TON (Toncoin)           m/44'/607'/0'            UQ...")
		fmt.Println("  Zcash (ZEC)             m/44'/133'/0'/0/x        t1...")
		fmt.Println()
		fmt.Println("  ============================================================")
		fmt.Println("  REDES EVM (18 redes - mesmo endereco 0x):")
		fmt.Println("  ============================================================")
		fmt.Println()
		fmt.Println("  NOTA: Em 9 redes EVM (Ethereum, Polygon, Optimism, Base, Gnosis,")
		fmt.Println("  Scroll, Linea, zkSync, Celo) o programa busca automaticamente")
		fmt.Println("  QUALQUER token ERC-20 via Blockscout, incluindo memecoins e")
		fmt.Println("  tokens novos. Nas demais redes EVM, sao verificados os tokens")
		fmt.Println("  listados abaixo.")
		fmt.Println()
		fmt.Println("  Ethereum (ETH + 24 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDT, USDC, DAI, BUSD, TUSD, FRAX, WBTC, WETH, LINK, UNI,")
		fmt.Println("            AAVE, MKR, COMP, SNX, CRV, LDO, SHIB, PEPE, POL, APE,")
		fmt.Println("            SAND, MANA, GRT, ENS")
		fmt.Println("    + Descoberta automatica via Blockscout (qualquer token ERC-20)")
		fmt.Println()
		fmt.Println("  BNB Smart Chain / BSC (BNB + 14 tokens):")
		fmt.Println("    Nativo: BNB")
		fmt.Println("    Tokens: USDT, USDC, BUSD, DAI, CAKE, WBNB, XVS, BTCB, ETH, LINK,")
		fmt.Println("            DOT, ADA, DOGE, FLOKI")
		fmt.Println()
		fmt.Println("  Polygon (POL + 13 tokens):")
		fmt.Println("    Nativo: POL")
		fmt.Println("    Tokens: USDT, USDC, USDC.e, DAI, WPOL, AAVE, LINK, WBTC, WETH,")
		fmt.Println("            CRV, BAL, GNS, QUICK")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  Arbitrum (ETH + 12 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDT, USDC, USDC.e, DAI, WETH, ARB, LINK, GRT, GMX, MAGIC,")
		fmt.Println("            RDNT, PENDLE")
		fmt.Println()
		fmt.Println("  Avalanche (AVAX + 11 tokens):")
		fmt.Println("    Nativo: AVAX")
		fmt.Println("    Tokens: USDT, USDC, DAI.e, WAVAX, WETH.e, WBTC.e, AAVE.e, LINK.e,")
		fmt.Println("            JOE, GMX, sAVAX")
		fmt.Println()
		fmt.Println("  Optimism (ETH + 10 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDT, USDC, DAI, OP, WETH, LINK, SNX, VELO, AAVE, PERP")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  Base (ETH + 9 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDC, DAI, WETH, cbETH, AERO, BRETT, DEGEN, TOSHI, USDbC")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  Sonic (S + 6 tokens):")
		fmt.Println("    Nativo: S")
		fmt.Println("    Tokens: USDC.e, WETH, wS, BRUSH, EQUAL, USDT")
		fmt.Println()
		fmt.Println("  Mantle (MNT + 6 tokens):")
		fmt.Println("    Nativo: MNT")
		fmt.Println("    Tokens: USDT, USDC, WMNT, WETH, mETH, PUFF")
		fmt.Println()
		fmt.Println("  Flare (FLR + 5 tokens):")
		fmt.Println("    Nativo: FLR")
		fmt.Println("    Tokens: WFLR, USDC.e, USDT.e, sFLR, WETH")
		fmt.Println()
		fmt.Println("  Linea (ETH + 5 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDC, USDT, WETH, DAI, wstETH")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  Scroll (ETH + 5 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDC, USDT, WETH, wstETH, DAI")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  Gnosis (xDAI + 5 tokens):")
		fmt.Println("    Nativo: xDAI")
		fmt.Println("    Tokens: GNO, USDC, USDT, WETH, sDAI")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  zkSync Era (ETH + 4 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDC, USDT, WETH, ZK")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  Blast (ETH + 3 tokens):")
		fmt.Println("    Nativo: ETH")
		fmt.Println("    Tokens: USDB, WETH, USDT")
		fmt.Println()
		fmt.Println("  Cronos (CRO + 4 tokens):")
		fmt.Println("    Nativo: CRO")
		fmt.Println("    Tokens: USDC, USDT, WETH, WCRO")
		fmt.Println()
		fmt.Println("  Celo (CELO + 5 tokens):")
		fmt.Println("    Nativo: CELO")
		fmt.Println("    Tokens: cUSD, cEUR, USDT, USDC, WETH")
		fmt.Println("    + Descoberta automatica via Blockscout")
		fmt.Println()
		fmt.Println("  Berachain (BERA + 5 tokens):")
		fmt.Println("    Nativo: BERA")
		fmt.Println("    Tokens: HONEY, USDC, USDT, WETH, WBTC")
		fmt.Println()
		fmt.Println("  ============================================================")
		fmt.Println("  REDES NAO-EVM:")
		fmt.Println("  ============================================================")
		fmt.Println()
		fmt.Println("  Solana (SOL + 13 tokens SPL + descoberta automatica):")
		fmt.Println("    Nativo: SOL")
		fmt.Println("    Tokens: USDC, USDT, JUP, BONK, RAY, ORCA, JTO, PYTH, W, RENDER,")
		fmt.Println("            WIF, mSOL, jitoSOL")
		fmt.Println("    + Descoberta automatica via RPC (qualquer token SPL)")
		fmt.Println()
		fmt.Println("  Tron (TRX + 9 tokens TRC-20 + descoberta automatica):")
		fmt.Println("    Nativo: TRX")
		fmt.Println("    Tokens: USDT, USDC, USDD, TUSD, WTRX, BTT, JST, SUN, WIN")
		fmt.Println("    + Descoberta automatica via TronGrid (qualquer token TRC-20)")
		fmt.Println()
		fmt.Println("  XRP / Ripple: Saldo nativo XRP")
		fmt.Println("  Stellar (XLM): Saldo nativo + todos os tokens automaticamente")
		fmt.Println("  Algorand (ALGO): Saldo nativo + ASA tokens automaticamente")
		fmt.Println("  Sui (SUI): Saldo nativo + todos os tokens automaticamente")
		fmt.Println("  Near (NEAR): Saldo nativo")
		fmt.Println("  Bitcoin (BTC): 4 tipos de endereco (Legacy, SegWit, NativeSegWit, Taproot)")
		fmt.Println("  Bitcoin Cash (BCH): CashAddr + Legacy")
		fmt.Println("  Litecoin (LTC): Legacy, SegWit, Native SegWit")
		fmt.Println("  Dogecoin (DOGE): Saldo nativo")
		fmt.Println("  TON (Toncoin): Saldo nativo")
		fmt.Println("  Zcash (ZEC): Enderecos transparentes (t1...)")
		fmt.Println()
		fmt.Println("  NOTA: A seed phrase e universal e identica para todas")
		fmt.Println("  as redes blockchain. O derivation path determina qual rede")
		fmt.Println("  e tipo de endereco sera gerado a partir da mesma seed.")
		fmt.Println("  Voce pode selecionar os paths desejados no proximo menu.")
	} else {
		fmt.Println("  This module searches for balances and transaction history")
		fmt.Println("  across multiple blockchain networks from BIP39 seed phrases.")
		fmt.Println()
		fmt.Println("  FEATURES:")
		fmt.Println("  - Native balance and token search on 18 EVM networks")
		fmt.Println("    (ETH, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Sonic,")
		fmt.Println("     Mantle, Flare, Linea, Scroll, Gnosis, zkSync, Blast, Cronos,")
		fmt.Println("     Celo, Berachain)")
		fmt.Println("  - AUTOMATIC ERC-20 token discovery via Blockscout")
		fmt.Println("    (finds any token, including memecoins, without configuration)")
		fmt.Println("  - BTC balance search (Legacy, SegWit, Native SegWit, Taproot)")
		fmt.Println("  - BCH balance search (CashAddr and Legacy)")
		fmt.Println("  - TRX balance and TRC-20 token search on Tron")
		fmt.Println("  - SOL balance and SPL token search on Solana")
		fmt.Println("  - XRP (Ripple) balance search")
		fmt.Println("  - Stellar (XLM) balance and token search")
		fmt.Println("  - Algorand (ALGO) balance and ASA token search")
		fmt.Println("  - Sui (SUI) balance and token search")
		fmt.Println("  - Near (NEAR) balance search")
		fmt.Println("  - TON balance search (Toncoin)")
		fmt.Println("  - ZEC balance search (Zcash transparent)")
		fmt.Println("  - Litecoin (LTC) and Dogecoin (DOGE) balance search")
		fmt.Println("  - Import Excel files from Unmixer Seed module")
		fmt.Println("  - Specific derivation path selection per network")
		fmt.Println("  - API key configuration for faster scanning")
		fmt.Println("  - Excel report with 3 tabs (Summary, With Balance, With History)")
		fmt.Println()
		fmt.Println("  IMPORTANT NOTE: This module only accepts complete seed phrases.")
		fmt.Println("  To recover seeds with missing or shuffled words, use the")
		fmt.Println("  Crypto Hunter Pro - Unmixer Seed module.")
		fmt.Println()
		fmt.Println("  ============================================================")
		fmt.Println("  SUPPORTED NETWORKS AND DERIVATIONS (28 blockchains):")
		fmt.Println("  ============================================================")
		fmt.Println()
		fmt.Println("  NETWORK                 DERIVATION (PATH)        ADDRESS")
		fmt.Println("  ----------------------  -----------------------  --------")
		fmt.Println("  BTC (Legacy)            m/44'/0'/0'/0/x          1...")
		fmt.Println("  BTC (SegWit)            m/49'/0'/0'/0/x          3...")
		fmt.Println("  BTC (Native SegWit)     m/84'/0'/0'/0/x          bc1q...")
		fmt.Println("  BTC (Taproot)           m/86'/0'/0'/0/x          bc1p...")
		fmt.Println("  BCH (CashAddr)          m/44'/145'/0'/0/x        bitcoincash:q...")
		fmt.Println("  BCH (Legacy)            m/44'/145'/0'/0/x        1...")
		fmt.Println("    (same address, different formats - counted 1x)")
		fmt.Println("  18 EVM networks         m/44'/60'/0'/0/x         0x...")
		fmt.Println("  LTC (Legacy)            m/44'/2'/0'/0/x          L...")
		fmt.Println("  LTC (SegWit)            m/49'/2'/0'/0/x          3... / M...")
		fmt.Println("  LTC (Native SegWit)     m/84'/2'/0'/0/x          ltc1q...")
		fmt.Println("  Dogecoin (DOGE)         m/44'/3'/0'/0/x          D...")
		fmt.Println("  Tron (TRX)              m/44'/195'/0'/0/x        T...")
		fmt.Println("  Solana (SOL)            m/44'/501'/0'/0'         Base58...")
		fmt.Println("  XRP (Ripple)            m/44'/144'/0'/0/x        r...")
		fmt.Println("  Stellar (XLM)           m/44'/148'/0'            G...")
		fmt.Println("  Algorand (ALGO)         m/44'/283'/0'/0/x        Base32...")
		fmt.Println("  Sui (SUI)               m/44'/784'/0'/0'/0'      0x...")
		fmt.Println("  Near (NEAR)             m/44'/397'/0'            hex...")
		fmt.Println("  TON (Toncoin)           m/44'/607'/0'            UQ...")
		fmt.Println("  Zcash (ZEC)             m/44'/133'/0'/0/x        t1...")
		fmt.Println()
		fmt.Println("  ============================================================")
		fmt.Println("  EVM NETWORKS (18 networks - same 0x address):")
		fmt.Println("  ============================================================")
		fmt.Println()
		fmt.Println("  NOTE: On 9 EVM networks (Ethereum, Polygon, Optimism, Base, Gnosis,")
		fmt.Println("  Scroll, Linea, zkSync, Celo) the program automatically searches")
		fmt.Println("  for ANY ERC-20 token via Blockscout, including memecoins and")
		fmt.Println("  new tokens. On other EVM networks, the tokens listed below")
		fmt.Println("  are checked.")
		fmt.Println()
		fmt.Println("  Ethereum (ETH + 24 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDT, USDC, DAI, BUSD, TUSD, FRAX, WBTC, WETH, LINK, UNI,")
		fmt.Println("            AAVE, MKR, COMP, SNX, CRV, LDO, SHIB, PEPE, POL, APE,")
		fmt.Println("            SAND, MANA, GRT, ENS")
		fmt.Println("    + Automatic discovery via Blockscout (any ERC-20 token)")
		fmt.Println()
		fmt.Println("  BNB Smart Chain / BSC (BNB + 14 tokens):")
		fmt.Println("    Native: BNB")
		fmt.Println("    Tokens: USDT, USDC, BUSD, DAI, CAKE, WBNB, XVS, BTCB, ETH, LINK,")
		fmt.Println("            DOT, ADA, DOGE, FLOKI")
		fmt.Println()
		fmt.Println("  Polygon (POL + 13 tokens):")
		fmt.Println("    Native: POL")
		fmt.Println("    Tokens: USDT, USDC, USDC.e, DAI, WPOL, AAVE, LINK, WBTC, WETH,")
		fmt.Println("            CRV, BAL, GNS, QUICK")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  Arbitrum (ETH + 12 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDT, USDC, USDC.e, DAI, WETH, ARB, LINK, GRT, GMX, MAGIC,")
		fmt.Println("            RDNT, PENDLE")
		fmt.Println()
		fmt.Println("  Avalanche (AVAX + 11 tokens):")
		fmt.Println("    Native: AVAX")
		fmt.Println("    Tokens: USDT, USDC, DAI.e, WAVAX, WETH.e, WBTC.e, AAVE.e, LINK.e,")
		fmt.Println("            JOE, GMX, sAVAX")
		fmt.Println()
		fmt.Println("  Optimism (ETH + 10 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDT, USDC, DAI, OP, WETH, LINK, SNX, VELO, AAVE, PERP")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  Base (ETH + 9 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDC, DAI, WETH, cbETH, AERO, BRETT, DEGEN, TOSHI, USDbC")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  Sonic (S + 6 tokens):")
		fmt.Println("    Native: S")
		fmt.Println("    Tokens: USDC.e, WETH, wS, BRUSH, EQUAL, USDT")
		fmt.Println()
		fmt.Println("  Mantle (MNT + 6 tokens):")
		fmt.Println("    Native: MNT")
		fmt.Println("    Tokens: USDT, USDC, WMNT, WETH, mETH, PUFF")
		fmt.Println()
		fmt.Println("  Flare (FLR + 5 tokens):")
		fmt.Println("    Native: FLR")
		fmt.Println("    Tokens: WFLR, USDC.e, USDT.e, sFLR, WETH")
		fmt.Println()
		fmt.Println("  Linea (ETH + 5 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDC, USDT, WETH, DAI, wstETH")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  Scroll (ETH + 5 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDC, USDT, WETH, wstETH, DAI")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  Gnosis (xDAI + 5 tokens):")
		fmt.Println("    Native: xDAI")
		fmt.Println("    Tokens: GNO, USDC, USDT, WETH, sDAI")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  zkSync Era (ETH + 4 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDC, USDT, WETH, ZK")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  Blast (ETH + 3 tokens):")
		fmt.Println("    Native: ETH")
		fmt.Println("    Tokens: USDB, WETH, USDT")
		fmt.Println()
		fmt.Println("  Cronos (CRO + 4 tokens):")
		fmt.Println("    Native: CRO")
		fmt.Println("    Tokens: USDC, USDT, WETH, WCRO")
		fmt.Println()
		fmt.Println("  Celo (CELO + 5 tokens):")
		fmt.Println("    Native: CELO")
		fmt.Println("    Tokens: cUSD, cEUR, USDT, USDC, WETH")
		fmt.Println("    + Automatic discovery via Blockscout")
		fmt.Println()
		fmt.Println("  Berachain (BERA + 5 tokens):")
		fmt.Println("    Native: BERA")
		fmt.Println("    Tokens: HONEY, USDC, USDT, WETH, WBTC")
		fmt.Println()
		fmt.Println("  ============================================================")
		fmt.Println("  NON-EVM NETWORKS:")
		fmt.Println("  ============================================================")
		fmt.Println()
		fmt.Println("  Solana (SOL + 13 SPL tokens + automatic discovery):")
		fmt.Println("    Native: SOL")
		fmt.Println("    Tokens: USDC, USDT, JUP, BONK, RAY, ORCA, JTO, PYTH, W, RENDER,")
		fmt.Println("            WIF, mSOL, jitoSOL")
		fmt.Println("    + Automatic discovery via RPC (any SPL token)")
		fmt.Println()
		fmt.Println("  Tron (TRX + 9 TRC-20 tokens + automatic discovery):")
		fmt.Println("    Native: TRX")
		fmt.Println("    Tokens: USDT, USDC, USDD, TUSD, WTRX, BTT, JST, SUN, WIN")
		fmt.Println("    + Automatic discovery via TronGrid (any TRC-20 token)")
		fmt.Println()
		fmt.Println("  XRP / Ripple: Native XRP balance")
		fmt.Println("  Stellar (XLM): Native balance + all tokens automatically")
		fmt.Println("  Algorand (ALGO): Native balance + ASA tokens automatically")
		fmt.Println("  Sui (SUI): Native balance + all tokens automatically")
		fmt.Println("  Near (NEAR): Native balance")
		fmt.Println("  Bitcoin (BTC): 4 address types (Legacy, SegWit, NativeSegWit, Taproot)")
		fmt.Println("  Bitcoin Cash (BCH): CashAddr + Legacy")
		fmt.Println("  Litecoin (LTC): Legacy, SegWit, Native SegWit")
		fmt.Println("  Dogecoin (DOGE): Native balance")
		fmt.Println("  TON (Toncoin): Native balance")
		fmt.Println("  Zcash (ZEC): Transparent addresses (t1...)")
		fmt.Println()
		fmt.Println("  NOTE: The seed phrase is universal and identical for all")
		fmt.Println("  blockchain networks. The derivation path determines which network")
		fmt.Println("  and address type will be generated from the same seed.")
		fmt.Println("  You can select the desired paths in the next menu.")
	}

	fmt.Println()
	fmt.Println(t("  Pressione Enter para continuar...", "  Press Enter to continue..."))
	getUserInput("")
}

// ============================================================================
// MENU PRINCIPAL
// ============================================================================

func showMainMenu() int {
	for {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("  " + t("MENU PRINCIPAL", "MAIN MENU"))
		fmt.Println("================================================================")
		fmt.Println()
		fmt.Println(t(
			"  1. Digitar seed phrase manualmente",
			"  1. Enter seed phrase manually"))
		fmt.Println(t(
			"  2. Importar arquivos Excel do Unmixer Seed",
			"  2. Import Excel files from Unmixer Seed"))
		fmt.Println(t(
			"  0. Sair",
			"  0. Exit"))
		fmt.Println()

		choice := getUserInput("> ")
		switch choice {
		case "0":
			return 0
		case "1":
			return 1
		case "2":
			return 2
		default:
			fmt.Println()
			fmt.Println(t(
				"  [!] Opcao invalida. Digite 0, 1 ou 2.",
				"  [!] Invalid option. Enter 0, 1 or 2."))
		}
	}
}

// ============================================================================
// VALIDACAO DE SEED
// ============================================================================

func chooseSeedValidation() bool {
	for {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("  " + t("VALIDACAO DA SEED", "SEED VALIDATION"))
		fmt.Println("================================================================")
		fmt.Println()
		fmt.Println(t(
			"  [1] BIP39 (Padrao)",
			"  [1] BIP39 (Standard)"))
		fmt.Println(t(
			"      - Validacao por checksum (padrao da industria)",
			"      - Checksum validation (industry standard)"))
		fmt.Println(t(
			"      - Compativel com a maioria das carteiras",
			"      - Compatible with most wallets"))
		fmt.Println()
		fmt.Println(t(
			"  [2] Ignorar checksum (Forca Bruta)",
			"  [2] Skip checksum (Brute Force)"))
		fmt.Println(t(
			"      - Ignora validacao de checksum",
			"      - Ignores checksum validation"))
		fmt.Println(t(
			"      - Util para seeds de carteiras nao-padrao",
			"      - Useful for non-standard wallet seeds"))
		fmt.Println(t(
			"      - Util para seeds importadas no Electron Cash",
			"      - Useful for seeds imported in Electron Cash"))
		fmt.Println()

		choice := getUserInput(t("  Escolha (1/2): ", "  Choose (1/2): "))
		switch choice {
		case "1":
			return false
		case "2":
			return true
		default:
			fmt.Println(t(
				"  [!] Opcao invalida. Digite 1 ou 2.",
				"  [!] Invalid option. Enter 1 or 2."))
		}
	}
}

// ============================================================================
// INPUT DE SEED MANUAL
// ============================================================================

func getSeedManual(skipChecksum bool) string {
	for {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("  " + t("DIGITE A SEED PHRASE", "ENTER THE SEED PHRASE"))
		fmt.Println("================================================================")
		fmt.Println()
		fmt.Println(t(
			"  Digite ou cole as palavras da seed phrase (separadas por espaco):",
			"  Type or paste the seed phrase words (separated by space):"))
		fmt.Println(t(
			"  (digite 'voltar' para retornar ao menu)",
			"  (type 'back' to return to menu)"))
		fmt.Println()

		seed := getUserInput("  > ")
		seed = strings.TrimSpace(seed)

		if strings.ToLower(seed) == "voltar" || strings.ToLower(seed) == "back" {
			return ""
		}

		words := strings.Fields(seed)
		seed = strings.Join(words, " ")

		wordCount := len(words)

		validCounts := map[int]bool{12: true, 15: true, 18: true, 21: true, 24: true}
		if !validCounts[wordCount] {
			fmt.Println()
			fmt.Printf(t(
				"  [!] ERRO: Seed com %d palavras. Deve ter 12, 15, 18, 21 ou 24 palavras.\n",
				"  [!] ERROR: Seed with %d words. Must have 12, 15, 18, 21 or 24 words.\n"), wordCount)
			fmt.Println(t(
				"  Tente novamente.",
				"  Try again."))
			continue
		}

		if !skipChecksum {
			if !ValidateSeedPhrase(seed) {
				fmt.Println()
				fmt.Println(t(
					"  [!] ERRO: Seed invalida (checksum BIP39 incorreto).",
					"  [!] ERROR: Invalid seed (incorrect BIP39 checksum)."))
				fmt.Println(t(
					"  Verifique as palavras ou use a opcao 'Ignorar checksum'.",
					"  Check the words or use the 'Skip checksum' option."))
				fmt.Println(t(
					"  Tente novamente.",
					"  Try again."))
				continue
			}
			fmt.Println()
			fmt.Printf(t(
				"  [OK] Seed BIP39 valida com %d palavras.\n",
				"  [OK] Valid BIP39 seed with %d words.\n"), wordCount)
		} else {
			fmt.Println()
			fmt.Printf(t(
				"  [OK] Seed com %d palavras (checksum ignorado).\n",
				"  [OK] Seed with %d words (checksum skipped).\n"), wordCount)
		}

		return seed
	}
}

// ============================================================================
// PASSPHRASE (OPCIONAL)
// ============================================================================

func askPassphrase() string {
	for {
		fmt.Println()
		fmt.Println(t(
			"  Deseja adicionar uma passphrase (25a palavra)?",
			"  Do you want to add a passphrase (25th word)?"))
		fmt.Println(t(
			"  [1] Nao - Continuar sem passphrase",
			"  [1] No - Continue without passphrase"))
		fmt.Println(t(
			"  [2] Sim - Adicionar uma passphrase extra",
			"  [2] Yes - Add an extra passphrase"))
		fmt.Println()

		choice := getUserInput("  > ")
		switch choice {
		case "1":
			return ""
		case "2":
			fmt.Println()
			passphrase := getUserInput(t("  Digite a passphrase: ", "  Enter the passphrase: "))
			return passphrase
		default:
			fmt.Println(t(
				"  [!] Opcao invalida. Digite 1 ou 2.",
				"  [!] Invalid option. Enter 1 or 2."))
		}
	}
}

// ============================================================================
// SELECAO DE REDES
// ============================================================================

func selectNetworks() []NetworkGroup {
	networks := getDefaultNetworkGroups()

	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("SELECAO DE REDES", "NETWORK SELECTION"))
	fmt.Println("================================================================")
	fmt.Println()
	fmt.Println(t(
		"  Selecione as redes para buscar (digite os numeros separados por virgula):",
		"  Select networks to search (enter numbers separated by comma):"))
	fmt.Println()

	for i, ng := range networks {
		name := ng.NameEN
		if uiLanguage == "pt" {
			name = ng.NamePT
		}
		fmt.Printf("  %d. %s\n", i+1, name)
	}
	fmt.Println()
	fmt.Printf("  %d. %s\n", len(networks)+1, t("TODAS as redes", "ALL networks"))
	fmt.Println()

	choice := getUserInput("  > ")

	allChoice := strconv.Itoa(len(networks) + 1)
	if strings.TrimSpace(choice) == allChoice {
		for i := range networks {
			networks[i].Enabled = true
		}
		return networks
	}

	parts := strings.Split(choice, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		idx, err := strconv.Atoi(p)
		if err == nil && idx >= 1 && idx <= len(networks) {
			networks[idx-1].Enabled = true
		}
	}

	anySelected := false
	for _, ng := range networks {
		if ng.Enabled {
			anySelected = true
			break
		}
	}
	if !anySelected {
		fmt.Println()
		fmt.Println(t(
			"  [!] Nenhuma rede selecionada. Selecionando EVM por padrao.",
			"  [!] No network selected. Selecting EVM by default."))
		networks[0].Enabled = true
	}

	return networks
}

// ============================================================================
// SELECAO DE DERIVACOES POR REDE
// ============================================================================

func selectDerivations(networks []NetworkGroup) []NetworkGroup {
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("SELECAO DE DERIVACOES", "DERIVATION SELECTION"))
	fmt.Println("================================================================")

	for i, ng := range networks {
		if !ng.Enabled {
			continue
		}

		if len(ng.Derivations) == 1 {
			networks[i].Derivations[0].Enabled = true
			continue
		}

		name := ng.NameEN
		if uiLanguage == "pt" {
			name = ng.NamePT
		}

		fmt.Println()
		fmt.Printf("  --- %s ---\n", name)
		fmt.Println(t(
			"  Selecione as derivacoes (numeros separados por virgula):",
			"  Select derivations (numbers separated by comma):"))
		fmt.Println()

		for j, dp := range ng.Derivations {
			dpName := dp.NameEN
			if uiLanguage == "pt" {
				dpName = dp.NamePT
			}
			fmt.Printf("    %d. %s  [%s]\n", j+1, dpName, dp.Path)
		}
		fmt.Println()
		fmt.Printf("    %d. %s\n", len(ng.Derivations)+1, t("TODAS", "ALL"))
		fmt.Println()

		choice := getUserInput("    > ")

		allChoice := strconv.Itoa(len(ng.Derivations) + 1)
		if strings.TrimSpace(choice) == allChoice {
			for j := range networks[i].Derivations {
				networks[i].Derivations[j].Enabled = true
			}
		} else {
			for j := range networks[i].Derivations {
				networks[i].Derivations[j].Enabled = false
			}
			parts := strings.Split(choice, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				idx, err := strconv.Atoi(p)
				if err == nil && idx >= 1 && idx <= len(ng.Derivations) {
					networks[i].Derivations[idx-1].Enabled = true
				}
			}
		}

		anySelected := false
		for _, dp := range networks[i].Derivations {
			if dp.Enabled {
				anySelected = true
				break
			}
		}
		if !anySelected {
			fmt.Println(t(
				"    [!] Nenhuma derivacao selecionada. Selecionando a primeira por padrao.",
				"    [!] No derivation selected. Selecting the first one by default."))
			networks[i].Derivations[0].Enabled = true
		}
	}

	return networks
}

// ============================================================================
// SELECAO DE RANGE DE INDICES
// ============================================================================

func selectIndexRange() (int, int) {
	for {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("  " + t("RANGE DE INDICES", "INDEX RANGE"))
		fmt.Println("================================================================")
		fmt.Println()
		fmt.Println(t(
			"  Selecione o range de indices de endereco para buscar:",
			"  Select the address index range to search:"))
		fmt.Println()
		fmt.Println(t(
			"  [1] 0 a 20  (padrao - rapido)",
			"  [1] 0 to 20  (default - fast)"))
		fmt.Println(t(
			"  [2] 0 a 50  (abrangente)",
			"  [2] 0 to 50  (comprehensive)"))
		fmt.Println(t(
			"  [3] 0 a 100 (extensivo)",
			"  [3] 0 to 100 (extensive)"))
		fmt.Println(t(
			"  [4] Personalizado",
			"  [4] Custom"))
		fmt.Println()

		choice := getUserInput("  > ")

		switch choice {
		case "1":
			return 0, 20
		case "2":
			return 0, 50
		case "3":
			return 0, 100
		case "4":
			for {
				fmt.Println()
				startStr := getUserInput(t("  Indice inicial: ", "  Start index: "))
				endStr := getUserInput(t("  Indice final: ", "  End index: "))
				start, _ := strconv.Atoi(startStr)
				end, _ := strconv.Atoi(endStr)
				if end >= start {
					return start, end
				}
				fmt.Println(t(
					"  [!] Range invalido! O indice final deve ser maior ou igual ao inicial.",
					"  [!] Invalid range! End index must be greater than or equal to start index."))
			}
		default:
			fmt.Println(t(
				"  [!] Opcao invalida. Digite 1, 2, 3 ou 4.",
				"  [!] Invalid option. Enter 1, 2, 3 or 4."))
		}
	}
}

// ============================================================================
// CONFIGURACAO DE API KEYS (MENU ANTES DO MENU PRINCIPAL)
// ============================================================================

func showAPIConfigMenu() {
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("CONFIGURACAO DE APIs", "API CONFIGURATION"))
	fmt.Println("================================================================")
	fmt.Println()
	fmt.Println(t(
		"  O programa utiliza APIs publicas gratuitas por padrao.",
		"  The program uses free public APIs by default."))
	fmt.Println(t(
		"  Para buscas em larga escala (milhares de seeds por dias seguidos),",
		"  For large-scale searches (thousands of seeds over days),"))
	fmt.Println(t(
		"  voce pode configurar chaves de API pagas para maior velocidade",
		"  you can configure paid API keys for higher speed"))
	fmt.Println(t(
		"  e sem interrupcoes por rate-limit.",
		"  and no rate-limit interruptions."))
	fmt.Println()
	fmt.Println(t(
		"  Caso a API paga falhe, o programa continua automaticamente",
		"  If the paid API fails, the program automatically continues"))
	fmt.Println(t(
		"  usando as APIs publicas sem interrupcao.",
		"  using public APIs without interruption."))
	fmt.Println()

	for {
		fmt.Println(t(
			"  [1] Continuar com APIs publicas gratuitas (padrao)",
			"  [1] Continue with free public APIs (default)"))
		fmt.Println(t(
			"  [2] Configurar chaves de API (pagas)",
			"  [2] Configure API keys (paid)"))
		fmt.Println()

		choice := getUserInput("> ")
		switch choice {
		case "1":
			fmt.Println()
			fmt.Println(t(
				"  [OK] Usando APIs publicas gratuitas.",
				"  [OK] Using free public APIs."))
			return
		case "2":
			configureAPIKeys()
			return
		default:
			fmt.Println()
			fmt.Println(t(
				"  [!] Opcao invalida. Digite 1 ou 2.",
				"  [!] Invalid option. Enter 1 or 2."))
			fmt.Println()
		}
	}
}

func configureAPIKeys() {
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("CONFIGURACAO DE CHAVES DE API (PAGAS)", "API KEY CONFIGURATION (PAID)"))
	fmt.Println("================================================================")
	fmt.Println()
	fmt.Println(t(
		"  Configure as APIs desejadas. Pressione Enter para pular.",
		"  Configure desired APIs. Press Enter to skip."))
	fmt.Println(t(
		"  Voce pode configurar 1 ou as 2 APIs abaixo.",
		"  You can configure 1 or both APIs below."))
	fmt.Println()
	fmt.Println(t(
		"  IMPORTANTE: Se a API paga falhar durante o escaneamento",
		"  IMPORTANT: If the paid API fails during scanning"))
	fmt.Println(t(
		"  (creditos acabaram, timeout, erro de conexao, etc),",
		"  (credits ran out, timeout, connection error, etc),"))
	fmt.Println(t(
		"  o programa detecta automaticamente e continua a busca",
		"  the program automatically detects it and continues"))
	fmt.Println(t(
		"  usando as APIs publicas gratuitas, sem interrupcao.",
		"  using the free public APIs, without interruption."))
	fmt.Println(t(
		"  Voce nao perde nenhum resultado.",
		"  You don't lose any results."))
	fmt.Println()

	// --- ALCHEMY ---
	fmt.Println("  ============================================================")
	fmt.Println(t("  1. ALCHEMY (Recomendada)", "  1. ALCHEMY (Recommended)"))
	fmt.Println("  ============================================================")
	fmt.Println(t(
		"     Cobre 16 redes EVM: Ethereum, BSC, Polygon, Arbitrum,",
		"     Covers 16 EVM networks: Ethereum, BSC, Polygon, Arbitrum,"))
	fmt.Println(t(
		"     Avalanche, Optimism, Base, Linea, Scroll, Gnosis, zkSync,",
		"     Avalanche, Optimism, Base, Linea, Scroll, Gnosis, zkSync,"))
	fmt.Println(t(
		"     Blast, Celo, Berachain, Sonic, Mantle",
		"     Blast, Celo, Berachain, Sonic, Mantle"))
	fmt.Println(t(
		"     + 1 rede nao-EVM: Solana",
		"     + 1 non-EVM network: Solana"))
	fmt.Println()
	fmt.Println(t(
		"     NAO cobre (usam APIs publicas automaticamente):",
		"     NOT covered (use public APIs automatically):"))
	fmt.Println(t(
		"     EVM: Cronos, Flare",
		"     EVM: Cronos, Flare"))
	fmt.Println(t(
		"     Nao-EVM: Bitcoin, Bitcoin Cash, Tron, Litecoin, Dogecoin,",
		"     Non-EVM: Bitcoin, Bitcoin Cash, Tron, Litecoin, Dogecoin,"))
	fmt.Println(t(
		"     TON, Zcash, XRP, Stellar, Algorand, Sui, Near",
		"     TON, Zcash, XRP, Stellar, Algorand, Sui, Near"))
	fmt.Println()
	fmt.Println("     " + t("Obter em:", "Get at:") + " https://www.alchemy.com/pricing")
	fmt.Printf("     %s: %s\n", t("Atual", "Current"), maskKey(apiKeys.AlchemyKey))
	key := getUserInput(t("     Chave (Enter para pular): ", "     Key (Enter to skip): "))
	if key != "" {
		apiKeys.AlchemyKey = key
	}
	fmt.Println()

	// --- TRONGRID ---
	fmt.Println("  ============================================================")
	fmt.Println(t("  2. TRONGRID", "  2. TRONGRID"))
	fmt.Println("  ============================================================")
	fmt.Println(t(
		"     Cobre: Tron (TRX + todos os tokens TRC-20)",
		"     Covers: Tron (TRX + all TRC-20 tokens)"))
	fmt.Println()
	fmt.Println(t(
		"     NAO cobre (usam APIs publicas automaticamente):",
		"     NOT covered (use public APIs automatically):"))
	fmt.Println(t(
		"     Todas as demais 30 redes (18 EVM + Bitcoin, Bitcoin Cash,",
		"     All other 30 networks (18 EVM + Bitcoin, Bitcoin Cash,"))
	fmt.Println(t(
		"     Solana, Litecoin, Dogecoin, TON, Zcash, XRP, Stellar,",
		"     Solana, Litecoin, Dogecoin, TON, Zcash, XRP, Stellar,"))
	fmt.Println(t(
		"     Algorand, Sui, Near)",
		"     Algorand, Sui, Near)"))
	fmt.Println()
	fmt.Println("     " + t("Obter em:", "Get at:") + " https://www.trongrid.io/")
	fmt.Printf("     %s: %s\n", t("Atual", "Current"), maskKey(apiKeys.TrongridKey))
	key = getUserInput(t("     Chave (Enter para pular): ", "     Key (Enter to skip): "))
	if key != "" {
		apiKeys.TrongridKey = key
	}

	// --- RESUMO ---
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("RESUMO DAS APIs CONFIGURADAS", "CONFIGURED APIs SUMMARY"))
	fmt.Println("================================================================")
	fmt.Println()
	if apiKeys.AlchemyKey != "" {
		fmt.Println(t(
			"  [OK] Alchemy: Configurada (16 EVM + Solana)",
			"  [OK] Alchemy: Configured (16 EVM + Solana)"))
	} else {
		fmt.Println(t(
			"  [ ] Alchemy: Nao configurada (usando APIs publicas)",
			"  [ ] Alchemy: Not configured (using public APIs)"))
	}
	if apiKeys.TrongridKey != "" {
		fmt.Println(t(
			"  [OK] TronGrid: Configurada (Tron + TRC-20)",
			"  [OK] TronGrid: Configured (Tron + TRC-20)"))
	} else {
		fmt.Println(t(
			"  [ ] TronGrid: Nao configurada (usando APIs publicas)",
			"  [ ] TronGrid: Not configured (using public APIs)"))
	}
	fmt.Println()
	fmt.Println(t(
		"  Redes sem API paga configurada usam APIs publicas gratuitas.",
		"  Networks without paid API use free public APIs."))
	fmt.Println(t(
		"  Se uma API paga falhar (creditos acabaram, timeout, etc),",
		"  If a paid API fails (credits ran out, timeout, etc),"))
	fmt.Println(t(
		"  o programa continua automaticamente pelas APIs publicas",
		"  the program automatically continues using public APIs"))
	fmt.Println(t(
		"  sem interrupcao e sem perda de resultados.",
		"  without interruption and without losing results."))
}

func maskKey(key string) string {
	if key == "" {
		return t("(nao configurada)", "(not configured)")
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

// ============================================================================
// RESUMO DO SCAN
// ============================================================================

func showScanSummary(config *ScanConfig) bool {
	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  " + t("RESUMO DO ESCANEAMENTO", "SCAN SUMMARY"))
	fmt.Println("================================================================")
	fmt.Println()

	seedCount := len(config.Seeds)
	fmt.Printf("  %s: %d\n", t("Seeds para escanear", "Seeds to scan"), seedCount)
	fmt.Printf("  %s: %s\n", t("Origem", "Source"), config.SeedSource)
	fmt.Printf("  %s: %d - %d\n", t("Range de indices", "Index range"), config.StartIndex, config.EndIndex)

	if config.Passphrase != "" {
		fmt.Printf("  %s: %s\n", t("Passphrase", "Passphrase"), "***")
	}

	fmt.Println()
	fmt.Println(t("  Redes e derivacoes selecionadas:", "  Selected networks and derivations:"))
	fmt.Println()

	totalChecks := 0
	for _, ng := range config.Networks {
		if !ng.Enabled {
			continue
		}
		name := ng.NameEN
		if uiLanguage == "pt" {
			name = ng.NamePT
		}
		fmt.Printf("    [X] %s\n", name)
		for _, dp := range ng.Derivations {
			if !dp.Enabled {
				continue
			}
			dpName := dp.NameEN
			if uiLanguage == "pt" {
				dpName = dp.NamePT
			}
			fmt.Printf("        - %s [%s]\n", dpName, dp.Path)
			totalChecks++
		}
	}

	indexRange := config.EndIndex - config.StartIndex + 1
	totalAddresses := seedCount * totalChecks * indexRange

	fmt.Println()
	fmt.Printf("  %s: %s\n", t("Total de enderecos a verificar", "Total addresses to check"),
		formatNumber(totalAddresses))
	fmt.Println()

	for {
		choice := getUserInput(t("  Iniciar escaneamento? (S/N): ", "  Start scanning? (Y/N): "))
		upper := strings.ToUpper(choice)
		if upper == "S" || upper == "Y" {
			return true
		}
		if upper == "N" {
			return false
		}
		fmt.Println(t(
			"  [!] Opcao invalida. Digite S (sim) ou N (nao).",
			"  [!] Invalid option. Enter Y (yes) or N (no)."))
	}
}

func formatNumber(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}

	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
