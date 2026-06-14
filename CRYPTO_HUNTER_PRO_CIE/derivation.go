package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// ============================================================================
// BIP32 - Derivacao HD Wallet (secp256k1)
// ============================================================================

type ExtendedKey struct {
	Key       []byte
	ChainCode []byte
	IsPrivate bool
}

func SeedToMasterKey(seed []byte) *ExtendedKey {
	hmacKey := []byte("Bitcoin seed")
	mac := hmac.New(sha512.New, hmacKey)
	mac.Write(seed)
	result := mac.Sum(nil)
	return &ExtendedKey{
		Key:       result[:32],
		ChainCode: result[32:],
		IsPrivate: true,
	}
}

func (ek *ExtendedKey) DeriveChild(index uint32) (*ExtendedKey, error) {
	var data []byte
	if index >= 0x80000000 {
		data = append([]byte{0x00}, ek.Key...)
	} else {
		pubKey := privateKeyToCompressedPubKey(ek.Key)
		data = pubKey
	}
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	data = append(data, indexBytes...)

	mac := hmac.New(sha512.New, ek.ChainCode)
	mac.Write(data)
	result := mac.Sum(nil)

	curve := btcec.S256()
	keyInt := new(big.Int).SetBytes(result[:32])
	parentInt := new(big.Int).SetBytes(ek.Key)
	keyInt.Add(keyInt, parentInt)
	keyInt.Mod(keyInt, curve.N)

	childKey := keyInt.Bytes()
	if len(childKey) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(childKey):], childKey)
		childKey = padded
	}

	return &ExtendedKey{
		Key:       childKey,
		ChainCode: result[32:],
		IsPrivate: true,
	}, nil
}

func DerivePath(masterKey *ExtendedKey, purpose, coinType, account, change, index uint32) (*ExtendedKey, error) {
	key := masterKey
	var err error

	indices := []uint32{
		purpose + 0x80000000,
		coinType + 0x80000000,
		account + 0x80000000,
		change,
		index,
	}

	for _, idx := range indices {
		key, err = key.DeriveChild(idx)
		if err != nil {
			return nil, err
		}
	}

	return key, nil
}

// ============================================================================
// SLIP-0010 - Derivacao Ed25519 (Solana, TON)
// ============================================================================

type Ed25519Key struct {
	Key       []byte // 32 bytes private key
	ChainCode []byte // 32 bytes chain code
}

func SeedToEd25519MasterKey(seed []byte) *Ed25519Key {
	hmacKey := []byte("ed25519 seed")
	mac := hmac.New(sha512.New, hmacKey)
	mac.Write(seed)
	result := mac.Sum(nil)
	return &Ed25519Key{
		Key:       result[:32],
		ChainCode: result[32:],
	}
}

func (ek *Ed25519Key) DeriveChild(index uint32) *Ed25519Key {
	// Ed25519 SLIP-0010 only supports hardened derivation
	index = index | 0x80000000

	data := make([]byte, 1+32+4)
	data[0] = 0x00
	copy(data[1:33], ek.Key)
	binary.BigEndian.PutUint32(data[33:], index)

	mac := hmac.New(sha512.New, ek.ChainCode)
	mac.Write(data)
	result := mac.Sum(nil)

	return &Ed25519Key{
		Key:       result[:32],
		ChainCode: result[32:],
	}
}

// DeriveSolanaKeypair derives Solana keypair from seed phrase
// Path: m/44'/501'/0'/0' (all hardened for Ed25519)
func DeriveSolanaKeypair(seedPhrase string, passphrase string, index int) (string, string, error) {
	seed := bip39.NewSeed(seedPhrase, passphrase)
	masterKey := SeedToEd25519MasterKey(seed)

	// m/44'/501'/account'/0'
	key := masterKey.DeriveChild(44)
	key = key.DeriveChild(501)
	key = key.DeriveChild(uint32(index))
	key = key.DeriveChild(0)

	// Generate ed25519 keypair
	privKey := ed25519.NewKeyFromSeed(key.Key)
	pubKey := privKey.Public().(ed25519.PublicKey)

	// Solana address = Base58 of public key (32 bytes)
	address := base58Encode([]byte(pubKey))

	// Private key = Base58 of full 64-byte keypair
	fullKey := make([]byte, 64)
	copy(fullKey[:32], key.Key)
	copy(fullKey[32:], pubKey)
	privKeyStr := base58Encode(fullKey)

	return address, privKeyStr, nil
}

