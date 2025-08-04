package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encipher(input string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(input))
	if err != nil {
		panic(err.Error())
	}

	hashedBytes := hasher.Sum(nil)            // 计算哈希摘要
	md5Str := hex.EncodeToString(hashedBytes) // 将摘要转换成十六进制格式的字符串
	return md5Str
}
