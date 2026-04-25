package main

import (
	"bytes"
	"io"
	"os"

	pdf "github.com/ledongthuc/pdf"
)

func ExtractText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	pr, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(&buf, pr); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ExtractMDText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
