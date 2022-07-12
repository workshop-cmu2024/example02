package sonicwebsocket

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
)

func mask(mask, b []byte) {
	for i := range b {
		b[i] ^= mask[i&3]
	}
}

func genMask(b []byte) {
	rand.Read(b)
}

func makeRequestKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func makeResponseKey(reqKey []byte) string {
	var resKey []byte
	resKey = append(resKey, reqKey...)
	resKey = append(resKey, GUID...)

	hasher := sha1.New()
	hasher.Write(resKey)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}