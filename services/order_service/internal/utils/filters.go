package utils

import (
	"bufio"
	"fmt"
	lg "github.com/randnull/Lessons/internal/logger"
	"os"
	"strings"
)

var badWords map[string]bool

func LoadBadWords() error {
	badWords = make(map[string]bool)

	file, err := os.Open("./ban_words.txt")
	if err != nil {
		lg.Error(fmt.Sprintf("[Filters] Error init filters: %v", err.Error()))
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		lg.Error(fmt.Sprintf("[Filters] Add word: %v", word))
		badWords[word] = true
	}

	return scanner.Err()
}

func ContainsBadWords(text string) bool {
	text = strings.ToLower(text)
	words := strings.Fields(text)

	for _, word := range words {
		_, exists := badWords[word]

		if exists {
			lg.Info(fmt.Sprintf("[Filters] Detected ban word: %v", word))
			return true
		}
	}
	return false
}
