package terse

import (
	"bytes"
	"hash/crc32"
)

const ALPHABET string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
const BASE uint32 = uint32(len(ALPHABET))

func GetShortCode(url []byte) string {
	var code bytes.Buffer
	num := crc32.ChecksumIEEE(url)
	for num > 0 {
		next := (num % BASE)
		code.WriteRune(rune(ALPHABET[next]))
		num = num / 62
	}
	return code.String()
}
