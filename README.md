<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue?style=for-the-badge" />
  <img src="https://img.shields.io/badge/License-Open%20Source-green?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Networks-31%2B%20Blockchains-orange?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Tokens-160%2B-purple?style=for-the-badge" />
</p>

<h1 align="center">CRYPTO HUNTER PRO</h1>

<pre align="center">
 ██████╗██████╗ ██╗   ██╗██████╗ ████████╗ ██████╗ 
██╔════╝██╔══██╗╚██╗ ██╔╝██╔══██╗╚══██╔══╝██╔═══██╗
██║     ██████╔╝ ╚████╔╝ ██████╔╝   ██║   ██║   ██║
██║     ██╔══██╗  ╚██╔╝  ██╔═══╝    ██║   ██║   ██║
╚██████╗██║  ██║   ██║   ██║        ██║   ╚██████╔╝
 ╚═════╝╚═╝  ╚═╝   ╚═╝   ╚═╝        ╚═╝    ╚═════╝ 
██╗  ██╗██╗   ██╗███╗   ██╗████████╗███████╗██████╗ 
██║  ██║██║   ██║████╗  ██║╚══██╔══╝██╔════╝██╔══██╗
███████║██║   ██║██╔██╗ ██║   ██║   █████╗  ██████╔╝
██╔══██║██║   ██║██║╚██╗██║   ██║   ██╔══╝  ██╔══██╗
██║  ██║╚██████╔╝██║ ╚████║   ██║   ███████╗██║  ██║
╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
██████╗ ██████╗  ██████╗ 
██╔══██╗██╔══██╗██╔═══██╗
██████╔╝██████╔╝██║   ██║
██╔═══╝ ██╔══██╗██║   ██║
██║     ██║  ██║╚██████╔╝
╚═╝     ╚═╝  ╚═╝ ╚═════╝ 
</pre>

<p align="center">
  <img src="logo.png" alt="Crypto Hunter Pro" width="200"/>
</p>

<p align="center">
  <strong>Este README esta disponivel em:</strong><br/>
  <a href="#portugues-brasil">Portugues :brazil:</a> · <a href="#english">English :us:</a>
</p>

---

---

<a name="portugues-brasil"></a>

# :brazil: PORTUGUES

---

## O que e o Crypto Hunter Pro?

**Crypto Hunter Pro** e um conjunto de ferramentas open-source projetado para ajudar detentores de criptomoedas a recuperar seed phrases perdidas ou esquecidas e escanear multiplas blockchains em busca de fundos. Construido em Go para performance maxima, combina dois modulos poderosos que funcionam juntos como um pipeline completo de recuperacao.

> Perdeu a ordem da sua seed phrase? Esqueceu qual carteira usou? Tem palavras parciais? O Crypto Hunter Pro foi feito para voce.

---

## Modulos

### Modulo 1: Unmixer Seed Search

O **Unmixer Seed Search** e um motor de permutacao de seed phrases BIP39 que gera todas as combinacoes validas de seeds a partir de informacoes parciais ou embaralhadas.

**Funcionalidades Principais:**

| Funcionalidade | Descricao |
|----------------|-----------|
| **4 Modos de Input** | Simples (ordem conhecida), Avancado Parcial (wildcards + ordem desconhecida), Avancado Completo (palavras completas + ordem desconhecida), Descrambler (testa TODAS as permutacoes) |
| **Suporte a Wildcards** | Use `*` para partes desconhecidas: `aban*`, `bo*`, `*tion`, `*` (palavra inteira desconhecida) |
| **9 Idiomas BIP39** | Ingles, Espanhol, Frances, Italiano, Portugues, Japones, Coreano, Chines Simplificado, Chines Tradicional |
| **Correcao Inteligente de Typos** | Sugestoes automaticas quando voce digita uma palavra errada (distancia de Levenshtein + mapa de teclas adjacentes + deteccao de transposicao) |
| **Validacao de Checksum BIP39** | Apenas gera seeds que passam na verificacao criptografica de checksum |
| **Exportacao Excel** | Resultados salvos em arquivos `.xlsx` organizados, prontos para importacao no CIE |
| **Interface Bilingue** | Suporte completo em Portugues e Ingles |

