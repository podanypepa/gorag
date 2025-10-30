package main

import "strings"

func splitText(text string, size int) []string {
	var chunks []string
	words := strings.Fields(text)
	for i := 0; i < len(words); i += size {
		end := i + size
		end = min(end, len(words))
		chunks = append(chunks, strings.Join(words[i:end], " "))
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
