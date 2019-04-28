package neo

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"math"
	"math/big"
)

const NEO_ADDRESS_VERSION  = 0x17

func Hash256(str string) []byte {
	s := sha256.New()
	s.Write([]byte(str))
	res := s.Sum(nil)
	s.Reset()
	s.Write(res)
	return s.Sum(nil)
}

func HexEncode(str string) string {
	return hex.EncodeToString([]byte(str))
}

func HexDecode(str string) string {
	b, _ := hex.DecodeString(str)
	return string(b)
}

func HexDecodeBytes(bys []byte) []byte {
	var dst = make([]byte, hex.EncodedLen(len(bys)))
	hex.Decode(dst, bys)
	return dst
}

func HexEncodeBytes(bys []byte) []byte {
	var dst = make([]byte, hex.EncodedLen(len(bys)))
	hex.Encode(dst, bys)
	return dst
}

func HexDecodeString(hexStr string) []byte {
	b, _ := hex.DecodeString(hexStr)
	return b
}

func BytesReverse(bytes []byte) []byte {
	ret := make([]byte, len(bytes))
	copy(ret, bytes)
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret
}

func Script2ScriptHash(script []byte) []byte {
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(script)
	pub_hash_1 := sha256_h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)

	return pub_hash_2
}

func ReverseBigLitterEndian(hexStr string) string {
	bys := []byte(hexStr)
	length := len(bys)
	for i := 0; i < length/2; i++ {
		if i%2 == 0 {
			bys[i], bys[length-2-i] = bys[length-2-i], bys[i]
		} else {
			bys[i], bys[length-i] = bys[length-i], bys[i]
		}
	}

	return string(bys)
}

func ScriptToHash(unhex []byte) string {
	sh := Script2ScriptHash(unhex)
	return ReverseBigLitterEndian(string(HexEncodeBytes(sh)))
}

func HexToUInt64(hexStr string) uint64 {
	bs := HexDecodeString(hexStr)
	for i := len(bs); i < 8; i++ {
		bs = append(bs, 0)
	}
	return binary.LittleEndian.Uint64(bs)
}

func Address2ScriptHash(address string) ([]byte, error) {
	_, data, err := Base58CheckDecode(address)
	if err != nil {
		return nil, err
	}
	return HexEncodeBytes(data), nil
}

func ScriptHash2Address(scriptHash []byte) (string, error) {
	length := len(scriptHash)
	if length != 20 {
		return "", fmt.Errorf("invalid scriptHash")
	}

	address := Base58CheckEncode(NEO_ADDRESS_VERSION, scriptHash)
	return address, nil
}

func FormatBigFloat(num string, base int, decimals int) (string, error) {
	f, _, err := new(big.Float).Parse(num, base)
	if err != nil {
		return "", nil
	}
	value := new(big.Float).Quo(f, big.NewFloat(math.Pow10(decimals)))
	return value.Text('f', -1), nil
}