**Como Funciona:**

```
Suas palavras embaralhadas/parciais → Unmixer → Todas as combinacoes BIP39 validas → Arquivo Excel
```

---

### Modulo 2: CIE - Crypto Intelligence Engine

O **CIE (Crypto Intelligence Engine)** e um scanner multi-chain de carteiras que recebe seed phrases e automaticamente verifica saldos em 31+ blockchains e 160+ tokens.

**Funcionalidades Principais:**

| Funcionalidade | Descricao |
|----------------|-----------|
| **31 Blockchains** | Bitcoin, Ethereum, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Solana, Tron, Litecoin, Dogecoin, TON, Zcash, XRP, Stellar, Algorand, Sui, Near, e mais |
| **160+ Tokens** | Verificacao automatica de saldos ERC-20, BEP-20, SPL, TRC-20 |
| **Multi-Derivation Paths** | Testa TODOS os formatos de endereco automaticamente (Legacy, SegWit, Native SegWit, Taproot) |
| **Pool de APIs Multiplas** | Sistema de failover com multiplos provedores por rede para maxima confiabilidade |
| **Importacao Excel** | Importacao direta dos resultados do Unmixer Seed Search |
| **Exibicao em Tempo Real** | Acompanhe os enderecos sendo verificados ao vivo com indicadores de saldo |
| **Relatorios Detalhados** | Saida em Excel com abas Resumo, Com Saldo e Com Historico |
| **Suporte a Passphrase** | Teste de passphrase BIP39 opcional (25a palavra) |

**Redes Suportadas:**

```
EVM: Ethereum, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Linea,
     Scroll, Gnosis, zkSync, Blast, Cronos, Celo, Berachain, Sonic, Mantle, Flare

UTXO: Bitcoin (4 tipos de endereco), Bitcoin Cash, Litecoin (3 tipos), Dogecoin, Zcash

Outras: Solana, Tron, TON, XRP, Stellar, Algorand, Sui, Near Protocol
```

---

## O Pipeline: Unmixer → CIE

Os dois modulos funcionam juntos como um pipeline completo de recuperacao:

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  1. UNMIXER SEED SEARCH                                         │
│     Entrada: Suas palavras embaralhadas/parciais                │
│     Saida: Todas as combinacoes BIP39 validas (Excel)           │
│                                                                 │
│                          ↓                                      │
│                                                                 │
│  2. CIE - CRYPTO INTELLIGENCE ENGINE                            │
│     Entrada: Seeds validas do Unmixer (ou entrada manual)       │
│     Saida: Quais seeds tem fundos, em quais redes,              │
│            com detalhes completos de saldo (relatorio Excel)     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Inicio Rapido

### Requisitos

- **Go 1.21+** (para compilacao)
- **Windows, Linux ou macOS**
- **Conexao com internet** (apenas CIE, para consultas na blockchain)

### Compilacao

```bash
# Clone o repositorio
git clone https://github.com/YOUR_USERNAME/crypto-hunter-pro.git
cd crypto-hunter-pro

# Compilar Unmixer Seed Search
cd CRYPTO_HUNTER_PRO_Unmixer_Seed_Search
go build -o unmixer_seed_search .

# Compilar CIE
cd ../CRYPTO_HUNTER_PRO_CIE
go build -o crypto_hunter_pro_cie .
```

### Windows (usando o batch incluso)

```batch
# Basta dar duplo clique no compilador.bat em cada pasta do modulo
```

### Linux e macOS (via codigo-fonte)

> **Nota:** Os binarios pre-compilados (.exe) sao para Windows. Para Linux e macOS, compile diretamente a partir do codigo-fonte:

