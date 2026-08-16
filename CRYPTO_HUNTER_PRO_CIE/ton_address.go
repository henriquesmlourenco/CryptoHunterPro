package main

import (
"crypto/sha256"
"encoding/base64"
"fmt"
)

// ============================================================================
// TON WALLET V4R2 - DERIVACAO DE ENDERECO
// ============================================================================

const walletV4R2CodeBOC = "te6ccgECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8QERITAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNBgcCASAICQB4AfoA9AQw+CdvIjBQCqEhvvLgUIIQcGx1Z4MesXCAGFAEywUmzxZY+gIZ9ADLaRfLH1Jgyz8gyYBA+wAGAIpQBIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UAXKwjiOCEGRzdHKDHrFwgBhQBcsFUAPPFiP6AhPLassfyz/JgED7AJJfA+ICASAKCwBZvSQrb2omhAgKBrkPoCGEcNQICEekk30pkQzmkD6f+YN4EoAbeBAUiYcVnzGEAgFYDA0AEbjJftRNDXCx+AA9sp37UTQgQFA1yH0BDACyMoHy//J0AGBAQj0Cm+hMYAIBIA4PABmtznaiaEAga5Drhf/AABmvHfaiaEAQa5DrhY/AAG7SB/oA1NQi+QAFyMoHFcv/ydB3dIAYyMsFywIizxZQBfoCFMtrEszMyXP7AMhAFIEBCPRR8qcCAHCBAQjXGPoA0z/IVCBHgQEI9FHyp4IQbm90ZXB0gBjIywXLAlAGzxZQBPoCFMtqEssfyz/Jc/sAAgBsgQEI1xj6ANM/MFIkgQEI9Fnyp4IQZHN0cnB0gBjIywXLAlAFzxZQA/oCE8tqyx8Syz/Jc/sAAAr0AMntVA=="

const walletV4R2DefaultWalletID uint32 = 698983191

// tonCell representa uma cell TON simplificada
type tonCell struct {
data []byte
bits int
refs []*tonCell
}

func newTonCell() *tonCell {
return &tonCell{data: make([]byte, 0, 128), bits: 0, refs: make([]*tonCell, 0, 4)}
}

func (c *tonCell) writeBit(bit int) {
byteIdx := c.bits / 8
bitIdx := 7 - (c.bits % 8)
if byteIdx >= len(c.data) {
c.data = append(c.data, 0)
}
if bit == 1 {
c.data[byteIdx] |= (1 << bitIdx)
}
c.bits++
}

func (c *tonCell) writeUint(value uint64, bits int) {
for i := bits - 1; i >= 0; i-- {
if (value>>uint(i))&1 == 1 {
c.writeBit(1)
} else {
c.writeBit(0)
}
}
}

func (c *tonCell) writeBytes(data []byte) {
for _, b := range data {
for i := 7; i >= 0; i-- {
if (b>>uint(i))&1 == 1 {
c.writeBit(1)
} else {
c.writeBit(0)
}
}
}
}

func (c *tonCell) representationHash() [32]byte {
repr := c.representation()
return sha256.Sum256(repr)
}

func (c *tonCell) representation() []byte {
d1 := byte(len(c.refs))
d2 := byte(c.bits/8) * 2
if c.bits%8 != 0 {
d2 = byte(c.bits/8)*2 + 1
}

var repr []byte
repr = append(repr, d1, d2)

if c.bits%8 == 0 {
repr = append(repr, c.data[:c.bits/8]...)
} else {
dataLen := (c.bits + 7) / 8
paddedData := make([]byte, dataLen)
copy(paddedData, c.data[:dataLen])
lastByteIdx := c.bits / 8
lastBitIdx := 7 - (c.bits % 8)
paddedData[lastByteIdx] |= (1 << lastBitIdx)
repr = append(repr, paddedData...)
}

for _, ref := range c.refs {
depth := ref.depth()
repr = append(repr, byte(depth>>8), byte(depth&0xFF))
}

for _, ref := range c.refs {
hash := ref.representationHash()
repr = append(repr, hash[:]...)
}

return repr
}

func (c *tonCell) depth() int {
if len(c.refs) == 0 {
return 0
}
maxDepth := 0
for _, ref := range c.refs {
d := ref.depth()
if d > maxDepth {
maxDepth = d
}
}
return maxDepth + 1
}

