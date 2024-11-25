package tools

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"io"
)

func DoDecrypt(deckey string, src []byte) ([]byte, error) {
	key, err := hex.DecodeString(deckey)
	if err != nil {
		return nil, fmt.Errorf("could not parse provided key as a 256-bit hexadecimal number")
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 256 bits (32 bytes)")
	}
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %v", err)
	}
	dst := make([]byte, 0)
	tmp := make([]byte, cipherBlock.BlockSize())
	// TODO: MAKE IT MORE ELEGANT
	for index := 0; index < len(src); index += cipherBlock.BlockSize() {
		cipherBlock.Decrypt(tmp, src[index:index+cipherBlock.BlockSize()])
		dst = append(dst, tmp...)
	}
	return dst, nil
}

func DoEncrypt(deckey string, src []byte) ([]byte, error) {
	key, err := hex.DecodeString(deckey)
	if err != nil {
		return nil, fmt.Errorf("could not parse provided key as a 256-bit hexadecimal number")
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 256 bits (32 bytes)")
	}
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %v", err)
	}
	blockSize := cipherBlock.BlockSize()
	buffer := make([]byte, blockSize)
	encrypted := make([]byte, blockSize)
	dst := make([]byte, 0)
	lines := bytes.NewReader(src)
	for {
		n, err := lines.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading input file: %v", err)
		}

		if n == 0 {
			break
		}

		if n < blockSize {
			padding := bytes.Repeat([]byte{0}, blockSize-n)
			buffer = append(buffer[:n], padding...)
		}

		cipherBlock.Encrypt(encrypted, buffer)
		dst = append(dst, encrypted...)
	}
	return dst, nil
}