```bash
# Linux
cd CRYPTO_HUNTER_PRO_Unmixer_Seed_Search
GOOS=linux GOARCH=amd64 go build -o unmixer_seed_search .

cd ../CRYPTO_HUNTER_PRO_CIE
GOOS=linux GOARCH=amd64 go build -o crypto_hunter_pro_cie .

# macOS
cd CRYPTO_HUNTER_PRO_Unmixer_Seed_Search
GOOS=darwin GOARCH=amd64 go build -o unmixer_seed_search .

cd ../CRYPTO_HUNTER_PRO_CIE
GOOS=darwin GOARCH=amd64 go build -o crypto_hunter_pro_cie .
```

---

## Exemplos de Uso

### Exemplo 1: Modo Descrambler (Modo 4)

Voce tem todas as 12 palavras mas nao sabe a ordem:

```
Selecione o Modo: 4

Digite todas as palavras separadas por espaco:
> abandon ability able about above absent absorb abstract absurd abuse access acid

Processando todas as 479.001.600 permutacoes...
Seeds validas encontradas: 1
```

### Exemplo 2: Recuperacao com Wildcards (Modo 1)

Voce sabe a ordem mas esqueceu algumas palavras:

```
Selecione o Modo: 1

Digite a seed (use * para desconhecidas):
> abandon * able about * absent absorb abstract absurd abuse access acid

Expandindo wildcards contra a wordlist BIP39...
Seeds validas encontradas: 3
```

### Exemplo 3: CIE Scan Multi-Chain

```
Digite a seed phrase: abandon ability able about above absent absorb abstract absurd abuse access acid

Selecione as redes: TODAS
Testar TODOS os derivation paths? SIM

Escaneando 31 redes...
[BTC] bc1q... - Saldo: 0.00142 BTC ✓
[ETH] 0x...  - Saldo: 0.5 ETH ✓
[ETH] 0x...  - USDT: 150.00 ✓
```

---

## Derivation Paths Suportados

| Rede | Formato | Path |
|------|---------|------|
| Bitcoin | Legacy (1...) | m/44'/0'/0'/0/x |
| Bitcoin | SegWit (3...) | m/49'/0'/0'/0/x |
| Bitcoin | Native SegWit (bc1q...) | m/84'/0'/0'/0/x |
| Bitcoin | Taproot (bc1p...) | m/86'/0'/0'/0/x |
| EVM (todas) | Padrao (0x...) | m/44'/60'/0'/0/x |
| Solana | Ed25519 | m/44'/501'/0'/0' |
| Tron | Padrao (T...) | m/44'/195'/0'/0/x |
| Litecoin | Legacy/SegWit/Native | m/44'/2', m/49'/2', m/84'/2' |
| Dogecoin | Padrao (D...) | m/44'/3'/0'/0/x |
| TON | Ed25519 | m/44'/607'/0' |
| XRP | Padrao (r...) | m/44'/144'/0'/0/x |
| Stellar | Ed25519 (G...) | m/44'/148'/0' |

---

## Seguranca

- **100% Offline** (Unmixer): O motor de permutacao funciona completamente offline. Sem internet necessaria, sem dados enviados para lugar nenhum.
- **Open Source**: Cada linha de codigo e auditavel. Sem backdoors escondidos, sem telemetria, sem coleta de dados.
- **Processamento Local**: Todas as operacoes criptograficas acontecem na SUA maquina. Seeds nunca saem do seu computador.
- **Sem Dependencias Externas em Runtime**: O binario compilado e autossuficiente.

> **AVISO**: Nunca compartilhe suas seed phrases com ninguem. Nunca insira suas seed phrases em websites. Sempre verifique que voce esta rodando a versao oficial open-source.

---

## Criadores

<table>
  <tr>
    <td align="center">
      <strong>Henrique Lourenco</strong><br/>
      <em>Criador</em><br/><br/>
      <a href="https://www.linkedin.com/in/henriquelourenco">LinkedIn</a> · 
      <a href="https://www.instagram.com/henrique.web3">Instagram</a>
    </td>
    <td align="center">
      <strong>Alexandre Senra</strong><br/>
      <em>Criador</em><br/><br/>
      <a href="https://www.linkedin.com/in/alexandresenra">LinkedIn</a> · 
      <a href="https://www.instagram.com/alexandresenra_">Instagram</a>
    </td>
  </tr>
