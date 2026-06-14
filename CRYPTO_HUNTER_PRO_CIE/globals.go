package main

import (
	"bufio"
	"sync"
)

// ============================================================================
// VARIAVEIS GLOBAIS
// ============================================================================

var (
	uiLanguage    string = "en"
	globalScanner *bufio.Scanner
	apiKeys       APIKeys
)

// ============================================================================
// STRUCTS
// ============================================================================

type APIKeys struct {
	AlchemyKey  string
	TrongridKey string
}

type NetworkGroup struct {
	ID          string
	NamePT      string
	NameEN      string
	Enabled     bool
	Derivations []DerivationPath
	EVMNetworks []string
}

type DerivationPath struct {
	ID          string
	NamePT      string
	NameEN      string
	Path        string
	Purpose     uint32
	CoinType    uint32
	AddressType string
	Enabled     bool
}

type ScanConfig struct {
	Seeds        []string
	SeedSource   string
	SkipChecksum bool
	Networks     []NetworkGroup
	StartIndex   int
	EndIndex     int
	Passphrase   string
}

type ScanResult struct {
	SeedPhrase     string
	Network        string
	DerivationPath string
	Index          int
	Address        string
	PrivateKey     string
	NativeBalance  string
	NativeSymbol   string
	HasBalance     bool
	HasHistory     bool
	TxCount        int
	LastTxDate     string
	Tokens         []TokenResult
}

type TokenResult struct {
	Symbol   string
	Name     string
	Balance  string
	Contract string
}

type APIEndpoint struct {
	URL       string
	Name      string
	RateLimit int
	mu        sync.Mutex
	lastUsed  int64
}

type APIPool struct {
	endpoints []*APIEndpoint
	current   int
	mu        sync.Mutex
}

// ============================================================================
// CONSTANTES - REDES E DERIVACOES
// ============================================================================