// DeriveTONKeypair derives TON keypair from seed phrase
// Path: m/44'/607'/0' (simplified)
func DeriveTONKeypair(seedPhrase string, passphrase string, index int) (string, string, error) {
	seed := bip39.NewSeed(seedPhrase, passphrase)
	masterKey := SeedToEd25519MasterKey(seed)

	// m/44'/607'/account'
	key := masterKey.DeriveChild(44)
	key = key.DeriveChild(607)
	key = key.DeriveChild(uint32(index))

	// Generate ed25519 keypair
	privKey := ed25519.NewKeyFromSeed(key.Key)
	pubKey := privKey.Public().(ed25519.PublicKey)

	// TON address = simplified raw address (hex of pubkey for now)
	// Real TON uses workchain + state init hash, but for scanning we use pubkey hex
	address := hex.EncodeToString(pubKey)
	privKeyStr := hex.EncodeToString(key.Key)

	return address, privKeyStr, nil
}

// ============================================================================
// CONVERSAO DE CHAVES
// ============================================================================

func privateKeyToCompressedPubKey(privKey []byte) []byte {
	_, pub := btcec.PrivKeyFromBytes(privKey)
	return pub.SerializeCompressed()
}

func privateKeyToUncompressedPubKey(privKey []byte) []byte {
	_, pub := btcec.PrivKeyFromBytes(privKey)
	return pub.SerializeUncompressed()
}

func privateKeyToECDSA(privKeyBytes []byte) *ecdsa.PrivateKey {
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
	return privKey.ToECDSA()
}

// ============================================================================
// GERACAO DE ENDERECOS - BTC
// ============================================================================

func deriveBTCLegacy(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	addr, _ := btcutil.NewAddressPubKeyHash(hash160, &chaincfg.MainNetParams)
	wif, _ := btcutil.NewWIF(func() *btcec.PrivateKey {
		pk, _ := btcec.PrivKeyFromBytes(privKey)
		return pk
	}(), &chaincfg.MainNetParams, true)
	return addr.EncodeAddress(), wif.String()
}

func deriveBTCSegWit(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	witnessProg := append([]byte{0x00, 0x14}, hash160...)
	scriptHash := hash160Bytes(witnessProg)
	addr, _ := btcutil.NewAddressScriptHashFromHash(scriptHash, &chaincfg.MainNetParams)
	wif, _ := btcutil.NewWIF(func() *btcec.PrivateKey {
		pk, _ := btcec.PrivKeyFromBytes(privKey)
		return pk
	}(), &chaincfg.MainNetParams, true)
	return addr.EncodeAddress(), wif.String()
}

func deriveBTCNativeSegWit(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	addr, _ := btcutil.NewAddressWitnessPubKeyHash(hash160, &chaincfg.MainNetParams)
	wif, _ := btcutil.NewWIF(func() *btcec.PrivateKey {
		pk, _ := btcec.PrivKeyFromBytes(privKey)
		return pk
	}(), &chaincfg.MainNetParams, true)
	return addr.EncodeAddress(), wif.String()
}

func deriveBTCTaproot(privKey []byte) (string, string) {
	privKeyObj, pubKeyObj := btcec.PrivKeyFromBytes(privKey)
	taprootKey := txscript.ComputeTaprootKeyNoScript(pubKeyObj)
	addr, err := btcutil.NewAddressTaproot(taprootKey.SerializeCompressed()[1:], &chaincfg.MainNetParams)
	if err != nil {
		return "", ""
	}
	wif, _ := btcutil.NewWIF(privKeyObj, &chaincfg.MainNetParams, true)
	return addr.EncodeAddress(), wif.String()
}

// ============================================================================
// GERACAO DE ENDERECOS - BCH
// ============================================================================

func deriveBCHCashAddr(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	prefix := "bitcoincash"
	addr := encodeCashAddr(prefix, 0, hash160)
	wif, _ := btcutil.NewWIF(func() *btcec.PrivateKey {
		pk, _ := btcec.PrivKeyFromBytes(privKey)
		return pk
	}(), &chaincfg.MainNetParams, true)
	return addr, wif.String()
}

func deriveBCHLegacy(privKey []byte) (string, string) {
	return deriveBTCLegacy(privKey)
}