</table>

---

## Apoie o Projeto

Este e um **projeto gratuito e open-source** construido com dedicacao e incontaveis horas de trabalho. Se o Crypto Hunter Pro te ajudou a recuperar seus fundos ou voce simplesmente quer apoiar o desenvolvimento open-source independente, por favor considere fazer uma doacao.

**Toda doacao, por menor que seja, nos ajuda a manter este projeto vivo e ativamente mantido.**

### Henrique Lourenco

| Rede | Endereco |
|------|----------|
| **BTC** | `bc1qpq0cgvyxczetumdf87345zzk0zr0xz96ampmhs` |
| **ETH** | `henriquelourenco.eth` |
| **PIX** | `henriquesamuel@yahoo.com.br` |

---

## Contribuindo

Contribuicoes sao bem-vindas! Seja correcao de bugs, novas funcionalidades, melhorias na documentacao ou traducoes, sinta-se a vontade para abrir um Pull Request.

1. Faca um Fork do repositorio
2. Crie sua branch de feature (`git checkout -b feature/funcionalidade-incrivel`)
3. Commit suas alteracoes (`git commit -m 'Adiciona funcionalidade incrivel'`)
4. Push para a branch (`git push origin feature/funcionalidade-incrivel`)
5. Abra um Pull Request

---

## Aviso Legal

Este software e fornecido **apenas para propositos legitimos de recuperacao**. Ele foi projetado para ajudar usuarios a recuperar acesso as suas proprias carteiras de criptomoedas. Os criadores nao sao responsaveis por qualquer uso indevido desta ferramenta. Sempre garanta que voce tem autorizacao legal para acessar qualquer carteira que tente recuperar.

---

<p align="center">
  <strong>Construido com dedicacao por Henrique Lourenco & Alexandre Senra</strong><br/>
  <em>Ajude-nos a manter este projeto vivo! Doe qualquer valor.</em><br/>
  <em>Software livre, feito com dedicacao. Apoie os criadores!</em>
</p>

---

---

<a name="english"></a>

# :us: ENGLISH

---

## What is Crypto Hunter Pro?

**Crypto Hunter Pro** is an open-source suite of tools designed to help cryptocurrency holders recover lost or forgotten seed phrases and scan multiple blockchains for funds. Built in Go for maximum performance, it combines two powerful modules that work together as a complete recovery pipeline.

> Lost your seed phrase order? Forgot which wallet you used? Have partial words? Crypto Hunter Pro was built for you.

---

## Modules

### Module 1: Unmixer Seed Search

The **Unmixer Seed Search** is a BIP39 seed phrase permutation engine that generates all valid seed combinations from partial or scrambled information.

**Key Features:**

| Feature | Description |
|---------|-------------|
| **4 Input Modes** | Simple (known order), Advanced Partial (wildcards + unknown order), Advanced Complete (full words + unknown order), Descrambler (test ALL permutations) |
| **Wildcard Support** | Use `*` for unknown parts: `aban*`, `bo*`, `*tion`, `*` (full unknown) |
| **9 BIP39 Languages** | English, Spanish, French, Italian, Portuguese, Japanese, Korean, Chinese Simplified, Chinese Traditional |
| **Smart Typo Correction** | Automatic suggestions when you mistype a word (Levenshtein distance + keyboard adjacency maps + transposition detection) |
| **BIP39 Checksum Validation** | Only outputs seeds that pass cryptographic checksum verification |
| **Excel Export** | Results saved in organized `.xlsx` files ready for CIE import |
| **Bilingual Interface** | Full English and Portuguese support |

**How It Works:**

```
Your scrambled/partial words → Unmixer → All valid BIP39 seed combinations → Excel file
```

---

### Module 2: CIE - Crypto Intelligence Engine

The **CIE (Crypto Intelligence Engine)** is a multi-chain wallet scanner that takes seed phrases and automatically checks balances across 31+ blockchains and 160+ tokens.

