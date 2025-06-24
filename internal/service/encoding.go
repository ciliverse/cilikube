package service

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GBKToUTF8Reader wraps an io.Reader and converts GBK bytes to UTF-8
func GBKToUTF8Reader(r io.Reader) io.Reader {
	return transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
}

// DetectGBK tries to detect if the first N bytes are GBK encoded (very simple heuristic)
func DetectGBK(data []byte) bool {
	// If data contains bytes in 0x81-0xFE range, likely GBK
	for _, b := range data {
		if b >= 0x81 && b <= 0xFE {
			return true
		}
	}
	return false
}

// ConvertIfGBK wraps a log stream, if it looks like GBK, convert to UTF-8
func ConvertIfGBK(logStream io.ReadCloser) io.ReadCloser {
	peek := make([]byte, 512)
	n, _ := logStream.Read(peek)
	if n == 0 {
		return logStream
	}
	if DetectGBK(peek[:n]) {
		utf8Reader := GBKToUTF8Reader(io.MultiReader(bytes.NewReader(peek[:n]), logStream))
		return ioutil.NopCloser(utf8Reader)
	}
	return ioutil.NopCloser(io.MultiReader(bytes.NewReader(peek[:n]), logStream))
}
