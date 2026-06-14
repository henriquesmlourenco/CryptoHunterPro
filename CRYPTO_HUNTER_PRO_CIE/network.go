package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// POOL DE APIs COM ROUND-ROBIN E FALLBACK
// ============================================================================

func NewAPIPool(endpoints []*APIEndpoint) *APIPool {
	return &APIPool{endpoints: endpoints, current: 0}
}

func (p *APIPool) Next() *APIEndpoint {
	p.mu.Lock()
	defer p.mu.Unlock()
	ep := p.endpoints[p.current]
	p.current = (p.current + 1) % len(p.endpoints)
	return ep
}

func (p *APIPool) GetAll() []*APIEndpoint {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.endpoints
}

// ============================================================================
// ALCHEMY RPC ENDPOINTS POR REDE (requer API key paga)
// ============================================================================

var alchemyNetworkIDs = map[string]string{
	"ethereum":  "eth-mainnet",
	"bsc":       "bnb-mainnet",
	"polygon":   "polygon-mainnet",
	"arbitrum":  "arb-mainnet",
	"avalanche": "avax-mainnet",
	"optimism":  "opt-mainnet",
	"base":      "base-mainnet",
	"linea":     "linea-mainnet",
	"scroll":    "scroll-mainnet",
	"gnosis":    "gnosis-mainnet",
	"zksync":    "zksync-mainnet",
	"blast":     "blast-mainnet",
	"celo":      "celo-mainnet",
	"berachain": "berachain-mainnet",
	"sonic":     "sonic-mainnet",
	"mantle":    "mantle-mainnet",
	// cronos e flare NAO sao suportados pela Alchemy
}

// getEVMRPCs retorna a lista de RPCs para uma rede EVM.
// Se Alchemy estiver configurada e a rede for suportada, adiciona como primeiro RPC.
func getEVMRPCs(network string) []string {
	rpcs := make([]string, 0)

	// Se Alchemy configurada, adicionar como RPC prioritario
	if apiKeys.AlchemyKey != "" {
		if netID, ok := alchemyNetworkIDs[network]; ok {
			alchemyURL := fmt.Sprintf("https://%s.g.alchemy.com/v2/%s", netID, apiKeys.AlchemyKey)
			rpcs = append(rpcs, alchemyURL)
		}
	}

	// Adicionar RPCs publicos como fallback
	if publicRPCs, ok := evmRPCs[network]; ok {
		rpcs = append(rpcs, publicRPCs...)
	}

	return rpcs
}

// ============================================================================
// RPCs PUBLICOS POR REDE EVM (nao precisam de API key!)
// ============================================================================

var evmRPCs = map[string][]string{
	"ethereum": {
		"https://eth.llamarpc.com",
		"https://1rpc.io/eth",
		"https://ethereum-rpc.publicnode.com",
	},
	"bsc": {
		"https://bsc-dataseed.binance.org/",
		"https://bsc-dataseed1.defibit.io/",
		"https://bsc-rpc.publicnode.com",
	},
	"polygon": {
		"https://polygon-bor-rpc.publicnode.com",
		"https://1rpc.io/matic",
		"https://polygon.drpc.org",
	},
	"arbitrum": {
		"https://arb1.arbitrum.io/rpc",
		"https://arbitrum-one-rpc.publicnode.com",
	},
	"avalanche": {
		"https://api.avax.network/ext/bc/C/rpc",
		"https://avalanche-c-chain-rpc.publicnode.com",
		"https://1rpc.io/avax/c",
	},
	"optimism": {
		"https://mainnet.optimism.io",
		"https://optimism-rpc.publicnode.com",
	},
	"base": {
		"https://mainnet.base.org",
		"https://base-rpc.publicnode.com",
		"https://1rpc.io/base",
	},
	"sonic": {
		"https://rpc.soniclabs.com",
		"https://sonic.drpc.org",
	},
	"mantle": {
		"https://rpc.mantle.xyz",
		"https://mantle-rpc.publicnode.com",
	},
	"flare": {
		"https://flare-api.flare.network/ext/C/rpc",
		"https://flare.drpc.org",
	},
	"linea": {
		"https://rpc.linea.build",
		"https://linea-rpc.publicnode.com",
		"https://1rpc.io/linea",
	},
	"scroll": {
		"https://rpc.scroll.io",
		"https://scroll-rpc.publicnode.com",
		"https://1rpc.io/scroll",
	},
	"gnosis": {
		"https://rpc.gnosischain.com",
		"https://gnosis-rpc.publicnode.com",
		"https://1rpc.io/gnosis",
	},
	"zksync": {
		"https://mainnet.era.zksync.io",
		"https://zksync-era-rpc.publicnode.com",
		"https://1rpc.io/zksync2-era",
	},
	"blast": {
		"https://rpc.blast.io",
		"https://blast-rpc.publicnode.com",
	},
	"cronos": {
		"https://evm.cronos.org",
		"https://cronos-evm-rpc.publicnode.com",
	},
	"celo": {
		"https://forno.celo.org",
		"https://celo-rpc.publicnode.com",
		"https://1rpc.io/celo",
	},
	"berachain": {
		"https://rpc.berachain.com",
		"https://berachain-rpc.publicnode.com",
	},
}

// ============================================================================
// POOLS DE APIs POR REDE (BTC, BCH, TRX, LTC, DOGE)
// ============================================================================

var (
	btcPool      *APIPool
	bchPool      *APIPool
	trxPool      *APIPool
	ltcPool      *APIPool
	dogePool     *APIPool
	poolInitOnce sync.Once
)

