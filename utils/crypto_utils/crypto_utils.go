package crypto_utils

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func GetMd5(input string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetBcrypt(input string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hash)
}
