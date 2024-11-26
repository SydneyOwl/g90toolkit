package tools

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"io"
	"os"
	"strings"
)

func encrypt(deckey string, src []byte) ([]byte, error) {
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

func DoEncryptAndSave(key string, firmwarePath string, outputPath string) error {
	md5sum := CalcMD5([]byte(key))
	if md5sum != firmware_data.KnownKeyMD5 {
		fmt.Println("WARNING: THE KEY PROVIDED MISMATCH WITH KNOWN KEY. USE AT YOUR OWN RISK.")
	}
	data, err := os.ReadFile(firmwarePath)
	if err != nil {
		return fmt.Errorf("failed to read firmware: %v\n", err)
	}
	if !CheckDecrypted(data) {
		fmt.Println("Seems like the firmware you provided is already encrypted. Do you wish to continue? [y/n]")
		choice := "n"
		_, _ = fmt.Scanln(&choice)
		if strings.ToUpper(choice) != "Y" {
			return fmt.Errorf("user aborted")
		}
	}
	enc, err := encrypt(key, data)
	if err != nil {
		return fmt.Errorf("failed to Encrypt firmware: %v\n", err)
	}
	err = os.WriteFile(outputPath, enc, 0777)
	if err != nil {
		return fmt.Errorf("failed to write to output: %v\n", err)
	}
	return nil
}

func DoEncrypt(key string, firmware []byte) ([]byte, error) {
	md5sum := CalcMD5([]byte(key))
	if md5sum != firmware_data.KnownKeyMD5 {
		fmt.Println("WARNING: THE KEY PROVIDED MISMATCH WITH KNOWN KEY. USE AT YOUR OWN RISK.")
	}
	if !CheckDecrypted(firmware) {
		fmt.Println("Seems like the firmware you provided is already encrypted. Do you wish to continue? [y/n]")
		choice := "n"
		_, _ = fmt.Scanln(&choice)
		if strings.ToUpper(choice) != "Y" {
			return nil, fmt.Errorf("user aborted")
		}
	}
	enc, err := encrypt(key, firmware)
	if err != nil {
		return nil, fmt.Errorf("failed to Encrypt firmware: %v\n", err)
	}
	return enc, nil
}