func encodeCashAddr(prefix string, version byte, payload []byte) string {
	data := convertBits(append([]byte{version << 3}, payload...), 8, 5, true)
	prefixData := cashAddrPrefixExpand(prefix)
	values := append(prefixData, data...)
	values = append(values, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
	polymod := cashAddrPolymod(values) ^ 1
	checksum := make([]byte, 8)
	for i := 0; i < 8; i++ {
		checksum[i] = byte((polymod >> uint(5*(7-i))) & 0x1f)
	}
	data = append(data, checksum...)
	charset := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	result := prefix + ":"
	for _, d := range data {
		result += string(charset[d])
	}
	return result
}

func cashAddrPrefixExpand(prefix string) []byte {
	result := make([]byte, len(prefix)+1)
	for i, c := range prefix {
		result[i] = byte(c) & 0x1f
	}
	result[len(prefix)] = 0
	return result
}

func cashAddrPolymod(values []byte) uint64 {
	c := uint64(1)
	for _, d := range values {
		c0 := c >> 35
		c = ((c & 0x07ffffffff) << 5) ^ uint64(d)
		if c0&0x01 != 0 { c ^= 0x98f2bc8e61 }
		if c0&0x02 != 0 { c ^= 0x79b76d99e2 }
		if c0&0x04 != 0 { c ^= 0xf33e5fb3c4 }
		if c0&0x08 != 0 { c ^= 0xae2eabe2a8 }
		if c0&0x10 != 0 { c ^= 0x1e4f43e470 }
	}
	return c
}

func convertBits(data []byte, fromBits, toBits int, pad bool) []byte {
	acc := 0
	bits := 0
	var result []byte
	maxv := (1 << toBits) - 1
	for _, value := range data {
		acc = (acc << fromBits) | int(value)
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			result = append(result, byte((acc>>bits)&maxv))
		}
	}
	if pad && bits > 0 {
		result = append(result, byte((acc<<(toBits-bits))&maxv))
	}
	return result
}

// ============================================================================
// GERACAO DE ENDERECOS - ETH/EVM
// ============================================================================

func deriveEVMAddress(privKey []byte) (string, string) {
	pubKey := privateKeyToUncompressedPubKey(privKey)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey[1:])
	addrBytes := hash.Sum(nil)[12:]
	addr := eip55Checksum(addrBytes)
	return addr, hex.EncodeToString(privKey)
}

func eip55Checksum(addr []byte) string {
	hexAddr := hex.EncodeToString(addr)
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(hexAddr))
	hashBytes := hash.Sum(nil)
	result := "0x"
	for i, c := range hexAddr {
		if c >= '0' && c <= '9' {
			result += string(c)
		} else {
			hashByte := hashBytes[i/2]
			var nibble byte
			if i%2 == 0 {
				nibble = hashByte >> 4
			} else {
				nibble = hashByte & 0x0f
			}
			if nibble >= 8 {
				result += string(c - 32)
			} else {
				result += string(c)
			}
		}
	}
	return result
}

// ============================================================================
// GERACAO DE ENDERECOS - TRON (TRX)
// ============================================================================

func deriveTronAddress(privKey []byte) (string, string) {
	pubKey := privateKeyToUncompressedPubKey(privKey)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey[1:])
	addrBytes := hash.Sum(nil)[12:]
	tronAddr := append([]byte{0x41}, addrBytes...)
	addr := base58CheckEncode(tronAddr)
	return addr, hex.EncodeToString(privKey)
}

func base58CheckEncode(payload []byte) string {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	checksum := second[:4]
	full := append(payload, checksum...)
	return base58Encode(full)
}

func base58Encode(input []byte) string {
	alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	x := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)
	var result []byte
	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append([]byte{alphabet[mod.Int64()]}, result...)
	}
	for _, b := range input {
		if b != 0 {
			break
		}
		result = append([]byte{'1'}, result...)
	}
	return string(result)
}

// ============================================================================
// GERACAO DE ENDERECOS - LTC
// ============================================================================

func deriveLTCLegacy(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	payload := append([]byte{0x30}, hash160...)
	addr := base58CheckEncode(payload)
	wif := encodeLTCWIF(privKey)
	return addr, wif
}

func deriveLTCSegWit(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	witnessProg := append([]byte{0x00, 0x14}, hash160...)
	scriptHash := hash160Bytes(witnessProg)
	payload := append([]byte{0x32}, scriptHash...)
	addr := base58CheckEncode(payload)
	wif := encodeLTCWIF(privKey)
	return addr, wif
}

func deriveLTCNativeSegWit(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	conv := convertBits(hash160, 8, 5, true)
	addr, _ := bech32Encode("ltc", 0, conv)
	wif := encodeLTCWIF(privKey)
	return addr, wif
}

