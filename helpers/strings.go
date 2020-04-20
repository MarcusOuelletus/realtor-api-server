package helpers

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
)

func StartsWith(fullstring, substring string) bool {
	var b []byte
	valueLength := len(fullstring)
	fullstring = strings.ToLower(fullstring)
	substring = strings.ToLower(substring)

	for i, character := range substring {
		if valueLength == i {
			return false
		}

		b = []byte{fullstring[i]}

		if !bytes.ContainsRune(b, character) {
			return false
		}
	}

	return true
}

func GenerateToken(length int) string {

	var keySymbols = [62]string{
		"a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l",
		"m", "n", "o", "p", "q", "r",
		"s", "t", "u", "v", "w", "x",
		"y", "z", "A", "B", "C", "D",
		"E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U", "V",
		"W", "X", "Y", "Z", "1", "2",
		"3", "4", "5", "6", "7", "8",
		"9", "0",
	}

	rand.Seed(time.Now().Unix())

	var token string
	var symbolIndex int

	for i := 0; i < length; i++ {
		symbolIndex = rand.Int() % len(keySymbols)

		token += keySymbols[symbolIndex]
	}

	return token
}
