package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"io"
)

func CheckDecrypted(data []byte) bool {
	index := bytes.Index(data, firmware_data.ChkDecryptedBytes)
	return index != -1
}

func CalcMD5(data []byte) string {
	// calc md5
	hash := md5.New()
	_, _ = io.Copy(hash, bytes.NewReader(data))
	return hex.EncodeToString(hash.Sum(nil))
}
