package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
