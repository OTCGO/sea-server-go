package neo

import (
	"crypto/sha256"
	"encoding/hex"
)

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
