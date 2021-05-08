package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"github.com/shopspring/decimal"
)

func IsNumber(text string) (decimal.Decimal, bool) {
	value, err := decimal.NewFromString(text)
	return value, nil == err
}

// Get the MD5 hash value of bytes
func BytesMD5Hash(data []byte) string {
	mac := hmac.New(md5.New, nil)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// Get the MD5 hash value of string
func StringMD5Hash(data string) string {
	return BytesMD5Hash([]byte(data))
}

// Get the MD5 hash value of multi string
func MultiStringMD5Hash(data ...string) string {
	var temp string
	for _, item := range data {
		temp += item
	}
	return StringMD5Hash(temp)
}
