package main

import "strings"

func splitText(text string, size, overlap int) []string {
	var chunks []string
	words := strings.Fields(text)
	if size <= 0 {
		return chunks
	}
	if overlap >= size {
		overlap = size / 2
	}

	for i := 0; i < len(words); {
		end := i + size
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
		if end == len(words) {
			break
		}
		i += (size - overlap)
	}
	return chunks
}

func estimateTokens(s string) int {
	words := len(strings.Fields(s))
	runes := len([]rune(s))
	wTok := int(float64(words) / 0.75)
	cTok := runes / 4
	if wTok > cTok {
		return wTok
	}
	return cTok
}
