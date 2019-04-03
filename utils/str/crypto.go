package str

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

func MD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func Sha256(bytes []byte) string {
	h := sha256.New()
	h.Write(bytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}
