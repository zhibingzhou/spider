package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RandSleep(sec int) {
	time.Sleep(time.Duration(rand.Intn(sec)) * time.Second)
}

func ZhToUnicode(sText string) ([]byte, error) {
	textQuoted := strconv.QuoteToASCII(sText)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(textUnquoted)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