func initAPIPools() {
	poolInitOnce.Do(func() {
		btcPool = NewAPIPool([]*APIEndpoint{
			{URL: "https://mempool.space/api/address/", Name: "Mempool.space"},
			{URL: "https://blockchain.info/rawaddr/", Name: "Blockchain.info"},
			{URL: "https://api.blockchair.com/bitcoin/dashboards/address/", Name: "Blockchair"},
			{URL: "https://api.blockcypher.com/v1/btc/main/addrs/", Name: "BlockCypher"},
		})

		bchPool = NewAPIPool([]*APIEndpoint{
			{URL: "https://api.blockchain.info/bch/multiaddr?active=", Name: "Blockchain.info BCH"},
			{URL: "https://api.blockchair.com/bitcoin-cash/dashboards/address/", Name: "Blockchair BCH"},
		})

		trxPool = NewAPIPool([]*APIEndpoint{
			{URL: "https://api.trongrid.io/v1/accounts/", Name: "TronGrid"},
			{URL: "https://apilist.tronscanapi.com/api/accountv2?address=", Name: "TronScan"},
		})

		ltcPool = NewAPIPool([]*APIEndpoint{
			{URL: "https://litecoinspace.org/api/address/", Name: "Litecoinspace"},
			{URL: "https://api.blockchair.com/litecoin/dashboards/address/", Name: "Blockchair LTC"},
			{URL: "https://api.blockcypher.com/v1/ltc/main/addrs/", Name: "BlockCypher LTC"},
		})

		dogePool = NewAPIPool([]*APIEndpoint{
			{URL: "https://api.blockcypher.com/v1/doge/main/addrs/", Name: "BlockCypher DOGE"},
			{URL: "https://api.blockchair.com/dogecoin/dashboards/address/", Name: "Blockchair DOGE"},
			{URL: "https://dogechain.info/api/v1/address/balance/", Name: "Dogechain"},
		})
	})
}

// ============================================================================
// HTTP CLIENT COMPARTILHADO
// ============================================================================

var httpClient = &http.Client{
	Timeout: 12 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// Cliente com timeout curto para RPCs
var rpcClient = &http.Client{
	Timeout: 8 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	},
}

// httpGetWithTronKey faz GET com header TRON-PRO-API-KEY se TronGrid key estiver configurada
func httpGetWithTronKey(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	// Adicionar API key do TronGrid se configurada
	if apiKeys.TrongridKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", apiKeys.TrongridKey)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func httpGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func httpPost(url string, contentType string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// rpcPost usa o cliente com timeout curto para RPCs EVM
func rpcPost(url string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := rpcClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// ============================================================================
// CONSULTA DE SALDO - EVM VIA RPC PUBLICO (SEM API KEY!)
// ============================================================================

type EVMBalanceResult struct {
	NativeBalance string
	NativeSymbol  string
	Tokens        []TokenResult
	HasBalance    bool
	HasHistory    bool
	TxCount       int
	LastTxDate    string
}

func CheckEVMBalance(address string, network string) (*EVMBalanceResult, error) {
	result := &EVMBalanceResult{
		NativeBalance: "0",
		Tokens:        []TokenResult{},
	}

	nativeSymbols := map[string]string{
		"ethereum": "ETH", "bsc": "BNB", "polygon": "POL",
		"arbitrum": "ETH", "avalanche": "AVAX", "optimism": "ETH", "base": "ETH",
		"sonic": "S", "mantle": "MNT", "flare": "FLR",
		"linea": "ETH", "scroll": "ETH", "gnosis": "xDAI",
		"zksync": "ETH", "blast": "ETH", "cronos": "CRO",
		"celo": "CELO", "berachain": "BERA",
	}
	result.NativeSymbol = nativeSymbols[network]

	rpcs := getEVMRPCs(network)
	if len(rpcs) == 0 {
		return nil, fmt.Errorf("rede EVM nao suportada: %s", network)
	}

	// Tentar cada RPC endpoint (Alchemy primeiro se configurada, depois publicos)
	for _, rpcURL := range rpcs {
		success := false

		// 1. Consultar saldo nativo via eth_getBalance
		balPayload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["%s","latest"],"id":1}`, address)
		body, err := rpcPost(rpcURL, []byte(balPayload))
		if err != nil {
			continue
		}

		var balResp struct {
			Result string      `json:"result"`
			Error  interface{} `json:"error"`
		}
		if json.Unmarshal(body, &balResp) != nil || balResp.Error != nil || balResp.Result == "" {
			continue
		}

		// Parse hex balance
		balHex := strings.TrimPrefix(balResp.Result, "0x")
		if balHex == "" {
			balHex = "0"
		}
		bal := new(big.Int)
		bal.SetString(balHex, 16)

		if bal.Cmp(big.NewInt(0)) > 0 {
			result.NativeBalance = formatWeiBigInt(bal, 18)
			result.HasBalance = true
			result.HasHistory = true
		}

		// 2. Verificar historico via eth_getTransactionCount
		noncePayload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["%s","latest"],"id":2}`, address)
		nonceBody, nonceErr := rpcPost(rpcURL, []byte(noncePayload))
		if nonceErr == nil {
			var nonceResp struct {
				Result string `json:"result"`
			}
			if json.Unmarshal(nonceBody, &nonceResp) == nil && nonceResp.Result != "" {
				nonceHex := strings.TrimPrefix(nonceResp.Result, "0x")
				nonce := new(big.Int)
				nonce.SetString(nonceHex, 16)
				if nonce.Cmp(big.NewInt(0)) > 0 {
					result.HasHistory = true
					result.TxCount = int(nonce.Int64())
				}
			}
		}

		// 3. Verificar tokens ERC-20 via eth_call (balanceOf)
		tokens := getEVMTokens()[network]
		// balanceOf(address) selector = 0x70a08231
		addrClean := strings.TrimPrefix(strings.ToLower(address), "0x")
		addrPadded := fmt.Sprintf("000000000000000000000000%s", addrClean)
		callData := "0x70a08231" + addrPadded

		for _, token := range tokens {
			tokenPayload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_call","params":[{"to":"%s","data":"%s"},"latest"],"id":3}`,
				token.Contract, callData)
			tokenBody, tokenErr := rpcPost(rpcURL, []byte(tokenPayload))
			if tokenErr != nil {
				continue
			}

			var tokenResp struct {
				Result string      `json:"result"`
				Error  interface{} `json:"error"`
			}
			if json.Unmarshal(tokenBody, &tokenResp) != nil || tokenResp.Error != nil {
				continue
			}

			tokenResult := tokenResp.Result
			if tokenResult == "" || tokenResult == "0x" || tokenResult == "0x0" {
				continue
			}

			tokenHex := strings.TrimPrefix(tokenResult, "0x")
			if tokenHex == "" {
				continue
			}

			// Decodificar hex para bytes e verificar se tem conteudo
			tokenBytes, decErr := hex.DecodeString(tokenHex)
			if decErr != nil {
				continue
			}

			tokenBal := new(big.Int)
			tokenBal.SetBytes(tokenBytes)

			if tokenBal.Cmp(big.NewInt(0)) > 0 {
				result.Tokens = append(result.Tokens, TokenResult{
					Symbol:   token.Symbol,
					Name:     token.Name,
					Balance:  formatWeiBigInt(tokenBal, token.Decimals),
					Contract: token.Contract,
				})
				result.HasBalance = true
			}
		}

		success = true
		if success {
			break
		}
	}

	// 4. Descoberta automatica de tokens via Blockscout (encontra QUALQUER token ERC-20)
	// Isso complementa os tokens fixos - se Blockscout encontrar tokens extras, eles sao adicionados
	blockscoutTokens := fetchBlockscoutTokens(address, network)
	for _, bt := range blockscoutTokens {
		// Verificar se este token ja foi encontrado pelos tokens fixos
		alreadyFound := false
		for _, existing := range result.Tokens {
			if strings.EqualFold(existing.Contract, bt.Contract) {
				alreadyFound = true
				break
			}
		}
		if !alreadyFound {
			result.Tokens = append(result.Tokens, bt)
			result.HasBalance = true
		}
	}

	return result, nil
}

