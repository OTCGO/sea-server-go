package neo

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

const BitcoinBase58Table = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// b58encode encodes a byte slice b into a base-58 encoded string.
func b58encode(b []byte) (s string) {
	x := new(big.Int).SetBytes(b)

	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	s = ""

	for x.Cmp(zero) > 0 {
		x.QuoRem(x, m, r)
		s = string(BitcoinBase58Table[r.Int64()]) + s
	}

	return s
}

// b58decode decodes a base-58 encoded string into a byte slice b.
func b58decode(s string) (b []byte, err error) {
	x := big.NewInt(0)
	m := big.NewInt(58)

	for i := 0; i < len(s); i++ {
		b58index := strings.IndexByte(BitcoinBase58Table, s[i])
		if b58index == -1 {
			return nil, fmt.Errorf("invalid base58 character encountered: '%c', index %d", s[i], i)
		}
		b58value := big.NewInt(int64(b58index))
		x.Mul(x, m)
		x.Add(x, b58value)
	}
	return x.Bytes(), nil
}

// Base58CheckEncode encodes version ver and byte slice b into a base-58 check encoded string.
func Base58CheckEncode(ver uint8, b []byte) (s string) {
	bcpy := append([]byte{ver}, b...)

	h := sha256.New()

	h.Reset()
	h.Write(bcpy)
	hash1 := h.Sum(nil)

	h.Reset()
	h.Write(hash1)
	hash2 := h.Sum(nil)

	bcpy = append(bcpy, hash2[0:4]...)

	s = b58encode(bcpy)

	for _, v := range bcpy {
		if v != 0 {
			break
		}
		s = "1" + s
	}

	return s
}

// Base58CheckDecode decodes base-58 check encoded string s into a version ver and byte slice b.
func Base58CheckDecode(s string) (ver uint8, b []byte, err error) {
	b, err = b58decode(s)
	if err != nil {
		return 0, nil, err
	}

	for i := 0; i < len(s); i++ {
		if s[i] != '1' {
			break
		}
		b = append([]byte{0x00}, b...)
	}

	if len(b) < 5 {
		return 0, nil, fmt.Errorf("invalid base58 check string: missing checksum")
	}

	h := sha256.New()

	h.Reset()
	h.Write(b[:len(b)-4])
	hash1 := h.Sum(nil)

	h.Reset()
	h.Write(hash1)
	hash2 := h.Sum(nil)

	if bytes.Compare(hash2[0:4], b[len(b)-4:]) != 0 {
		return 0, nil, fmt.Errorf("invalid base58 check string: invalid checksum")
	}

	b = b[:len(b)-4]

	ver = b[0]
	b = b[1:]

	return ver, b, nil
}
