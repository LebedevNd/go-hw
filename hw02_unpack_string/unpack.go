package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	var resultStringBuilder strings.Builder
	resultString, e := iterateString(inputString, &resultStringBuilder)

	return resultString, e
}

func iterateString(mainString string, resultStringBuilder *strings.Builder) (string, error) {
	if mainString == "" {
		return resultStringBuilder.String(), nil
	}

	if len(mainString) == 1 {
		if unicode.IsDigit(rune(mainString[0])) {
			number, _ := strconv.Atoi(mainString[0:1])
			return processingNumber(number, mainString, resultStringBuilder)
		}
		resultStringBuilder.WriteString(mainString)
		return resultStringBuilder.String(), nil
	}

	firstChar := mainString[0:1]
	secondChar := mainString[1:2]

	firstCharNum, firstCharError := strconv.Atoi(firstChar)
	secondCharNum, secondCharError := strconv.Atoi(secondChar)

	if firstChar == `\` && (secondChar == `\` || secondCharError == nil) {
		resultStringBuilder.WriteString(secondChar)
		return iterateString(mainString[2:], resultStringBuilder)
	}

	if (firstCharError == nil && secondCharError == nil) || firstChar == `\` {
		return "", ErrInvalidString
	}

	if firstCharError != nil && secondCharError != nil {
		resultStringBuilder.WriteString(firstChar)
		return iterateString(mainString[1:], resultStringBuilder)
	}

	if firstCharError == nil {
		return processingNumber(firstCharNum, mainString, resultStringBuilder)
	}

	if secondCharError == nil {
		additionalString := strings.Repeat(firstChar, secondCharNum)
		resultStringBuilder.WriteString(additionalString)
	} else {
		resultStringBuilder.WriteString(firstChar + secondChar)
	}

	return iterateString(mainString[2:], resultStringBuilder)
}

func processingNumber(number int, mainString string, resultStringBuilder *strings.Builder) (string, error) {
	resultString := resultStringBuilder.String()
	var lastChar rune
	if len(resultString) > 0 {
		lastChar = rune(resultString[len(resultString)-1:][0])
		if string(lastChar) == `\` || unicode.IsDigit(lastChar) {
			if number > 0 {
				number--
			}
			additionalString := strings.Repeat(string(lastChar), number)
			resultStringBuilder.WriteString(additionalString)
			return iterateString(mainString[1:], resultStringBuilder)
		}
	}
	return "", ErrInvalidString
}
