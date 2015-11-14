package terse

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"net/url"
)

// URL safe characters, as definied by section 2.3 of RFC 3986 (http://www.ietf.org/rfc/rfc3986.txt)
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

func CleanURL(rawurl string) (string, error) {
	parsed, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return "", fmt.Errorf("Error parsing url \"%s\": %s", rawurl, err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("Scheme \"%s\" not allowed for url \"%s\"", parsed.Scheme, rawurl)
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("No hostname provided for url \"%s\"", rawurl)
	}

	return parsed.String(), nil
}