// ============================================================================
// BLOCKSCOUT - DESCOBERTA AUTOMATICA DE TOKENS ERC-20
// ============================================================================

var blockscoutURLs = map[string]string{
	"ethereum":  "https://eth.blockscout.com",
	"polygon":   "https://polygon.blockscout.com",
	"optimism":  "https://optimism.blockscout.com",
	"base":      "https://base.blockscout.com",
	"gnosis":    "https://gnosis.blockscout.com",
	"scroll":    "https://scroll.blockscout.com",
	"linea":     "https://linea.blockscout.com",
	"zksync":    "https://zksync.blockscout.com",
	"celo":      "https://celo.blockscout.com",
}

func fetchBlockscoutTokens(address string, network string) []TokenResult {
	var results []TokenResult

	baseURL, ok := blockscoutURLs[network]
	if !ok {
		return results // Rede sem Blockscout, retorna vazio
	}

	url := fmt.Sprintf("%s/api/v2/addresses/%s/token-balances", baseURL, address)
	body, err := httpGet(url)
	if err != nil {
		return results
	}

	var tokenBalances []struct {
		Value string `json:"value"`
		Token struct {
			Address  string `json:"address"`
			Symbol   string `json:"symbol"`
			Name     string `json:"name"`
			Decimals string `json:"decimals"`
			Type     string `json:"type"`
		} `json:"token"`
	}

	if json.Unmarshal(body, &tokenBalances) != nil {
		return results
	}

	for _, tb := range tokenBalances {
		// Apenas tokens ERC-20 com saldo > 0
		if tb.Token.Type != "ERC-20" {
			continue
		}
		if tb.Value == "" || tb.Value == "0" {
			continue
		}

		bal := new(big.Int)
		bal.SetString(tb.Value, 10)
		if bal.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		decimals := 18
		if tb.Token.Decimals != "" {
			fmt.Sscanf(tb.Token.Decimals, "%d", &decimals)
		}

		symbol := tb.Token.Symbol
		if symbol == "" {
			symbol = tb.Token.Address[:10] + "..."
		}

		results = append(results, TokenResult{
			Symbol:   symbol,
			Name:     tb.Token.Name,
			Balance:  formatWei(tb.Value, decimals),
			Contract: tb.Token.Address,
		})
	}

	return results
}

// formatWeiBigInt formata um big.Int de wei para a unidade com decimais
func formatWeiBigInt(wei *big.Int, decimals int) string {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	intPart := new(big.Int).Div(wei, divisor)
	remainder := new(big.Int).Mod(wei, divisor)
	if remainder.Cmp(big.NewInt(0)) == 0 {
		return intPart.String()
	}
	fracStr := fmt.Sprintf("%0*s", decimals, remainder.String())
	fracStr = strings.TrimRight(fracStr, "0")
	if fracStr == "" {
		return intPart.String()
	}
	return fmt.Sprintf("%s.%s", intPart.String(), fracStr)
}

// ============================================================================
// CONSULTA DE SALDO - BTC
// ============================================================================

type BTCBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
	LastTxDate string
}