**Key Features:**

| Feature | Description |
|---------|-------------|
| **31 Blockchains** | Bitcoin, Ethereum, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Solana, Tron, Litecoin, Dogecoin, TON, Zcash, XRP, Stellar, Algorand, Sui, Near, and more |
| **160+ Tokens** | Automatic balance check for ERC-20, BEP-20, SPL, TRC-20 tokens |
| **Multi-Derivation Paths** | Test ALL address formats automatically (Legacy, SegWit, Native SegWit, Taproot) |
| **Multiple API Pools** | Failover system with multiple providers per network for maximum reliability |
| **Excel Import** | Direct import from Unmixer Seed Search results |
| **Real-time Display** | Watch addresses being checked live with balance indicators |
| **Detailed Reports** | Excel output with Summary, With Balance, and With History tabs |
| **Passphrase Support** | Optional BIP39 passphrase (25th word) testing |

**Supported Networks:**

```
EVM: Ethereum, BSC, Polygon, Arbitrum, Avalanche, Optimism, Base, Linea,
     Scroll, Gnosis, zkSync, Blast, Cronos, Celo, Berachain, Sonic, Mantle, Flare

UTXO: Bitcoin (4 address types), Bitcoin Cash, Litecoin (3 types), Dogecoin, Zcash

Other: Solana, Tron, TON, XRP, Stellar, Algorand, Sui, Near Protocol
```

---

## The Pipeline: Unmixer → CIE

The two modules work together as a complete recovery pipeline:

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  1. UNMIXER SEED SEARCH                                         │
│     Input: Your scrambled/partial seed words                    │
│     Output: All valid BIP39 seed combinations (Excel)           │
│                                                                 │
│                          ↓                                      │
│                                                                 │
│  2. CIE - CRYPTO INTELLIGENCE ENGINE                            │
│     Input: Valid seeds from Unmixer (or manual entry)            │
│     Output: Which seeds have funds, on which networks,          │
│             with full balance details (Excel report)             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Quick Start

### Requirements

- **Go 1.21+** (for compilation)
- **Windows, Linux, or macOS**
- **Internet connection** (CIE only, for blockchain queries)

### Compilation

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/crypto-hunter-pro.git
cd crypto-hunter-pro

# Compile Unmixer Seed Search
cd CRYPTO_HUNTER_PRO_Unmixer_Seed_Search
go build -o unmixer_seed_search .

# Compile CIE
cd ../CRYPTO_HUNTER_PRO_CIE
go build -o crypto_hunter_pro_cie .
```

### Windows (using the included batch file)

```batch
# Just double-click compilador.bat in each module folder
```

### Linux and macOS (from source)

> **Note:** Pre-compiled binaries (.exe) are for Windows. For Linux and macOS, compile directly from source code:

```bash
# Linux
cd CRYPTO_HUNTER_PRO_Unmixer_Seed_Search
GOOS=linux GOARCH=amd64 go build -o unmixer_seed_search .

cd ../CRYPTO_HUNTER_PRO_CIE
GOOS=linux GOARCH=amd64 go build -o crypto_hunter_pro_cie .

# macOS
cd CRYPTO_HUNTER_PRO_Unmixer_Seed_Search
GOOS=darwin GOARCH=amd64 go build -o unmixer_seed_search .

cd ../CRYPTO_HUNTER_PRO_CIE
GOOS=darwin GOARCH=amd64 go build -o crypto_hunter_pro_cie .
```

---

## Usage Examples

### Example 1: Descrambler Mode (Mode 4)

You have all 12 words but don't know the order:

```
Select Mode: 4

Enter all words separated by space:
> abandon ability able about above absent absorb abstract absurd abuse access acid

Processing all 479,001,600 permutations...
Valid seeds found: 1
```

### Example 2: Wildcard Recovery (Mode 1)

You know the order but forgot some words:

```
Select Mode: 1

Enter seed (use * for unknown):
> abandon * able about * absent absorb abstract absurd abuse access acid

