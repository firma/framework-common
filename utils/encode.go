package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
	"strings"
)

func EncodeString(source string) string {
	src := []byte(source)
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	str := base64.StdEncoding.EncodeToString(in.Bytes())
	return encodeReplace(str)
}

func DecodeString(source string) string {
	source = dencodeReplace(source)
	compressSrc, _ := base64.StdEncoding.DecodeString(source)
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return string(out.String())
}

func encodeReplace(encode string) string {
	return strings.Replace(encode, "/", "_", -1)
}

func dencodeReplace(encode string) string {
	return strings.Replace(encode, "_", "/", -1)
}
