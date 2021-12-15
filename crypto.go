package main

import (
	"crypto/des"
	"crypto/sha1"
	"encoding/hex"
)

func hashPassw(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	bs := hex.EncodeToString(h.Sum(nil))
	return bs
}

func createDesKey(key []byte) []byte {
	if len(key) >= 8 {
		return key[:8]
	}
	length := len(key)
	for i := length; i < 8; i++ {
		key = append(key, key[i%length])
	}
	return key[:8]
}

func preparePlaintext(plaintext string) []byte {
	plain := []byte(plaintext)
	if len(plaintext) < 8 {
		for i := 0; i < 8-len(plaintext); i++ {
			plain = append([]byte{0}, plain...)
		}
	}
	return plain
}

func encryptDES(key []byte, plaintext string) string {
	// create cipher
	key = createDesKey(key)
	c, _ := des.NewCipher(key)
	// allocate space for ciphered data
	plain := preparePlaintext(plaintext)

	out := make([]byte, len(plain))

	// encrypt
	c.Encrypt(out, plain)
	// return hex string
	return hex.EncodeToString(out)
}

func decryptDES(key []byte, ct string) string {
	key = createDesKey(key)
	ciphertext, _ := hex.DecodeString(ct)

	c, _ := des.NewCipher(key)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	for i, p := range pt {
		if p != 0 {
			return string(pt[i:])
		}
	}
	return string(pt)
}
