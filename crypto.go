package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// ZoomContext ...
type ZoomContext struct {
	Type      string `json:"typ"`
	UserID    string `json:"uid"`
	MeetingID string `json:"mid"`
	Nonce     string `json:"dev"`
	Timestamp string `json:"ts"`
}

func unpackContext(context string) (iv, aad, tag, cipherText []byte, err error) {
	decoded, err := base64.URLEncoding.DecodeString(context)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	buf := bytes.NewBuffer(decoded)

	ivLength := int(buf.Next(1)[0])
	iv = buf.Next(ivLength)

	aadLength := int(binary.LittleEndian.Uint16(buf.Next(2)))
	aad = buf.Next(aadLength)

	cipherLength := int(binary.LittleEndian.Uint32(buf.Next(4)))
	cipher := buf.Next(cipherLength)

	tag, err = buf.ReadBytes(0x00)

	return iv, aad, cipher, tag, nil
}

func hashSecret(secret string) []byte {
	h := sha256.New()
	h.Write([]byte(secret))

	return h.Sum(nil)
}

func decrypt(key, iv, ciphertext, aad, tag []byte) (result []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, iv, []byte(string(ciphertext)+string(tag)), nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func decryptZoomContext(secret, context string) (meetingContext *ZoomContext, err error) {
	key := hashSecret(secret)
	iv, aad, cipher, tag, err := unpackContext(context)

	if err != nil {
		return nil, err
	}

	plaintext, err := decrypt(key, iv, cipher, aad, tag)

	if err != nil {
		return nil, err
	}

	res := new(ZoomContext)

	json.Unmarshal(plaintext, res)

	fmt.Println(string(plaintext))
	fmt.Println(res.Timestamp)

	return res, nil
}
