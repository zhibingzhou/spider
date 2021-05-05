package utils

import (
	"math/rand"
	"reflect"
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

// 利用反射将化为map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