func CheckBTCBalance(address string) (*BTCBalanceResult, error) {
	initAPIPools()
	result := &BTCBalanceResult{Balance: "0"}

	endpoints := btcPool.GetAll()
	for _, ep := range endpoints {
		success := false

		if strings.Contains(ep.URL, "mempool.space") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				ChainStats struct {
					FundedTxoSum int64 `json:"funded_txo_sum"`
					SpentTxoSum  int64 `json:"spent_txo_sum"`
					TxCount      int   `json:"tx_count"`
				} `json:"chain_stats"`
			}
			if json.Unmarshal(body, &resp) == nil {
				balance := resp.ChainStats.FundedTxoSum - resp.ChainStats.SpentTxoSum
				if balance > 0 {
					result.Balance = formatSatoshi(balance)
					result.HasBalance = true
				}
				if resp.ChainStats.TxCount > 0 {
					result.HasHistory = true
					result.TxCount = resp.ChainStats.TxCount
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "blockchain.info") {
			url := ep.URL + address + "?limit=1"
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				FinalBalance int64 `json:"final_balance"`
				NTx          int   `json:"n_tx"`
				Txs          []struct {
					Time int64 `json:"time"`
				} `json:"txs"`
			}
			if json.Unmarshal(body, &resp) == nil {
				if resp.FinalBalance > 0 {
					result.Balance = formatSatoshi(resp.FinalBalance)
					result.HasBalance = true
				}
				if resp.NTx > 0 {
					result.HasHistory = true
					result.TxCount = resp.NTx
					if len(resp.Txs) > 0 && resp.Txs[0].Time > 0 {
						result.LastTxDate = time.Unix(resp.Txs[0].Time, 0).Format("2006-01-02")
					}
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "blockchair") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Data map[string]struct {
					Address struct {
						Balance  int64  `json:"balance"`
						TxCount  int    `json:"transaction_count"`
						LastSeen string `json:"last_seen_receiving"`
					} `json:"address"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &resp) == nil {
				for _, v := range resp.Data {
					if v.Address.Balance > 0 {
						result.Balance = formatSatoshi(v.Address.Balance)
						result.HasBalance = true
					}
					if v.Address.TxCount > 0 {
						result.HasHistory = true
						result.TxCount = v.Address.TxCount
					}
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "blockcypher") {
			url := ep.URL + address + "/balance"
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Balance   int64 `json:"balance"`
				TotalSent int64 `json:"total_sent"`
				NTx       int   `json:"n_tx"`
			}
			if json.Unmarshal(body, &resp) == nil {
				if resp.Balance > 0 {
					result.Balance = formatSatoshi(resp.Balance)
					result.HasBalance = true
				}
				if resp.NTx > 0 {
					result.HasHistory = true
					result.TxCount = resp.NTx
				}
				success = true
			}
		}

		if success {
			break
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - BCH (blockchain.info/bch como primario)
// ============================================================================

func CheckBCHBalance(address string) (*BTCBalanceResult, error) {
	initAPIPools()
	result := &BTCBalanceResult{Balance: "0"}

	endpoints := bchPool.GetAll()
	for _, ep := range endpoints {
		success := false

		if strings.Contains(ep.URL, "blockchain.info") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Addresses []struct {
					FinalBalance  int64 `json:"final_balance"`
					NTx           int   `json:"n_tx"`
					TotalReceived int64 `json:"total_received"`
				} `json:"addresses"`
				Txs []struct {
					Time int64 `json:"time"`
				} `json:"txs"`
			}
			if json.Unmarshal(body, &resp) == nil && len(resp.Addresses) > 0 {
				addr := resp.Addresses[0]
				if addr.FinalBalance > 0 {
					result.Balance = formatSatoshi(addr.FinalBalance)
					result.HasBalance = true
				}
				if addr.NTx > 0 || addr.TotalReceived > 0 {
					result.HasHistory = true
					result.TxCount = addr.NTx
					if len(resp.Txs) > 0 && resp.Txs[0].Time > 0 {
						result.LastTxDate = time.Unix(resp.Txs[0].Time, 0).Format("2006-01-02")
					}
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "blockchair") {
			shortAddr := address
			if strings.Contains(address, ":") {
				parts := strings.SplitN(address, ":", 2)
				shortAddr = parts[1]
			}
			url := ep.URL + shortAddr
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Data map[string]struct {
					Address struct {
						Balance int64 `json:"balance"`
						TxCount int   `json:"transaction_count"`
					} `json:"address"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &resp) == nil {
				for _, v := range resp.Data {
					if v.Address.Balance > 0 {
						result.Balance = formatSatoshi(v.Address.Balance)
						result.HasBalance = true
					}
					if v.Address.TxCount > 0 {
						result.HasHistory = true
						result.TxCount = v.Address.TxCount
					}
				}
				success = true
			}
		}

		if success {
			break
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - TRX
// ============================================================================

type TRXBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
	LastTxDate string
	Tokens     []TokenResult
}

func CheckTRXBalance(address string) (*TRXBalanceResult, error) {
	initAPIPools()
	result := &TRXBalanceResult{Balance: "0", Tokens: []TokenResult{}}

	endpoints := trxPool.GetAll()
	for _, ep := range endpoints {
		success := false

		if strings.Contains(ep.URL, "trongrid") {
			url := ep.URL + address
			body, err := httpGetWithTronKey(url)
			if err != nil {
				continue
			}
			var resp struct {
				Data []struct {
					Balance int64               `json:"balance"`
					Trc20   []map[string]string `json:"trc20"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &resp) == nil && len(resp.Data) > 0 {
				if resp.Data[0].Balance > 0 {
					result.Balance = formatSunToTRX(resp.Data[0].Balance)
					result.HasBalance = true
					result.HasHistory = true
				}
				trxTokens := getTRXTokens()
				for _, tokenMap := range resp.Data[0].Trc20 {
					for contract, balance := range tokenMap {
						bal := new(big.Int)
						bal.SetString(balance, 10)
						if bal.Cmp(big.NewInt(0)) > 0 {
							symbol := contract[:8] + "..."
							name := contract
							for _, t := range trxTokens {
								if strings.EqualFold(t.Contract, contract) {
									symbol = t.Symbol
									name = t.Name
									balance = formatWei(balance, t.Decimals)
									break
								}
							}
							result.Tokens = append(result.Tokens, TokenResult{
								Symbol:   symbol,
								Name:     name,
								Balance:  balance,
								Contract: contract,
							})
							result.HasBalance = true
						}
					}
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "tronscan") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Balance               int64 `json:"balance"`
				TotalTransactionCount int   `json:"totalTransactionCount"`
				WithPriceTokens       []struct {
					TokenAbbr string `json:"tokenAbbr"`
					TokenName string `json:"tokenName"`
					Balance   string `json:"balance"`
					TokenId   string `json:"tokenId"`
				} `json:"withPriceTokens"`
			}
			if json.Unmarshal(body, &resp) == nil {
				if resp.Balance > 0 {
					result.Balance = formatSunToTRX(resp.Balance)
					result.HasBalance = true
				}
				if resp.TotalTransactionCount > 0 {
					result.HasHistory = true
					result.TxCount = resp.TotalTransactionCount
				}
				for _, t := range resp.WithPriceTokens {
					if t.TokenAbbr != "trx" && t.Balance != "0" && t.Balance != "" {
						result.Tokens = append(result.Tokens, TokenResult{
							Symbol:   t.TokenAbbr,
							Name:     t.TokenName,
							Balance:  t.Balance,
							Contract: t.TokenId,
						})
						result.HasBalance = true
					}
				}
				success = true
			}
		}

		if success {
			break
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - SOLANA (SOL + SPL Tokens)
// ============================================================================

type SOLBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
	Tokens     []TokenResult
}

func CheckSOLBalance(address string) (*SOLBalanceResult, error) {
	result := &SOLBalanceResult{Balance: "0", Tokens: []TokenResult{}}

	rpcEndpoints := []string{}

	// Se Alchemy configurada, adicionar como RPC prioritario para Solana
	if apiKeys.AlchemyKey != "" {
		rpcEndpoints = append(rpcEndpoints, fmt.Sprintf("https://solana-mainnet.g.alchemy.com/v2/%s", apiKeys.AlchemyKey))
	}

	// RPCs publicos como fallback
	rpcEndpoints = append(rpcEndpoints,
		"https://api.mainnet-beta.solana.com",
		"https://solana-mainnet.g.alchemy.com/v2/demo",
	)

	for _, rpcURL := range rpcEndpoints {
		balPayload := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"getBalance","params":["%s"]}`, address)
		body, err := httpPost(rpcURL, "application/json", []byte(balPayload))
		if err != nil {
			continue
		}

		var balResp struct {
			Result struct {
				Value int64 `json:"value"`
			} `json:"result"`
			Error interface{} `json:"error"`
		}
		if json.Unmarshal(body, &balResp) != nil || balResp.Error != nil {
			continue
		}

		if balResp.Result.Value > 0 {
			lamports := balResp.Result.Value
			sol := float64(lamports) / 1000000000.0
			result.Balance = fmt.Sprintf("%.9f", sol)
			result.HasBalance = true
			result.HasHistory = true
		}

		sigPayload := fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"getSignaturesForAddress","params":["%s",{"limit":1}]}`, address)
		sigBody, sigErr := httpPost(rpcURL, "application/json", []byte(sigPayload))
		if sigErr == nil {
			var sigResp struct {
				Result []struct {
					BlockTime int64 `json:"blockTime"`
				} `json:"result"`
			}
			if json.Unmarshal(sigBody, &sigResp) == nil && len(sigResp.Result) > 0 {
				result.HasHistory = true
				result.TxCount = 1
			}
		}

		tokenPayload := fmt.Sprintf(`{"jsonrpc":"2.0","id":3,"method":"getTokenAccountsByOwner","params":["%s",{"programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},{"encoding":"jsonParsed"}]}`, address)
		tokenBody, tokenErr := httpPost(rpcURL, "application/json", []byte(tokenPayload))
		if tokenErr == nil {
			var tokenResp struct {
				Result struct {
					Value []struct {
						Account struct {
							Data struct {
								Parsed struct {
									Info struct {
										Mint        string `json:"mint"`
										TokenAmount struct {
											UiAmountString string  `json:"uiAmountString"`
											UiAmount       float64 `json:"uiAmount"`
										} `json:"tokenAmount"`
									} `json:"info"`
								} `json:"parsed"`
							} `json:"data"`
						} `json:"account"`
					} `json:"value"`
				} `json:"result"`
			}
			if json.Unmarshal(tokenBody, &tokenResp) == nil {
				solTokens := getSOLTokens()
				for _, acct := range tokenResp.Result.Value {
					info := acct.Account.Data.Parsed.Info
					if info.TokenAmount.UiAmount > 0 {
						mint := info.Mint
						symbol := mint[:8] + "..."
						name := mint
						for _, t := range solTokens {
							if t.Contract == mint {
								symbol = t.Symbol
								name = t.Name
								break
							}
						}
						result.Tokens = append(result.Tokens, TokenResult{
							Symbol:   symbol,
							Name:     name,
							Balance:  info.TokenAmount.UiAmountString,
							Contract: mint,
						})
						result.HasBalance = true
					}
				}
			}
		}

		return result, nil
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - TON (Toncoin)
// ============================================================================

type TONBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
}

func CheckTONBalance(pubKeyHex string) (*TONBalanceResult, error) {
	result := &TONBalanceResult{Balance: "0"}

	url := fmt.Sprintf("https://toncenter.com/api/v2/getWalletInformation?address=%s", pubKeyHex)
	body, err := httpGet(url)
	if err != nil {
		url = fmt.Sprintf("https://tonapi.io/v2/accounts/%s", pubKeyHex)
		body, err = httpGet(url)
		if err != nil {
			return result, nil
		}
		var resp struct {
			Balance int64  `json:"balance"`
			Status  string `json:"status"`
		}
		if json.Unmarshal(body, &resp) == nil {
			if resp.Balance > 0 {
				ton := float64(resp.Balance) / 1000000000.0
				result.Balance = fmt.Sprintf("%.9f", ton)
				result.HasBalance = true
				result.HasHistory = true
			}
			if resp.Status == "active" {
				result.HasHistory = true
			}
		}
		return result, nil
	}

	var resp struct {
		Ok     bool `json:"ok"`
		Result struct {
			Balance      string `json:"balance"`
			AccountState string `json:"account_state"`
			LastTxLt     string `json:"last_transaction_id"`
		} `json:"result"`
	}
	if json.Unmarshal(body, &resp) == nil && resp.Ok {
		bal := new(big.Int)
		bal.SetString(resp.Result.Balance, 10)
		if bal.Cmp(big.NewInt(0)) > 0 {
			ton := float64(bal.Int64()) / 1000000000.0
			result.Balance = fmt.Sprintf("%.9f", ton)
			result.HasBalance = true
		}
		if resp.Result.AccountState == "active" {
			result.HasHistory = true
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - ZCASH (ZEC Transparent)
// ============================================================================

func CheckZECBalance(address string) (*BTCBalanceResult, error) {
	result := &BTCBalanceResult{Balance: "0"}

	urls := []string{
		fmt.Sprintf("https://api.blockchair.com/zcash/dashboards/address/%s", address),
		fmt.Sprintf("https://zcashblockexplorer.com/api/addr/%s", address),
	}

	for _, url := range urls {
		body, err := httpGet(url)
		if err != nil {
			continue
		}

		if strings.Contains(url, "blockchair") {
			var resp struct {
				Data map[string]struct {
					Address struct {
						Balance int64 `json:"balance"`
						TxCount int   `json:"transaction_count"`
					} `json:"address"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &resp) == nil {
				for _, v := range resp.Data {
					if v.Address.Balance > 0 {
						result.Balance = formatSatoshi(v.Address.Balance)
						result.HasBalance = true
					}
					if v.Address.TxCount > 0 {
						result.HasHistory = true
						result.TxCount = v.Address.TxCount
					}
				}
				return result, nil
			}
		} else {
			var resp struct {
				Balance       float64 `json:"balance"`
				TotalReceived float64 `json:"totalReceived"`
				TxApperances  int     `json:"txApperances"`
			}
			if json.Unmarshal(body, &resp) == nil {
				if resp.Balance > 0 {
					result.Balance = fmt.Sprintf("%.8f", resp.Balance)
					result.HasBalance = true
				}
				if resp.TxApperances > 0 || resp.TotalReceived > 0 {
					result.HasHistory = true
					result.TxCount = resp.TxApperances
				}
				return result, nil
			}
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - LTC
// ============================================================================

func CheckLTCBalance(address string) (*BTCBalanceResult, error) {
	initAPIPools()
	result := &BTCBalanceResult{Balance: "0"}

	endpoints := ltcPool.GetAll()
	for _, ep := range endpoints {
		success := false

		if strings.Contains(ep.URL, "litecoinspace") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				ChainStats struct {
					FundedTxoSum int64 `json:"funded_txo_sum"`
					SpentTxoSum  int64 `json:"spent_txo_sum"`
					TxCount      int   `json:"tx_count"`
				} `json:"chain_stats"`
			}
			if json.Unmarshal(body, &resp) == nil {
				balance := resp.ChainStats.FundedTxoSum - resp.ChainStats.SpentTxoSum
				if balance > 0 {
					result.Balance = formatSatoshi(balance)
					result.HasBalance = true
				}
				if resp.ChainStats.TxCount > 0 {
					result.HasHistory = true
					result.TxCount = resp.ChainStats.TxCount
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "blockchair") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Data map[string]struct {
					Address struct {
						Balance int64 `json:"balance"`
						TxCount int   `json:"transaction_count"`
					} `json:"address"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &resp) == nil {
				for _, v := range resp.Data {
					if v.Address.Balance > 0 {
						result.Balance = formatSatoshi(v.Address.Balance)
						result.HasBalance = true
					}
					if v.Address.TxCount > 0 {
						result.HasHistory = true
						result.TxCount = v.Address.TxCount
					}
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "blockcypher") {
			url := ep.URL + address + "/balance"
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Balance   int64 `json:"balance"`
				TotalSent int64 `json:"total_sent"`
				NTx       int   `json:"n_tx"`
			}
			if json.Unmarshal(body, &resp) == nil {
				if resp.Balance > 0 {
					result.Balance = formatSatoshi(resp.Balance)
					result.HasBalance = true
				}
				if resp.NTx > 0 {
					result.HasHistory = true
					result.TxCount = resp.NTx
				}
				success = true
			}
		}

		if success {
			break
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - DOGE
// ============================================================================

func CheckDOGEBalance(address string) (*BTCBalanceResult, error) {
	initAPIPools()
	result := &BTCBalanceResult{Balance: "0"}

	endpoints := dogePool.GetAll()
	for _, ep := range endpoints {
		success := false

		if strings.Contains(ep.URL, "blockcypher") {
			url := ep.URL + address + "/balance"
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Balance   int64 `json:"balance"`
				TotalSent int64 `json:"total_sent"`
				NTx       int   `json:"n_tx"`
			}
			if json.Unmarshal(body, &resp) == nil {
				if resp.Balance > 0 {
					result.Balance = formatSatoshi(resp.Balance)
					result.HasBalance = true
				}
				if resp.NTx > 0 {
					result.HasHistory = true
					result.TxCount = resp.NTx
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "dogechain") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Balance string `json:"balance"`
				Success int    `json:"success"`
			}
			if json.Unmarshal(body, &resp) == nil && resp.Success == 1 {
				if resp.Balance != "0" && resp.Balance != "0.00000000" {
					result.Balance = resp.Balance
					result.HasBalance = true
					result.HasHistory = true
				}
				success = true
			}
		} else if strings.Contains(ep.URL, "blockchair") {
			url := ep.URL + address
			body, err := httpGet(url)
			if err != nil {
				continue
			}
			var resp struct {
				Data map[string]struct {
					Address struct {
						Balance int64 `json:"balance"`
						TxCount int   `json:"transaction_count"`
					} `json:"address"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &resp) == nil {
				for _, v := range resp.Data {
					if v.Address.Balance > 0 {
						result.Balance = formatSatoshi(v.Address.Balance)
						result.HasBalance = true
					}
					if v.Address.TxCount > 0 {
						result.HasHistory = true
						result.TxCount = v.Address.TxCount
					}
				}
				success = true
			}
		}

		if success {
			break
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - XRP (Ripple)
// ============================================================================

type XRPBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
}

func CheckXRPBalance(address string) (*XRPBalanceResult, error) {
	result := &XRPBalanceResult{Balance: "0"}

	urls := []string{
		"https://xrplcluster.com",
		"https://s1.ripple.com:51234",
		"https://s2.ripple.com:51234",
	}

	for _, rpcURL := range urls {
		payload := fmt.Sprintf(`{"method":"account_info","params":[{"account":"%s","ledger_index":"validated"}]}`, address)
		body, err := httpPost(rpcURL, "application/json", []byte(payload))
		if err != nil {
			continue
		}

		var resp struct {
			Result struct {
				AccountData struct {
					Balance  string `json:"Balance"`
					Sequence int    `json:"Sequence"`
				} `json:"account_data"`
				Status string `json:"status"`
				Error  string `json:"error"`
			} `json:"result"`
		}
		if json.Unmarshal(body, &resp) != nil {
			continue
		}

		if resp.Result.Status == "success" {
			// XRP balance is in drops (1 XRP = 1,000,000 drops)
			// Account reserve is 10 XRP, so subtract it
			bal := new(big.Int)
			bal.SetString(resp.Result.AccountData.Balance, 10)
			if bal.Cmp(big.NewInt(0)) > 0 {
				xrp := float64(bal.Int64()) / 1000000.0
				result.Balance = fmt.Sprintf("%.6f", xrp)
				result.HasBalance = true
				result.HasHistory = true
			}
			if resp.Result.AccountData.Sequence > 0 {
				result.HasHistory = true
				result.TxCount = resp.Result.AccountData.Sequence
			}
			return result, nil
		} else if resp.Result.Error == "actNotFound" {
			// Account not found = no balance, no history
			return result, nil
		}
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - STELLAR (XLM)
// ============================================================================

type XLMBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
	Tokens     []TokenResult
}

func CheckXLMBalance(address string) (*XLMBalanceResult, error) {
	result := &XLMBalanceResult{Balance: "0", Tokens: []TokenResult{}}

	url := fmt.Sprintf("https://horizon.stellar.org/accounts/%s", address)
	body, err := httpGet(url)
	if err != nil {
		// Account not found (404) = no balance
		return result, nil
	}

	var resp struct {
		Sequence string `json:"sequence"`
		Balances []struct {
			Balance   string `json:"balance"`
			AssetType string `json:"asset_type"`
			AssetCode string `json:"asset_code"`
			AssetIssuer string `json:"asset_issuer"`
		} `json:"balances"`
		LastModifiedTime string `json:"last_modified_time"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return result, nil
	}

	for _, b := range resp.Balances {
		if b.AssetType == "native" {
			if b.Balance != "0.0000000" && b.Balance != "0" {
				result.Balance = b.Balance
				result.HasBalance = true
				result.HasHistory = true
			}
		} else {
			if b.Balance != "0.0000000" && b.Balance != "0" {
				result.Tokens = append(result.Tokens, TokenResult{
					Symbol:   b.AssetCode,
					Name:     b.AssetCode,
					Balance:  b.Balance,
					Contract: b.AssetIssuer,
				})
				result.HasBalance = true
			}
		}
	}

	if resp.Sequence != "0" && resp.Sequence != "" {
		result.HasHistory = true
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - ALGORAND (ALGO)
// ============================================================================

type ALGOBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
	Tokens     []TokenResult
}

func CheckALGOBalance(address string) (*ALGOBalanceResult, error) {
	result := &ALGOBalanceResult{Balance: "0", Tokens: []TokenResult{}}

	urls := []string{
		fmt.Sprintf("https://mainnet-api.algonode.cloud/v2/accounts/%s", address),
		fmt.Sprintf("https://mainnet-idx.algonode.cloud/v2/accounts/%s", address),
	}

	for _, url := range urls {
		body, err := httpGet(url)
		if err != nil {
			continue
		}

		var resp struct {
			Amount             uint64 `json:"amount"`
			TotalAppsOptedIn   int    `json:"total-apps-opted-in"`
			TotalAssetsOptedIn int    `json:"total-assets-opted-in"`
			Assets             []struct {
				AssetId uint64 `json:"asset-id"`
				Amount  uint64 `json:"amount"`
			} `json:"assets"`
		}
		if json.Unmarshal(body, &resp) != nil {
			continue
		}

		if resp.Amount > 0 {
			algo := float64(resp.Amount) / 1000000.0
			result.Balance = fmt.Sprintf("%.6f", algo)
			result.HasBalance = true
			result.HasHistory = true
		}

		// Check ASA tokens
		for _, asset := range resp.Assets {
			if asset.Amount > 0 {
				result.Tokens = append(result.Tokens, TokenResult{
					Symbol:   fmt.Sprintf("ASA#%d", asset.AssetId),
					Name:     fmt.Sprintf("Algorand ASA %d", asset.AssetId),
					Balance:  fmt.Sprintf("%d", asset.Amount),
					Contract: fmt.Sprintf("%d", asset.AssetId),
				})
				result.HasBalance = true
			}
		}

		return result, nil
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - SUI
// ============================================================================

type SUIBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
	Tokens     []TokenResult
}

func CheckSUIBalance(address string) (*SUIBalanceResult, error) {
	result := &SUIBalanceResult{Balance: "0", Tokens: []TokenResult{}}

	rpcURLs := []string{
		"https://fullnode.mainnet.sui.io:443",
		"https://sui-mainnet.public.blastapi.io",
	}

	for _, rpcURL := range rpcURLs {
		// Get all coin balances (SUI + tokens)
		payload := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"suix_getAllBalances","params":["%s"]}`, address)
		body, err := httpPost(rpcURL, "application/json", []byte(payload))
		if err != nil {
			continue
		}

		var resp struct {
			Result []struct {
				CoinType     string `json:"coinType"`
				TotalBalance string `json:"totalBalance"`
			} `json:"result"`
			Error interface{} `json:"error"`
		}
		if json.Unmarshal(body, &resp) != nil || resp.Error != nil {
			continue
		}

		for _, coin := range resp.Result {
			bal := new(big.Int)
			bal.SetString(coin.TotalBalance, 10)
			if bal.Cmp(big.NewInt(0)) <= 0 {
				continue
			}

			if coin.CoinType == "0x2::sui::SUI" {
				// SUI has 9 decimals (1 SUI = 10^9 MIST)
				sui := float64(bal.Int64()) / 1000000000.0
				result.Balance = fmt.Sprintf("%.9f", sui)
				result.HasBalance = true
				result.HasHistory = true
			} else {
				// Other tokens
				result.Tokens = append(result.Tokens, TokenResult{
					Symbol:   coin.CoinType,
					Name:     coin.CoinType,
					Balance:  coin.TotalBalance,
					Contract: coin.CoinType,
				})
				result.HasBalance = true
			}
		}

		// Check transaction count
		txPayload := fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"suix_queryTransactionBlocks","params":[{"filter":{"FromAddress":"%s"},"options":{"showInput":false}},null,1,true]}`, address)
		txBody, txErr := httpPost(rpcURL, "application/json", []byte(txPayload))
		if txErr == nil {
			var txResp struct {
				Result struct {
					Data []interface{} `json:"data"`
				} `json:"result"`
			}
			if json.Unmarshal(txBody, &txResp) == nil && len(txResp.Result.Data) > 0 {
				result.HasHistory = true
				result.TxCount = 1
			}
		}

		return result, nil
	}

	return result, nil
}

// ============================================================================
// CONSULTA DE SALDO - NEAR
// ============================================================================

type NEARBalanceResult struct {
	Balance    string
	HasBalance bool
	HasHistory bool
	TxCount    int
}

func CheckNEARBalance(address string) (*NEARBalanceResult, error) {
	result := &NEARBalanceResult{Balance: "0"}

	rpcURLs := []string{
		"https://rpc.mainnet.near.org",
		"https://near.lava.build",
	}

	for _, rpcURL := range rpcURLs {
		payload := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"query","params":{"request_type":"view_account","finality":"final","account_id":"%s"}}`, address)
		body, err := httpPost(rpcURL, "application/json", []byte(payload))
		if err != nil {
			continue
		}

		var resp struct {
			Result struct {
				Amount     string `json:"amount"`
				BlockHash  string `json:"block_hash"`
			} `json:"result"`
			Error interface{} `json:"error"`
		}
		if json.Unmarshal(body, &resp) != nil || resp.Error != nil {
			continue
		}

		if resp.Result.Amount != "" && resp.Result.Amount != "0" {
			// NEAR has 24 decimals (yoctoNEAR)
			bal := new(big.Int)
			bal.SetString(resp.Result.Amount, 10)
			if bal.Cmp(big.NewInt(0)) > 0 {
				// Convert yoctoNEAR to NEAR (24 decimals)
				result.Balance = formatWei(resp.Result.Amount, 24)
				result.HasBalance = true
				result.HasHistory = true
			}
		}

		return result, nil
	}

	return result, nil
}

// ============================================================================
// FUNCOES DE FORMATACAO
// ============================================================================

func formatWei(weiStr string, decimals int) string {
	wei := new(big.Int)
	wei.SetString(weiStr, 10)
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	intPart := new(big.Int).Div(wei, divisor)
	remainder := new(big.Int).Mod(wei, divisor)
	if remainder.Cmp(big.NewInt(0)) == 0 {
		return intPart.String()
	}
	fracStr := fmt.Sprintf("%0*s", decimals, remainder.String())
	fracStr = strings.TrimRight(fracStr, "0")
	if fracStr == "" {
		return intPart.String()
	}
	return fmt.Sprintf("%s.%s", intPart.String(), fracStr)
}

func formatSatoshi(satoshi int64) string {
	btc := float64(satoshi) / 100000000.0
	return fmt.Sprintf("%.8f", btc)
}

func formatSunToTRX(sun int64) string {
	trx := float64(sun) / 1000000.0
	return fmt.Sprintf("%.6f", trx)
}

func parseTimestamp(ts string) (string, error) {
	n := new(big.Int)
	n.SetString(ts, 10)
	t := time.Unix(n.Int64(), 0)
	return t.Format("2006-01-02"), nil
}