func encodeLTCWIF(privKey []byte) string {
	payload := append([]byte{0xB0}, privKey...)
	payload = append(payload, 0x01)
	return base58CheckEncode(payload)
}

// ============================================================================
// GERACAO DE ENDERECOS - DOGE
// ============================================================================

func deriveDOGEAddress(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)
	payload := append([]byte{0x1E}, hash160...)
	addr := base58CheckEncode(payload)
	wifPayload := append([]byte{0x9E}, privKey...)
	wifPayload = append(wifPayload, 0x01)
	wif := base58CheckEncode(wifPayload)
	return addr, wif
}

// ============================================================================
// GERACAO DE ENDERECOS - ZCASH (Transparent t1...)
// ============================================================================

func deriveZcashTransparent(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)

	// Zcash transparent uses 2-byte version prefix: 0x1CB8 for t1...
	payload := make([]byte, 0, 22)
	payload = append(payload, 0x1C, 0xB8)
	payload = append(payload, hash160...)
	addr := base58CheckEncode(payload)

	// Zcash WIF uses version byte 0x80 (same as BTC)
	wifPayload := append([]byte{0x80}, privKey...)
	wifPayload = append(wifPayload, 0x01)
	wif := base58CheckEncode(wifPayload)

	return addr, wif
}

// ============================================================================
// GERACAO DE ENDERECOS - XRP (Ripple)
// ============================================================================

func deriveXRPAddress(privKey []byte) (string, string) {
	pubKey := privateKeyToCompressedPubKey(privKey)
	hash160 := hash160Bytes(pubKey)

	// XRP uses its own Base58 alphabet and version byte 0x00
	payload := append([]byte{0x00}, hash160...)

	// Double SHA-256 checksum
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	checksum := second[:4]
	full := append(payload, checksum...)

	// XRP Base58 alphabet (different from Bitcoin!)
	xrpAlphabet := "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
	addr := base58EncodeCustom(full, xrpAlphabet)

	wif, _ := btcutil.NewWIF(func() *btcec.PrivateKey {
		pk, _ := btcec.PrivKeyFromBytes(privKey)
		return pk
	}(), &chaincfg.MainNetParams, true)

	return addr, wif.String()
}

func base58EncodeCustom(input []byte, alphabet string) string {
	x := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)
	var result []byte
	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append([]byte{alphabet[mod.Int64()]}, result...)
	}
	for _, b := range input {
		if b != 0 {
			break
		}
		result = append([]byte{alphabet[0]}, result...)
	}
	return string(result)
}

// ============================================================================
// GERACAO DE ENDERECOS - STELLAR (XLM)
// ============================================================================

func DeriveStellarKeypair(seedPhrase string, passphrase string, index int) (string, string, error) {
	seed := bip39.NewSeed(seedPhrase, passphrase)
	masterKey := SeedToEd25519MasterKey(seed)

	// m/44'/148'/index' (SLIP-10 Ed25519, all hardened)
	key := masterKey.DeriveChild(44)
	key = key.DeriveChild(148)
	key = key.DeriveChild(uint32(index))

	// Generate ed25519 keypair
	privKey := ed25519.NewKeyFromSeed(key.Key)
	pubKey := privKey.Public().(ed25519.PublicKey)

	// Stellar address = StrKey encoding with version byte 6<<3 = 48 (G...)
	address := stellarEncodeAddress(pubKey)

	// Stellar secret = StrKey encoding with version byte 18<<3 = 144 (S...)
	secretKey := stellarEncodeSecret(key.Key)

	return address, secretKey, nil
}

func stellarEncodeAddress(pubKey []byte) string {
	// Version byte for public key: 6 << 3 = 48
	payload := append([]byte{6 << 3}, pubKey...)
	checksum := stellarCRC16(payload)
	full := append(payload, checksum...)
	return base32Encode(full)
}

func stellarEncodeSecret(privKey []byte) string {
	// Version byte for secret key: 18 << 3 = 144
	payload := append([]byte{18 << 3}, privKey...)
	checksum := stellarCRC16(payload)
	full := append(payload, checksum...)
	return base32Encode(full)
}

func stellarCRC16(data []byte) []byte {
	crc := uint16(0x0000)
	poly := uint16(0x1021)
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc = crc << 1
			}
		}
	}
	return []byte{byte(crc & 0xFF), byte((crc >> 8) & 0xFF)}
}