func getDefaultNetworkGroups() []NetworkGroup {
	return []NetworkGroup{
		{
			ID:     "evm",
			NamePT: "EVM (ETH, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Linea, Scroll, Gnosis, zkSync, Blast, Cronos, Celo, Berachain, Sonic, Mantle, Flare)",
			NameEN: "EVM (ETH, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Linea, Scroll, Gnosis, zkSync, Blast, Cronos, Celo, Berachain, Sonic, Mantle, Flare)",
			Enabled: false,
			EVMNetworks: []string{"ethereum", "bsc", "polygon", "arbitrum", "avalanche", "optimism", "base", "linea", "scroll", "gnosis", "zksync", "blast", "cronos", "celo", "berachain", "sonic", "mantle", "flare"},
			Derivations: []DerivationPath{
				{ID: "evm_standard", NamePT: "Padrao EVM - Todas as redes EVM usam o mesmo endereco (0x...)", NameEN: "Standard EVM - All EVM networks share the same address (0x...)",
					Path: "m/44'/60'/0'/0/x", Purpose: 44, CoinType: 60,
					AddressType: "evm", Enabled: true},
			},
		},
		{
			ID:     "btc",
			NamePT: "Bitcoin (BTC)",
			NameEN: "Bitcoin (BTC)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "btc_legacy", NamePT: "Legacy - Enderecos antigos (1...)", NameEN: "Legacy - Old addresses (1...)",
					Path: "m/44'/0'/0'/0/x", Purpose: 44, CoinType: 0,
					AddressType: "legacy", Enabled: false},
				{ID: "btc_segwit", NamePT: "SegWit - P2SH-SegWit (3...)", NameEN: "SegWit - P2SH-SegWit (3...)",
					Path: "m/49'/0'/0'/0/x", Purpose: 49, CoinType: 0,
					AddressType: "segwit", Enabled: false},
				{ID: "btc_native", NamePT: "Native SegWit - Bech32 (bc1q...)", NameEN: "Native SegWit - Bech32 (bc1q...)",
					Path: "m/84'/0'/0'/0/x", Purpose: 84, CoinType: 0,
					AddressType: "native", Enabled: false},
				{ID: "btc_taproot", NamePT: "Taproot - Mais recente (bc1p...)", NameEN: "Taproot - Newest (bc1p...)",
					Path: "m/86'/0'/0'/0/x", Purpose: 86, CoinType: 0,
					AddressType: "taproot", Enabled: false},
			},
		},
		{
			ID:     "bch",
			NamePT: "Bitcoin Cash (BCH)",
			NameEN: "Bitcoin Cash (BCH)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "bch_cashaddr", NamePT: "CashAddr - Formato padrao BCH (bitcoincash:q...)", NameEN: "CashAddr - Standard BCH format (bitcoincash:q...)",
					Path: "m/44'/145'/0'/0/x", Purpose: 44, CoinType: 145,
					AddressType: "cashaddr", Enabled: false},
				{ID: "bch_legacy", NamePT: "Legacy - Formato antigo compartilhado com BTC (1...)", NameEN: "Legacy - Old format shared with BTC (1...)",
					Path: "m/44'/145'/0'/0/x", Purpose: 44, CoinType: 145,
					AddressType: "legacy", Enabled: false},
			},
		},
		{
			ID:     "trx",
			NamePT: "Tron (TRX)",
			NameEN: "Tron (TRX)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "trx_standard", NamePT: "Padrao Tron (T...)", NameEN: "Standard Tron (T...)",
					Path: "m/44'/195'/0'/0/x", Purpose: 44, CoinType: 195,
					AddressType: "tron", Enabled: true},
			},
		},
		{
			ID:     "sol",
			NamePT: "Solana (SOL)",
			NameEN: "Solana (SOL)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "sol_standard", NamePT: "Padrao Solana - Ed25519 (Base58, 32-44 chars)", NameEN: "Standard Solana - Ed25519 (Base58, 32-44 chars)",
					Path: "m/44'/501'/0'/0'", Purpose: 44, CoinType: 501,
					AddressType: "solana", Enabled: true},
			},
		},
		{
			ID:     "ltc",
			NamePT: "Litecoin (LTC)",
			NameEN: "Litecoin (LTC)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "ltc_legacy", NamePT: "Legacy (L...)", NameEN: "Legacy (L...)",
					Path: "m/44'/2'/0'/0/x", Purpose: 44, CoinType: 2,
					AddressType: "ltc_legacy", Enabled: false},
				{ID: "ltc_segwit", NamePT: "SegWit (M...)", NameEN: "SegWit (M...)",
					Path: "m/49'/2'/0'/0/x", Purpose: 49, CoinType: 2,
					AddressType: "ltc_segwit", Enabled: false},
				{ID: "ltc_native", NamePT: "Native SegWit (ltc1q...)", NameEN: "Native SegWit (ltc1q...)",
					Path: "m/84'/2'/0'/0/x", Purpose: 84, CoinType: 2,
					AddressType: "ltc_native", Enabled: false},
			},
		},
		{
			ID:     "doge",
			NamePT: "Dogecoin (DOGE)",
			NameEN: "Dogecoin (DOGE)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "doge_standard", NamePT: "Padrao (D...)", NameEN: "Standard (D...)",
					Path: "m/44'/3'/0'/0/x", Purpose: 44, CoinType: 3,
					AddressType: "doge", Enabled: true},
			},
		},
		{
			ID:     "ton",
			NamePT: "TON (Toncoin)",
			NameEN: "TON (Toncoin)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "ton_standard", NamePT: "Padrao TON - Ed25519 (UQ... / EQ...)", NameEN: "Standard TON - Ed25519 (UQ... / EQ...)",
					Path: "m/44'/607'/0'", Purpose: 44, CoinType: 607,
					AddressType: "ton", Enabled: true},
			},
		},
		{
			ID:     "zec",
			NamePT: "Zcash (ZEC)",
			NameEN: "Zcash (ZEC)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "zec_transparent", NamePT: "Transparente (t1...)", NameEN: "Transparent (t1...)",
					Path: "m/44'/133'/0'/0/x", Purpose: 44, CoinType: 133,
					AddressType: "zcash", Enabled: true},
			},
		},
		{
			ID:     "xrp",
			NamePT: "XRP (Ripple)",
			NameEN: "XRP (Ripple)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "xrp_standard", NamePT: "Padrao XRP (r...)", NameEN: "Standard XRP (r...)",
					Path: "m/44'/144'/0'/0/x", Purpose: 44, CoinType: 144,
					AddressType: "xrp", Enabled: true},
			},
		},
		{
			ID:     "xlm",
			NamePT: "Stellar (XLM)",
			NameEN: "Stellar (XLM)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "xlm_standard", NamePT: "Padrao Stellar - Ed25519 (G...)", NameEN: "Standard Stellar - Ed25519 (G...)",
					Path: "m/44'/148'/0'", Purpose: 44, CoinType: 148,
					AddressType: "stellar", Enabled: true},
			},
		},
		{
			ID:     "algo",
			NamePT: "Algorand (ALGO)",
			NameEN: "Algorand (ALGO)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "algo_standard", NamePT: "Padrao Algorand - Ed25519 (Base32)", NameEN: "Standard Algorand - Ed25519 (Base32)",
					Path: "m/44'/283'/0'/0/x", Purpose: 44, CoinType: 283,
					AddressType: "algorand", Enabled: true},
			},
		},
		{
			ID:     "sui",
			NamePT: "Sui (SUI)",
			NameEN: "Sui (SUI)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "sui_standard", NamePT: "Padrao Sui - Ed25519 (0x...)", NameEN: "Standard Sui - Ed25519 (0x...)",
					Path: "m/44'/784'/0'/0'/0'", Purpose: 44, CoinType: 784,
					AddressType: "sui", Enabled: true},
			},
		},
		{
			ID:     "near",
			NamePT: "Near Protocol (NEAR)",
			NameEN: "Near Protocol (NEAR)",
			Enabled: false,
			Derivations: []DerivationPath{
				{ID: "near_standard", NamePT: "Padrao NEAR - Ed25519 (hex implicit)", NameEN: "Standard NEAR - Ed25519 (hex implicit)",
					Path: "m/44'/397'/0'", Purpose: 44, CoinType: 397,
					AddressType: "near", Enabled: true},
			},
		},
	}
}