Expanding wildcards against BIP39 wordlist...
Valid seeds found: 3
```

### Example 3: CIE Multi-Chain Scan

```
Enter seed phrase: abandon ability able about above absent absorb abstract absurd abuse access acid

Select networks: ALL
Test ALL derivation paths? YES

Scanning 31 networks...
[BTC] bc1q... - Balance: 0.00142 BTC ✓
[ETH] 0x...  - Balance: 0.5 ETH ✓
[ETH] 0x...  - USDT: 150.00 ✓
```

---

## Derivation Paths Supported

| Network | Format | Path |
|---------|--------|------|
| Bitcoin | Legacy (1...) | m/44'/0'/0'/0/x |
| Bitcoin | SegWit (3...) | m/49'/0'/0'/0/x |
| Bitcoin | Native SegWit (bc1q...) | m/84'/0'/0'/0/x |
| Bitcoin | Taproot (bc1p...) | m/86'/0'/0'/0/x |
| EVM (all) | Standard (0x...) | m/44'/60'/0'/0/x |
| Solana | Ed25519 | m/44'/501'/0'/0' |
| Tron | Standard (T...) | m/44'/195'/0'/0/x |
| Litecoin | Legacy/SegWit/Native | m/44'/2', m/49'/2', m/84'/2' |
| Dogecoin | Standard (D...) | m/44'/3'/0'/0/x |
| TON | Ed25519 | m/44'/607'/0' |
| XRP | Standard (r...) | m/44'/144'/0'/0/x |
| Stellar | Ed25519 (G...) | m/44'/148'/0' |

---

## Security

- **100% Offline** (Unmixer): The seed permutation engine works completely offline. No internet required, no data sent anywhere.
- **Open Source**: Every line of code is auditable. No hidden backdoors, no telemetry, no data collection.
- **Local Processing**: All cryptographic operations happen on YOUR machine. Seeds never leave your computer.
- **No External Dependencies at Runtime**: Compiled binary is self-contained.

> **WARNING**: Never share your seed phrases with anyone. Never enter your seed phrases on websites. Always verify you are running the official open-source version.

---

## Creators

<table>
  <tr>
    <td align="center">
      <strong>Henrique Lourenco</strong><br/>
      <em>Creator</em><br/><br/>
      <a href="https://www.linkedin.com/in/henriquelourenco">LinkedIn</a> · 
      <a href="https://www.instagram.com/henrique.web3">Instagram</a>
    </td>
    <td align="center">
      <strong>Alexandre Senra</strong><br/>
      <em>Creator</em><br/><br/>
      <a href="https://www.linkedin.com/in/alexandresenra">LinkedIn</a> · 
      <a href="https://www.instagram.com/alexandresenra_">Instagram</a>
    </td>
  </tr>
</table>

---

## Support the Project

This is a **free, open-source project** built with dedication and countless hours of work. If Crypto Hunter Pro helped you recover your funds or you simply want to support independent open-source development, please consider making a donation.

**Every donation, no matter how small, helps us keep this project alive and actively maintained.**

### Henrique Lourenco

| Network | Address |
|---------|---------|
| **BTC** | `bc1qpq0cgvyxczetumdf87345zzk0zr0xz96ampmhs` |
| **ETH** | `henriquelourenco.eth` |
| **PIX** | `henriquesamuel@yahoo.com.br` |

---

## Contributing

We welcome contributions! Whether it's bug fixes, new features, documentation improvements, or translations, feel free to open a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## Disclaimer

This software is provided for **legitimate recovery purposes only**. It is designed to help users recover access to their own cryptocurrency wallets. The creators are not responsible for any misuse of this tool. Always ensure you have legal authorization to access any wallet you attempt to recover.

---

## Star History

If this project helped you, please give it a star! It helps others find this tool.

---

<p align="center">
  <strong>Built with dedication by Henrique Lourenco & Alexandre Senra</strong><br/>
  <em>Help us keep this project alive! Donate any amount.</em><br/>
  <em>Free software, made with dedication. Support the creators!</em>
</p>
