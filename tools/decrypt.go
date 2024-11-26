package tools

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"os"
	"strings"
)

func decrypt(deckey string, src []byte) ([]byte, error) {
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

func DoDecryptAndSave(key string, firmwarePath string, outputPath string) error {
	md5sum := CalcMD5([]byte(key))
	if md5sum != firmware_data.KnownKeyMD5 {
		fmt.Println("WARNING: THE KEY PROVIDED MISMATCH WITH KNOWN KEY. USE AT YOUR OWN RISK.")
	}
	data, err := os.ReadFile(firmwarePath)
	if err != nil {
		return fmt.Errorf("failed to read firmware: %v\n", err)
	}
	if CheckDecrypted(data) {
		fmt.Println("Seems like the firmware you provided is already decrypted. Do you wish to continue? [y/n]")
		choice := "n"
		_, _ = fmt.Scanln(&choice)
		if strings.ToUpper(choice) != "Y" {
			return fmt.Errorf("user aborted.")
		}
	}
	dec, err := decrypt(key, data)
	if err != nil {
		return fmt.Errorf("failed to decrypt firmware: %v\n", err)
	}
	err = os.WriteFile(outputPath, dec, 0777)
	if err != nil {
		return fmt.Errorf("failed to write to output: %v\n", err)
	}
	return nil
}

func DoDecrypt(key string, firmware []byte) ([]byte, error) {
	md5sum := CalcMD5([]byte(key))
	if md5sum != firmware_data.KnownKeyMD5 {
		fmt.Println("WARNING: THE KEY PROVIDED MISMATCH WITH KNOWN KEY. USE AT YOUR OWN RISK.")
	}
	if CheckDecrypted(firmware) {
		fmt.Println("Seems like the firmware you provided is already decrypted. Do you wish to continue? [y/n]")
		choice := "n"
		_, _ = fmt.Scanln(&choice)
		if strings.ToUpper(choice) != "Y" {
			return nil, fmt.Errorf("user aborted.")
		}
	}
	dec, err := decrypt(key, firmware)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt firmware: %v\n", err)
	}
	return dec, nil
}
