package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GeneratePasswordSalt() string {
	ts := time.Now().Unix()
	str := fmt.Sprintf("%x", ts)
	return fmt.Sprintf("%x\n", sha256.Sum256([]byte(str)))
}

func EncryptPassword(password string, salt string) string {
	return fmt.Sprintf("%x\n", sha256.Sum256([]byte(password+salt)))
}

func GenerateToken(username string) string {
	str := fmt.Sprintf("%x", time.Now().Unix())
	_md5 := md5.New()
	_md5.Write([]byte(username + str))
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum(nil))
}