// ============================================================================
// TOKENS
// ============================================================================

type TokenInfo struct {
	Symbol   string
	Name     string
	Contract string
	Decimals int
}

func getEVMTokens() map[string][]TokenInfo {
	return map[string][]TokenInfo{
		"ethereum": {
			{Symbol: "USDT", Name: "Tether USD", Contract: "0xdAC17F958D2ee523a2206206994597C13D831ec7", Decimals: 6},
			{Symbol: "USDC", Name: "USD Coin", Contract: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", Decimals: 6},
			{Symbol: "DAI", Name: "Dai Stablecoin", Contract: "0x6B175474E89094C44Da98b954EedeAC495271d0F", Decimals: 18},
			{Symbol: "BUSD", Name: "Binance USD", Contract: "0x4Fabb145d64652a948d72533023f6E7A623C7C53", Decimals: 18},
			{Symbol: "TUSD", Name: "TrueUSD", Contract: "0x0000000000085d4780B73119b644AE5ecd22b376", Decimals: 18},
			{Symbol: "FRAX", Name: "Frax", Contract: "0x853d955aCEf822Db058eb8505911ED77F175b99e", Decimals: 18},
			{Symbol: "WBTC", Name: "Wrapped BTC", Contract: "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599", Decimals: 8},
			{Symbol: "WETH", Name: "Wrapped ETH", Contract: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", Decimals: 18},
			{Symbol: "LINK", Name: "Chainlink", Contract: "0x514910771AF9Ca656af840dff83E8264EcF986CA", Decimals: 18},
			{Symbol: "UNI", Name: "Uniswap", Contract: "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984", Decimals: 18},
			{Symbol: "AAVE", Name: "Aave", Contract: "0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9", Decimals: 18},
			{Symbol: "MKR", Name: "Maker", Contract: "0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2", Decimals: 18},
			{Symbol: "COMP", Name: "Compound", Contract: "0xc00e94Cb662C3520282E6f5717214004A7f26888", Decimals: 18},
			{Symbol: "SNX", Name: "Synthetix", Contract: "0xC011a73ee8576Fb46F5E1c5751cA3B9Fe0af2a6F", Decimals: 18},
			{Symbol: "CRV", Name: "Curve DAO", Contract: "0xD533a949740bb3306d119CC777fa900bA034cd52", Decimals: 18},
			{Symbol: "LDO", Name: "Lido DAO", Contract: "0x5A98FcBEA516Cf06857215779Fd812CA3beF1B32", Decimals: 18},
			{Symbol: "SHIB", Name: "Shiba Inu", Contract: "0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE", Decimals: 18},
			{Symbol: "PEPE", Name: "Pepe", Contract: "0x6982508145454Ce325dDbE47a25d4ec3d2311933", Decimals: 18},
			{Symbol: "POL", Name: "Polygon (ERC-20)", Contract: "0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0", Decimals: 18},
			{Symbol: "APE", Name: "ApeCoin", Contract: "0x4d224452801ACEd8B2F0aebE155379bb5D594381", Decimals: 18},
			{Symbol: "SAND", Name: "The Sandbox", Contract: "0x3845badAde8e6dFF049820680d1F14bD3903a5d0", Decimals: 18},
			{Symbol: "MANA", Name: "Decentraland", Contract: "0x0F5D2fB29fb7d3CFeE444a200298f468908cC942", Decimals: 18},
			{Symbol: "GRT", Name: "The Graph", Contract: "0xc944E90C64B2c07662A292be6244BDf05Cda44a7", Decimals: 18},
			{Symbol: "ENS", Name: "Ethereum Name Service", Contract: "0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72", Decimals: 18},
		},
		"bsc": {
			{Symbol: "USDT", Name: "Tether USD (BSC)", Contract: "0x55d398326f99059fF775485246999027B3197955", Decimals: 18},
			{Symbol: "USDC", Name: "USD Coin (BSC)", Contract: "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d", Decimals: 18},
			{Symbol: "BUSD", Name: "Binance USD (BSC)", Contract: "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56", Decimals: 18},
			{Symbol: "DAI", Name: "Dai (BSC)", Contract: "0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3", Decimals: 18},
			{Symbol: "CAKE", Name: "PancakeSwap", Contract: "0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82", Decimals: 18},
			{Symbol: "WBNB", Name: "Wrapped BNB", Contract: "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", Decimals: 18},
			{Symbol: "XVS", Name: "Venus", Contract: "0xcF6BB5389c92Bdda8a3747Ddb454cB7a64626C63", Decimals: 18},
			{Symbol: "BTCB", Name: "Bitcoin BEP2", Contract: "0x7130d2A12B9BCbFAe4f2634d864A1Ee1Ce3Ead9c", Decimals: 18},
			{Symbol: "ETH", Name: "Ethereum (BSC)", Contract: "0x2170Ed0880ac9A755fd29B2688956BD959F933F8", Decimals: 18},
			{Symbol: "LINK", Name: "Chainlink (BSC)", Contract: "0xF8A0BF9cF54Bb92F17374d9e9A321E6a111a51bD", Decimals: 18},
			{Symbol: "DOT", Name: "Polkadot (BSC)", Contract: "0x7083609fCE4d1d8Dc0C979AAb8c869Ea2C873402", Decimals: 18},
			{Symbol: "ADA", Name: "Cardano (BSC)", Contract: "0x3EE2200Efb3400fAbB9AacF31297cBdD1d435D47", Decimals: 18},
			{Symbol: "DOGE", Name: "Dogecoin (BSC)", Contract: "0xbA2aE424d960c26247Dd6c32edC70B295c744C43", Decimals: 8},
			{Symbol: "FLOKI", Name: "Floki Inu (BSC)", Contract: "0xfb5B838b6cfEEdC2873aB27866079AC55363D37E", Decimals: 9},
		},
		"polygon": {
			{Symbol: "USDT", Name: "Tether USD (Polygon)", Contract: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F", Decimals: 6},
			{Symbol: "USDC", Name: "USD Coin (Polygon)", Contract: "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359", Decimals: 6},
			{Symbol: "USDC.e", Name: "Bridged USDC (Polygon)", Contract: "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", Decimals: 6},
			{Symbol: "DAI", Name: "Dai (Polygon)", Contract: "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063", Decimals: 18},
			{Symbol: "WPOL", Name: "Wrapped POL", Contract: "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270", Decimals: 18},
			{Symbol: "AAVE", Name: "Aave (Polygon)", Contract: "0xD6DF932A45C0f255f85145f286eA0b292B21C90B", Decimals: 18},
			{Symbol: "LINK", Name: "Chainlink (Polygon)", Contract: "0x53E0bca35eC356BD5ddDFebbD1Fc0fD03FaBad39", Decimals: 18},
			{Symbol: "WBTC", Name: "Wrapped BTC (Polygon)", Contract: "0x1BFD67037B42Cf73acF2047067bd4F2C47D9BfD6", Decimals: 8},
			{Symbol: "WETH", Name: "Wrapped ETH (Polygon)", Contract: "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619", Decimals: 18},
			{Symbol: "CRV", Name: "Curve DAO (Polygon)", Contract: "0x172370d5Cd63279eFa6d502DAB29171933a610AF", Decimals: 18},
			{Symbol: "BAL", Name: "Balancer (Polygon)", Contract: "0x9a71012B13CA4d3D0Cdc72A177DF3ef03b0E76A3", Decimals: 18},
			{Symbol: "GNS", Name: "Gains Network (Polygon)", Contract: "0xE5417Af564e4bFDA1c483642db72007871397896", Decimals: 18},
			{Symbol: "QUICK", Name: "QuickSwap (Polygon)", Contract: "0xB5C064F955D8e7F38fE0460C556a72987494eE17", Decimals: 18},
		},
		"arbitrum": {
			{Symbol: "USDT", Name: "Tether USD (Arbitrum)", Contract: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9", Decimals: 6},
			{Symbol: "USDC", Name: "USD Coin (Arbitrum)", Contract: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831", Decimals: 6},
			{Symbol: "USDC.e", Name: "Bridged USDC (Arbitrum)", Contract: "0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8", Decimals: 6},
			{Symbol: "DAI", Name: "Dai (Arbitrum)", Contract: "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1", Decimals: 18},
			{Symbol: "WETH", Name: "Wrapped ETH (Arbitrum)", Contract: "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1", Decimals: 18},
			{Symbol: "ARB", Name: "Arbitrum", Contract: "0x912CE59144191C1204E64559FE8253a0e49E6548", Decimals: 18},
			{Symbol: "LINK", Name: "Chainlink (Arbitrum)", Contract: "0xf97f4df75117a78c1A5a0DBb814Af92458539FB4", Decimals: 18},
			{Symbol: "GRT", Name: "The Graph (Arbitrum)", Contract: "0x9623063377AD1B27544C965cCd7342f7EA7e88C7", Decimals: 18},
			{Symbol: "GMX", Name: "GMX", Contract: "0xfc5A1A6EB076a2C7aD06eD22C90d7E710E35ad0a", Decimals: 18},
			{Symbol: "MAGIC", Name: "Magic (Treasure)", Contract: "0x539bdE0d7Dbd336b79148AA742883198BBF60342", Decimals: 18},
			{Symbol: "RDNT", Name: "Radiant Capital", Contract: "0x3082CC23568eA640225c2467653dB90e9250AaA0", Decimals: 18},
			{Symbol: "PENDLE", Name: "Pendle", Contract: "0x0c880f6761F1af8d9Aa9C466984b80DAb9a8c9e8", Decimals: 18},
		},
		"avalanche": {
			{Symbol: "USDT", Name: "Tether USD (Avalanche)", Contract: "0x9702230A8Ea53601f5cD2dc00fDBc13d4dF4A8c7", Decimals: 6},
			{Symbol: "USDC", Name: "USD Coin (Avalanche)", Contract: "0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E", Decimals: 6},
			{Symbol: "DAI.e", Name: "Dai (Avalanche)", Contract: "0xd586E7F844cEa2F87f50152665BCbc2C279D8d70", Decimals: 18},
			{Symbol: "WAVAX", Name: "Wrapped AVAX", Contract: "0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7", Decimals: 18},
			{Symbol: "WETH.e", Name: "Wrapped ETH (Avalanche)", Contract: "0x49D5c2BdFfac6CE2BFdB6640F4F80f226bc10bAB", Decimals: 18},
			{Symbol: "WBTC.e", Name: "Wrapped BTC (Avalanche)", Contract: "0x50b7545627a5162F82A992c33b87aDc75187B218", Decimals: 8},
			{Symbol: "AAVE.e", Name: "Aave (Avalanche)", Contract: "0x63a72806098Bd3D9520cC43356dD78afe5D386D9", Decimals: 18},
			{Symbol: "LINK.e", Name: "Chainlink (Avalanche)", Contract: "0x5947BB275c521040051D82396192181b413227A3", Decimals: 18},
			{Symbol: "JOE", Name: "Trader Joe", Contract: "0x6e84a6216eA6dACC71eE8E6b0a5B7322EEbC0fDd", Decimals: 18},
			{Symbol: "GMX", Name: "GMX (Avalanche)", Contract: "0x62edc0692BD897D2295872a9FFCac5425011c661", Decimals: 18},
			{Symbol: "sAVAX", Name: "Staked AVAX (Benqi)", Contract: "0x2b2C81e08f1Af8835a78Bb2A90AE924ACE0eA4bE", Decimals: 18},
		},
		"optimism": {
			{Symbol: "USDT", Name: "Tether USD (Optimism)", Contract: "0x94b008aA00579c1307B0EF2c499aD98a8ce58e58", Decimals: 6},
			{Symbol: "USDC", Name: "USD Coin (Optimism)", Contract: "0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85", Decimals: 6},
			{Symbol: "DAI", Name: "Dai (Optimism)", Contract: "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1", Decimals: 18},
			{Symbol: "OP", Name: "Optimism", Contract: "0x4200000000000000000000000000000000000042", Decimals: 18},
			{Symbol: "WETH", Name: "Wrapped ETH (Optimism)", Contract: "0x4200000000000000000000000000000000000006", Decimals: 18},
			{Symbol: "LINK", Name: "Chainlink (Optimism)", Contract: "0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6", Decimals: 18},
			{Symbol: "SNX", Name: "Synthetix (Optimism)", Contract: "0x8700dAec35aF8Ff88c16BdF0418774CB3D7599B4", Decimals: 18},
			{Symbol: "VELO", Name: "Velodrome", Contract: "0x9560e827aF36c94D2Ac33a39bCE1Fe78631088Db", Decimals: 18},
			{Symbol: "AAVE", Name: "Aave (Optimism)", Contract: "0x76FB31fb4af56892A25e32cFC43De717950c9278", Decimals: 18},
			{Symbol: "PERP", Name: "Perpetual Protocol", Contract: "0x9e1028F5F1D5eDE59748FFceE5532509976840E0", Decimals: 18},
		},
		"base": {
			{Symbol: "USDC", Name: "USD Coin (Base)", Contract: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", Decimals: 6},
			{Symbol: "DAI", Name: "Dai (Base)", Contract: "0x50c5725949A6F0c72E6C4a641F24049A917DB0Cb", Decimals: 18},
			{Symbol: "WETH", Name: "Wrapped ETH (Base)", Contract: "0x4200000000000000000000000000000000000006", Decimals: 18},
			{Symbol: "cbETH", Name: "Coinbase Wrapped Staked ETH", Contract: "0x2Ae3F1Ec7F1F5012CFEab0185bfc7aa3cf0DEc22", Decimals: 18},
			{Symbol: "AERO", Name: "Aerodrome Finance", Contract: "0x940181a94A35A4569E4529A3CDfB74e38FD98631", Decimals: 18},
			{Symbol: "BRETT", Name: "Brett", Contract: "0x532f27101965dd16442E59d40670FaF5eBB142E4", Decimals: 18},
			{Symbol: "DEGEN", Name: "Degen", Contract: "0x4ed4E862860beD51a9570b96d89aF5E1B0Efefed", Decimals: 18},
			{Symbol: "TOSHI", Name: "Toshi", Contract: "0xAC1Bd2486aAf3B5C0fc3Fd868558b082a531B2B4", Decimals: 18},
			{Symbol: "USDbC", Name: "Bridged USDC (Base)", Contract: "0xd9aAEc86B65D86f6A7B5B1b0c42FFA531710b6CA", Decimals: 6},
		},
		"sonic": {
			{Symbol: "USDC.e", Name: "Bridged USDC (Sonic)", Contract: "0x29219dd400f2Bf60E5a23d13Be72B486D4038894", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (Sonic)", Contract: "0x50c42dEAcD8Fc9773493ED674b675bE577f2634b", Decimals: 18},
			{Symbol: "wS", Name: "Wrapped Sonic", Contract: "0x039e2fB66102314Ce7b64Ce5Ce3E5183bc94aD38", Decimals: 18},
			{Symbol: "BRUSH", Name: "PaintSwap", Contract: "0x85dec8c4B2680793661bCA91a8F129607571863d", Decimals: 18},
			{Symbol: "EQUAL", Name: "Equalizer DEX", Contract: "0x3Fd3A0c85B70754eFc07aC9Ac0cbBDCe664865A6", Decimals: 18},
			{Symbol: "USDT", Name: "Tether USD (Sonic)", Contract: "0xfA1FBb8Ef55A4855E5688C0eE13aC3f202486286", Decimals: 6},
		},
		"mantle": {
			{Symbol: "USDT", Name: "Tether USD (Mantle)", Contract: "0x201EBa5CC46D216Ce6DC03F6a759e8E766e956aE", Decimals: 6},
			{Symbol: "USDC", Name: "USD Coin (Mantle)", Contract: "0x09Bc4E0D864854c6aFB6eB9A9cdF58aC190D0dF9", Decimals: 6},
			{Symbol: "WMNT", Name: "Wrapped MNT", Contract: "0x78c1b0C915c4FAA5FffA6CAbf0219DA63d7f4cb8", Decimals: 18},
			{Symbol: "WETH", Name: "Wrapped ETH (Mantle)", Contract: "0xdEAddEaDdeadDEadDEADDEAddEADDEAddead1111", Decimals: 18},
			{Symbol: "mETH", Name: "Mantle Staked ETH", Contract: "0xcDA86A272531e8640cD7F1a92c01839911B90bb0", Decimals: 18},
			{Symbol: "PUFF", Name: "Puff", Contract: "0x26a6b0dcdCfb981362aFA56D581e4A7dBA3Be140", Decimals: 18},
		},
		"flare": {
			{Symbol: "WFLR", Name: "Wrapped FLR", Contract: "0x1D80c49BbBCd1C0911346656B529DF9E5c2F783d", Decimals: 18},
			{Symbol: "USDC.e", Name: "Bridged USDC (Flare)", Contract: "0xFbDa5F676cB37624f28265A144A48B0d6e87d3b6", Decimals: 6},
			{Symbol: "USDT.e", Name: "Bridged USDT (Flare)", Contract: "0x96B41289D90444B8adD57e6F265DB5aE8651c446", Decimals: 6},
			{Symbol: "sFLR", Name: "Staked FLR (Sceptre)", Contract: "0x12e605bc104e93B45e1aD99F9e555f659051c2BB", Decimals: 18},
			{Symbol: "WETH", Name: "Wrapped ETH (Flare)", Contract: "0x1502FA4be69d526124D28A3A863C1E03b7C47E8c", Decimals: 18},
		},
		"linea": {
			{Symbol: "USDC", Name: "USD Coin (Linea)", Contract: "0x176211869cA2b568f2A7D4EE941E073a821EE1ff", Decimals: 6},
			{Symbol: "USDT", Name: "Tether USD (Linea)", Contract: "0xA219439258ca9da29E9Cc4cE5596924745e12B93", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (Linea)", Contract: "0xe5D7C2a44FfDDf6b295A15c148167daaAf5Cf34f", Decimals: 18},
			{Symbol: "DAI", Name: "Dai (Linea)", Contract: "0x4AF15ec2A0BD43Db75dd04E62FAA3B8EF36b00d5", Decimals: 18},
			{Symbol: "wstETH", Name: "Wrapped stETH (Linea)", Contract: "0xB5beDd42000b71FddE22D3eE8a79Bd49A568fC8F", Decimals: 18},
		},
		"scroll": {
			{Symbol: "USDC", Name: "USD Coin (Scroll)", Contract: "0x06eFdBFf2a14a7c8E15944D1F4A48F9F95F663A4", Decimals: 6},
			{Symbol: "USDT", Name: "Tether USD (Scroll)", Contract: "0xf55BEC9cafDbE8730f096Aa55dad6D22d44099Df", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (Scroll)", Contract: "0x5300000000000000000000000000000000000004", Decimals: 18},
			{Symbol: "wstETH", Name: "Wrapped stETH (Scroll)", Contract: "0xf610A9dfB7C89644979b4A0f27063E9e7d7Cda32", Decimals: 18},
			{Symbol: "SCR", Name: "Scroll Token", Contract: "0xd29687c813D741E2F938F4aC377128810E217b1b", Decimals: 18},
		},
		"gnosis": {
			{Symbol: "USDC", Name: "USD Coin (Gnosis)", Contract: "0xDDAfbb505ad214D7b80b1f830fcCc89B60fb7A83", Decimals: 6},
			{Symbol: "USDT", Name: "Tether USD (Gnosis)", Contract: "0x4ECaBa5870353805a9F068101A40E0f32ed605C6", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (Gnosis)", Contract: "0x6A023CCd1ff6F2045C3309768eAd9E68F978f6e1", Decimals: 18},
			{Symbol: "GNO", Name: "Gnosis Token", Contract: "0x9C58BAcC331c9aa871AFD802DB6379a98e80CEdb", Decimals: 18},
			{Symbol: "sDAI", Name: "Savings xDAI", Contract: "0xaf204776c7245bF4147c2612BF6e5972Ee483701", Decimals: 18},
			{Symbol: "wstETH", Name: "Wrapped stETH (Gnosis)", Contract: "0x6C76971f98945AE98dD7d4DFcA8711ebea946eA6", Decimals: 18},
		},
		"zksync": {
			{Symbol: "USDC", Name: "USD Coin (zkSync)", Contract: "0x1d17CBcF0D6D143135aE902365D2E5e2A16538D4", Decimals: 6},
			{Symbol: "USDT", Name: "Tether USD (zkSync)", Contract: "0x493257fD37EDB34451f62EDf8D2a0C418852bA4C", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (zkSync)", Contract: "0x5AEa5775959fBC2557Cc8789bC1bf90A239D9a91", Decimals: 18},
			{Symbol: "ZK", Name: "ZKsync Token", Contract: "0x5A7d6b2F92C77FAD6CCaBd7EE0624E64907Eaf3E", Decimals: 18},
		},
		"blast": {
			{Symbol: "USDB", Name: "USDB (Blast)", Contract: "0x4300000000000000000000000000000000000003", Decimals: 18},
			{Symbol: "WETH", Name: "Wrapped ETH (Blast)", Contract: "0x4300000000000000000000000000000000000004", Decimals: 18},
			{Symbol: "USDT", Name: "Tether USD (Blast)", Contract: "0x0C0Cf4ECa0110b1b9a0DE7125aF1E13e6D314a0F", Decimals: 6},
		},
		"cronos": {
			{Symbol: "USDC", Name: "USD Coin (Cronos)", Contract: "0xc21223249CA28397B4B6541dfFaEcC539BfF0c59", Decimals: 6},
			{Symbol: "USDT", Name: "Tether USD (Cronos)", Contract: "0x66e428c3f67a68878562e79A0234c1F83c208770", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (Cronos)", Contract: "0xe44Fd7fCb2b1581822D0c862B68222998a0c299a", Decimals: 18},
			{Symbol: "WCRO", Name: "Wrapped CRO", Contract: "0x5C7F8A570d578ED60E9120223616C988BD4FBBA0", Decimals: 18},
		},
		"celo": {
			{Symbol: "cUSD", Name: "Celo Dollar", Contract: "0x765DE816845861e75A25fCA122bb6898B8B1282a", Decimals: 18},
			{Symbol: "cEUR", Name: "Celo Euro", Contract: "0xD8763CBa276a3738E6DE85b4b3bF5FDed6D6cA73", Decimals: 18},
			{Symbol: "USDC", Name: "USD Coin (Celo)", Contract: "0xcebA9300f2b948710d2653dD7B07f33A8B32118C", Decimals: 6},
			{Symbol: "USDT", Name: "Tether USD (Celo)", Contract: "0x48065fbBE25f71C9282ddf5e1cD6D6A887483D5e", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (Celo)", Contract: "0x122013fd7dF1C6F636a5bb8f03108E876548b455", Decimals: 18},
		},
		"berachain": {
			{Symbol: "HONEY", Name: "Honey (Berachain)", Contract: "0x0E4aaF1351de4c0264C5c7056Ef3777b41BD8e03", Decimals: 18},
			{Symbol: "USDC", Name: "USD Coin (Berachain)", Contract: "0xd6D83aF58a19Cd14eF3CF6fe848C9A4d21e5727c", Decimals: 6},
			{Symbol: "USDT", Name: "Tether USD (Berachain)", Contract: "0x05D0dD5135E3eF3aDE32a9eF9Cb06DeFAEA36927", Decimals: 6},
			{Symbol: "WETH", Name: "Wrapped ETH (Berachain)", Contract: "0x6E1E9896e93F7A71ECB33d4386b49DeeD67a231A", Decimals: 18},
			{Symbol: "WBTC", Name: "Wrapped BTC (Berachain)", Contract: "0x2577D24a26f8FA19c1058a8b0106E2c7303454a4", Decimals: 8},
			{Symbol: "WBERA", Name: "Wrapped BERA", Contract: "0x7507c1dc16935B82698e4C63f2746A2fCf994dF8", Decimals: 18},
		},
	}
}

func getTRXTokens() []TokenInfo {
	return []TokenInfo{
		{Symbol: "USDT", Name: "Tether USD (TRC-20)", Contract: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", Decimals: 6},
		{Symbol: "USDC", Name: "USD Coin (TRC-20)", Contract: "TEkxiTehnzSmSe2XqrBj4w32RUN966rdz8", Decimals: 6},
		{Symbol: "USDD", Name: "USDD", Contract: "TPYmHEhy5n8TCEfYGqW2rPxsghSfzghPDn", Decimals: 18},
		{Symbol: "TUSD", Name: "TrueUSD (TRC-20)", Contract: "TUpMhErZL2fhh4sVNULAbNKLokS4GjC1F4", Decimals: 18},
		{Symbol: "WTRX", Name: "Wrapped TRX", Contract: "TNUC9Qb1rRpS5CbWLmNMxXBjyFoydXjWFR", Decimals: 6},
		{Symbol: "BTT", Name: "BitTorrent", Contract: "TAFjULxiVgT4qWk6UZwjqwZXTSaGaqnVp4", Decimals: 18},
		{Symbol: "JST", Name: "JUST", Contract: "TCFLL5dx5ZJdKnWuesXxi1VPwjLVmWZZy9", Decimals: 18},
		{Symbol: "SUN", Name: "SUN", Contract: "TSSMHYeV2uE9qYH95DqyoCuNCzEL1NvU3S", Decimals: 18},
		{Symbol: "WIN", Name: "WINkLink", Contract: "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7", Decimals: 6},
	}
}

func getSOLTokens() []TokenInfo {
	return []TokenInfo{
		{Symbol: "USDC", Name: "USD Coin (SPL)", Contract: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", Decimals: 6},
		{Symbol: "USDT", Name: "Tether USD (SPL)", Contract: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", Decimals: 6},
		{Symbol: "JUP", Name: "Jupiter", Contract: "JUPyiwrYJFskUPiHa7hkeR8VUtAeFoSYbKedZNsDvCN", Decimals: 6},
		{Symbol: "BONK", Name: "Bonk", Contract: "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263", Decimals: 5},
		{Symbol: "RAY", Name: "Raydium", Contract: "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R", Decimals: 6},
		{Symbol: "ORCA", Name: "Orca", Contract: "orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE", Decimals: 6},
		{Symbol: "JTO", Name: "Jito", Contract: "jtojtomepa8beP8AuQc6eXt5FriJwfFMwQx2v2f9mCL", Decimals: 9},
		{Symbol: "PYTH", Name: "Pyth Network", Contract: "HZ1JovNiVvGrGNiiYvEozEVgZ58xaU3RKwX8eACQBCt3", Decimals: 6},
		{Symbol: "W", Name: "Wormhole", Contract: "85VBFQZC9TZkfaptBWjvUw7YbZjy52A6mjtPGjstQAmQ", Decimals: 6},
		{Symbol: "RENDER", Name: "Render Token", Contract: "rndrizKT3MK1iimdxRdWabcF7Zg7AR5T4nud4EkHBof", Decimals: 8},
		{Symbol: "WIF", Name: "dogwifhat", Contract: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm", Decimals: 6},
		{Symbol: "mSOL", Name: "Marinade Staked SOL", Contract: "mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So", Decimals: 9},
		{Symbol: "jitoSOL", Name: "Jito Staked SOL", Contract: "J1toso1uCk3RLmjorhTtrVwY9HJ7X8V9yYac6Y7kGCPn", Decimals: 9},
	}
}
