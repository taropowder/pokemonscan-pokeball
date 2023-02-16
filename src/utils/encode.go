package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func DecodeB64(message string) (retour string, err error) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	_, err = base64.StdEncoding.Decode(base64Text, []byte(message))
	if err != nil {
		return "", err
	}
	//fmt.Printf("base64: %s\n", base64Text)
	return string(base64Text), nil
}

func EncodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(base64Text, []byte(message))
	return string(base64Text)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