// parseBOC parseia um Bag of Cells e retorna a root cell
func parseBOC(bocBase64 string) (*tonCell, error) {
bocBytes, err := base64.StdEncoding.DecodeString(bocBase64)
if err != nil {
return nil, err
}

if len(bocBytes) < 6 {
return nil, fmt.Errorf("BOC too short")
}

offset := 4 // skip magic b5ee9c72
flagsByte := bocBytes[offset]
offset++

hasIdx := (flagsByte >> 7) & 1
refSize := int(flagsByte & 0x07)

offsetSize := int(bocBytes[offset])
offset++

cellsCount := readNBytesInt(bocBytes[offset:offset+refSize], refSize)
offset += refSize

_ = readNBytesInt(bocBytes[offset:offset+refSize], refSize) // roots_count
offset += refSize

_ = readNBytesInt(bocBytes[offset:offset+refSize], refSize) // absent_count
offset += refSize

_ = readNBytesInt(bocBytes[offset:offset+offsetSize], offsetSize) // total_cells_size
offset += offsetSize

_ = readNBytesInt(bocBytes[offset:offset+refSize], refSize) // root index
offset += refSize

if hasIdx == 1 {
offset += cellsCount * offsetSize
}

cells := make([]*tonCell, cellsCount)
cellRefs := make([][]int, cellsCount)

for i := 0; i < cellsCount; i++ {
cells[i] = newTonCell()

d1 := bocBytes[offset]
offset++
refsCount := int(d1 & 0x07)

d2 := bocBytes[offset]
offset++
dataByteSize := int(d2) / 2
if d2%2 != 0 {
dataByteSize++
}

totalBits := int(d2) / 2 * 8
cellData := bocBytes[offset : offset+dataByteSize]
offset += dataByteSize

if d2%2 != 0 && dataByteSize > 0 {
lastByte := cellData[dataByteSize-1]
paddingBits := 0
for bit := 0; bit < 8; bit++ {
if (lastByte>>uint(bit))&1 == 1 {
paddingBits = bit + 1
break
}
}
totalBits = dataByteSize*8 - paddingBits
}

cells[i].data = make([]byte, len(cellData))
copy(cells[i].data, cellData)
cells[i].bits = totalBits

refs := make([]int, refsCount)
for r := 0; r < refsCount; r++ {
refs[r] = readNBytesInt(bocBytes[offset:offset+refSize], refSize)
offset += refSize
}
cellRefs[i] = refs
}

for i := 0; i < cellsCount; i++ {
for _, refIdx := range cellRefs[i] {
if refIdx < cellsCount {
cells[i].refs = append(cells[i].refs, cells[refIdx])
}
}
}

return cells[0], nil
}

func readNBytesInt(data []byte, n int) int {
result := 0
for i := 0; i < n; i++ {
result = (result << 8) | int(data[i])
}
return result
}

// tonWalletV4R2Address gera o endereço TON wallet V4R2 user-friendly (UQ...)
func tonWalletV4R2Address(pubKey []byte) (string, error) {
codeCell, err := parseBOC(walletV4R2CodeBOC)
if err != nil {
return "", fmt.Errorf("erro ao parsear BOC do wallet V4R2: %v", err)
}

dataCell := newTonCell()
dataCell.writeUint(0, 32)                                 // seqno = 0
dataCell.writeUint(uint64(walletV4R2DefaultWalletID), 32) // wallet_id
dataCell.writeBytes(pubKey)                                // 256-bit public key
dataCell.writeBit(0)                                       // empty plugins dict

stateInit := newTonCell()
stateInit.writeBit(0) // split_depth: absent
stateInit.writeBit(0) // tick_tock: absent
stateInit.writeBit(1) // code: present
stateInit.writeBit(1) // data: present
stateInit.writeBit(0) // library: absent
stateInit.refs = append(stateInit.refs, codeCell, dataCell)

hash := stateInit.representationHash()

// Non-bounceable user-friendly address (UQ...)
var addrBytes [36]byte
addrBytes[0] = 0x51 // non-bounceable
addrBytes[1] = 0x00 // workchain 0
copy(addrBytes[2:34], hash[:])

crc := crc16XMODEM(addrBytes[:34])
addrBytes[34] = byte(crc >> 8)
addrBytes[35] = byte(crc & 0xFF)

address := base64.URLEncoding.EncodeToString(addrBytes[:])
return address, nil
}

// tonWalletV4R2RawAddress gera o raw address (0:hex)
func tonWalletV4R2RawAddress(pubKey []byte) (string, error) {
codeCell, err := parseBOC(walletV4R2CodeBOC)
if err != nil {
return "", fmt.Errorf("erro ao parsear BOC do wallet V4R2: %v", err)
}

dataCell := newTonCell()
dataCell.writeUint(0, 32)
dataCell.writeUint(uint64(walletV4R2DefaultWalletID), 32)
dataCell.writeBytes(pubKey)
dataCell.writeBit(0)

stateInit := newTonCell()
stateInit.writeBit(0)
stateInit.writeBit(0)
stateInit.writeBit(1)
stateInit.writeBit(1)
stateInit.writeBit(0)
stateInit.refs = append(stateInit.refs, codeCell, dataCell)

hash := stateInit.representationHash()
rawAddr := fmt.Sprintf("0:%x", hash[:])
return rawAddr, nil
}

// crc16XMODEM calcula CRC16 com polinômio XMODEM (0x1021)
func crc16XMODEM(data []byte) uint16 {
crc := uint16(0x0000)
for _, b := range data {
crc ^= uint16(b) << 8
for i := 0; i < 8; i++ {
if crc&0x8000 != 0 {
crc = (crc << 1) ^ 0x1021
} else {
crc <<= 1
}
}
}
return crc
}