func base32Encode(data []byte) string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	var result strings.Builder
	buffer := 0
	bitsLeft := 0
	for _, b := range data {
		buffer = (buffer << 8) | int(b)
		bitsLeft += 8
		for bitsLeft >= 5 {
			bitsLeft -= 5
			result.WriteByte(alphabet[(buffer>>bitsLeft)&0x1F])
		}
	}
	if bitsLeft > 0 {
		result.WriteByte(alphabet[(buffer<<(5-bitsLeft))&0x1F])
	}
	return result.String()
}

// ============================================================================
// GERACAO DE ENDERECOS - ALGORAND (ALGO)
// ============================================================================

func DeriveAlgorandKeypair(seedPhrase string, passphrase string, index int) (string, string, error) {
	seed := bip39.NewSeed(seedPhrase, passphrase)
	masterKey := SeedToEd25519MasterKey(seed)

	// m/44'/283'/index'/0/0 (SLIP-10 Ed25519)
	key := masterKey.DeriveChild(44)
	key = key.DeriveChild(283)
	key = key.DeriveChild(uint32(index))
	key = key.DeriveChild(0)
	key = key.DeriveChild(0)

	// Generate ed25519 keypair
	privKey := ed25519.NewKeyFromSeed(key.Key)
	pubKey := privKey.Public().(ed25519.PublicKey)

	// Algorand address = Base32(pubkey + last 4 bytes of SHA-512/256(pubkey))
	checksum := sha512_256(pubKey)
	addrBytes := append([]byte(pubKey), checksum[28:]...)
	address := base32Encode(addrBytes)

	// Private key = Base64 of 32-byte seed
	privKeyStr := hex.EncodeToString(key.Key)

	return address, privKeyStr, nil
}

func sha512_256(data []byte) []byte {
	h := sha512.New512_256()
	h.Write(data)
	return h.Sum(nil)
}

// ============================================================================
// GERACAO DE ENDERECOS - SUI
// ============================================================================

func DeriveSuiKeypair(seedPhrase string, passphrase string, index int) (string, string, error) {
	seed := bip39.NewSeed(seedPhrase, passphrase)
	masterKey := SeedToEd25519MasterKey(seed)

	// m/44'/784'/0'/0'/index' (SLIP-10 Ed25519, all hardened)
	key := masterKey.DeriveChild(44)
	key = key.DeriveChild(784)
	key = key.DeriveChild(0)
	key = key.DeriveChild(0)
	key = key.DeriveChild(uint32(index))

	// Generate ed25519 keypair
	privKey := ed25519.NewKeyFromSeed(key.Key)
	pubKey := privKey.Public().(ed25519.PublicKey)

	// Sui address = 0x + hex(Blake2b-256(0x00 || pubkey))
	// 0x00 = Ed25519 flag
	data := append([]byte{0x00}, pubKey...)
	hash, _ := blake2b.New256(nil)
	hash.Write(data)
	digest := hash.Sum(nil)
	address := "0x" + hex.EncodeToString(digest)

	privKeyStr := hex.EncodeToString(key.Key)

	return address, privKeyStr, nil
}

// ============================================================================
// GERACAO DE ENDERECOS - NEAR
// ============================================================================

func DeriveNearKeypair(seedPhrase string, passphrase string, index int) (string, string, error) {
	seed := bip39.NewSeed(seedPhrase, passphrase)
	masterKey := SeedToEd25519MasterKey(seed)

	// m/44'/397'/index' (SLIP-10 Ed25519)
	key := masterKey.DeriveChild(44)
	key = key.DeriveChild(397)
	key = key.DeriveChild(uint32(index))

	// Generate ed25519 keypair
	privKey := ed25519.NewKeyFromSeed(key.Key)
	pubKey := privKey.Public().(ed25519.PublicKey)

	// NEAR implicit account = hex of public key (64 chars)
	address := hex.EncodeToString(pubKey)

	// Private key in ed25519:base58 format
	fullKey := make([]byte, 64)
	copy(fullKey[:32], key.Key)
	copy(fullKey[32:], pubKey)
	privKeyStr := "ed25519:" + base58Encode(fullKey)

	return address, privKeyStr, nil
}

// ============================================================================
// FUNCOES AUXILIARES
// ============================================================================

func hash160Bytes(data []byte) []byte {
	sha := sha256.Sum256(data)
	ripeHasher := ripemd160.New()
	ripeHasher.Write(sha[:])
	return ripeHasher.Sum(nil)
}

