package main

import (
	"bufio"
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
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		buf.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