func bech32Encode(hrp string, witnessVersion byte, data []byte) (string, error) {
	values := append([]byte{witnessVersion}, data...)
	checksum := bech32CreateChecksum(hrp, values, bech32Const)
	combined := append(values, checksum...)
	charset := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	result := hrp + "1"
	for _, v := range combined {
		result += string(charset[v])
	}
	return result, nil
}

const bech32Const = 1
const bech32mConst = 0x2bc830a3

func bech32Polymod(values []byte) uint32 {
	gen := []uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}
	chk := uint32(1)
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ uint32(v)
		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

func bech32HRPExpand(hrp string) []byte {
	result := make([]byte, 0, len(hrp)*2+1)
	for _, c := range hrp {
		result = append(result, byte(c>>5))
	}
	result = append(result, 0)
	for _, c := range hrp {
		result = append(result, byte(c&31))
	}
	return result
}

func bech32CreateChecksum(hrp string, data []byte, spec uint32) []byte {
	values := append(bech32HRPExpand(hrp), data...)
	values = append(values, []byte{0, 0, 0, 0, 0, 0}...)
	polymod := bech32Polymod(values) ^ spec
	checksum := make([]byte, 6)
	for i := 0; i < 6; i++ {
		checksum[i] = byte((polymod >> uint(5*(5-i))) & 31)
	}
	return checksum
}

// ============================================================================
// FUNCAO PRINCIPAL DE DERIVACAO
// ============================================================================

func DeriveAddress(seedPhrase string, passphrase string, dp DerivationPath, index int) (string, string, error) {
	// Ed25519 chains use SLIP-0010, not BIP32
	if dp.AddressType == "solana" {
		return DeriveSolanaKeypair(seedPhrase, passphrase, index)
	}
	if dp.AddressType == "ton" {
		return DeriveTONKeypair(seedPhrase, passphrase, index)
	}
	if dp.AddressType == "stellar" {
		return DeriveStellarKeypair(seedPhrase, passphrase, index)
	}
	if dp.AddressType == "algorand" {
		return DeriveAlgorandKeypair(seedPhrase, passphrase, index)
	}
	if dp.AddressType == "sui" {
		return DeriveSuiKeypair(seedPhrase, passphrase, index)
	}
	if dp.AddressType == "near" {
		return DeriveNearKeypair(seedPhrase, passphrase, index)
	}

	// Standard BIP32 secp256k1 derivation
	seed := bip39.NewSeed(seedPhrase, passphrase)
	masterKey := SeedToMasterKey(seed)

	derivedKey, err := DerivePath(masterKey, dp.Purpose, dp.CoinType, 0, 0, uint32(index))
	if err != nil {
		return "", "", fmt.Errorf("erro na derivacao: %v", err)
	}

	privKey := derivedKey.Key
	var address, privateKeyStr string

	switch dp.AddressType {
	case "legacy":
		if dp.CoinType == 0 {
			address, privateKeyStr = deriveBTCLegacy(privKey)
		} else if dp.CoinType == 145 {
			address, privateKeyStr = deriveBCHLegacy(privKey)
		}
	case "segwit":
		address, privateKeyStr = deriveBTCSegWit(privKey)
	case "native":
		address, privateKeyStr = deriveBTCNativeSegWit(privKey)
	case "taproot":
		address, privateKeyStr = deriveBTCTaproot(privKey)
	case "cashaddr":
		address, privateKeyStr = deriveBCHCashAddr(privKey)
	case "evm":
		address, privateKeyStr = deriveEVMAddress(privKey)
	case "tron":
		address, privateKeyStr = deriveTronAddress(privKey)
	case "ltc_legacy":
		address, privateKeyStr = deriveLTCLegacy(privKey)
	case "ltc_segwit":
		address, privateKeyStr = deriveLTCSegWit(privKey)
	case "ltc_native":
		address, privateKeyStr = deriveLTCNativeSegWit(privKey)
	case "doge":
		address, privateKeyStr = deriveDOGEAddress(privKey)
	case "zcash":
		address, privateKeyStr = deriveZcashTransparent(privKey)
	case "xrp":
		address, privateKeyStr = deriveXRPAddress(privKey)
	default:
		return "", "", fmt.Errorf("tipo de endereco nao suportado: %s", dp.AddressType)
	}

	return address, privateKeyStr, nil
}

func ValidateSeedPhrase(seedPhrase string) bool {
	return bip39.IsMnemonicValid(seedPhrase)
}
