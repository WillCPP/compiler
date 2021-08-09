// The entry point function is ScanFile
// The scanner opens a file and reads in
// the contents.  A token list is created
// where each token in the list is generated
// from the file's contents.

package scanner

import (
	"bufio"
	"compiler/src/types"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func OpenFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func CloseFile(file *os.File) {
	file.Close()
}

func ScanNextToken(byteScanner *bufio.Scanner, lineCounter *int, blockCommentCounter *int, atNextByte bool) (types.Token, types.ScanCode) {
	if !atNextByte {
		for {
			if byteScanner.Scan() {
				if byteScanner.Text() != " " && byteScanner.Text() != "\t" {
					break
				} else {
					continue
				}
			} else {
				return types.BuildToken(0, "", types.StopNoTokenScanCode)
			}
		}
	}

	word := ""
	firstByte := byteScanner.Text()
	firstByteAsRune, _ := utf8.DecodeRuneInString(firstByte)

	if firstByte == "\n" {
		return types.BuildToken(0, "", types.NewlineScanCode)
	} else if unicode.IsLetter(firstByteAsRune) {
		word += strings.ToLower(firstByte)
		for {
			if byteScanner.Scan() {
				nextByte := byteScanner.Text()
				nextByteAsRune, _ := utf8.DecodeRuneInString(nextByte)
				if unicode.IsLetter(nextByteAsRune) || unicode.IsDigit(nextByteAsRune) || nextByte == "_" {
					word += strings.ToLower(nextByte)
				} else if nextByte == " " || nextByte == "\t" {
					break
				} else if nextByte == "\n" {
					return types.BuildTokenFromAlphaNumeric(*lineCounter, word, types.AtNextByteScanCode)
				} else {
					return types.BuildTokenFromAlphaNumeric(*lineCounter, word, types.AtNextByteScanCode)
				}
			} else {
				return types.BuildTokenFromAlphaNumeric(*lineCounter, word, types.StopWithTokenScanCode)
			}
		}
		return types.BuildTokenFromAlphaNumeric(*lineCounter, word, types.TokenScanCode)
	} else if unicode.IsDigit(firstByteAsRune) {
		word += firstByte
		inDecimal := false
		for {
			if byteScanner.Scan() {
				nextByte := byteScanner.Text()
				nextByteAsRune, _ := utf8.DecodeRuneInString(nextByte)
				if unicode.IsDigit(nextByteAsRune) || nextByte == "_" {
					word += strings.ToLower(nextByte)
				} else if nextByte == "." && !inDecimal {
					word += strings.ToLower(nextByte)
					inDecimal = true
				} else if nextByte == " " || nextByte == "\t" {
					break
				} else if nextByte == "\n" {
					return types.BuildTokenFromNumeric(*lineCounter, word, types.AtNextByteScanCode)
				} else {
					return types.BuildTokenFromNumeric(*lineCounter, word, types.AtNextByteScanCode)
				}
			} else {
				return types.BuildTokenFromNumeric(*lineCounter, word, types.StopWithTokenScanCode)
			}
		}
		return types.BuildTokenFromNumeric(*lineCounter, word, types.TokenScanCode)
	} else if firstByte == "\"" {
		for {
			if byteScanner.Scan() {
				nextByte := byteScanner.Text()
				if nextByte == "\"" {
					return types.BuildTokenFromString(*lineCounter, word, types.TokenScanCode)
				} else if nextByte == "\n" {
					*lineCounter++
					word += nextByte
				} else {
					word += nextByte
				}
			} else {
				return types.BuildTokenFromString(*lineCounter, firstByte+word, types.ErrorScanCode)
			}
		}
	} else if firstByte == ":" {
		if byteScanner.Scan() {
			nextByte := byteScanner.Text()
			if nextByte == "=" {
				return types.BuildTokenFromSymbol(*lineCounter, firstByte+nextByte, types.TokenScanCode)
			}
		}
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	} else if firstByte == "<" {
		if byteScanner.Scan() {
			nextByte := byteScanner.Text()
			if nextByte == "=" {
				return types.BuildTokenFromSymbol(*lineCounter, firstByte+nextByte, types.TokenScanCode)
			}
		}
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	} else if firstByte == ">" {
		if byteScanner.Scan() {
			nextByte := byteScanner.Text()
			if nextByte == "=" {
				return types.BuildTokenFromSymbol(*lineCounter, firstByte+nextByte, types.TokenScanCode)
			}
		}
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	} else if firstByte == "=" {
		if byteScanner.Scan() {
			nextByte := byteScanner.Text()
			if nextByte == "=" {
				return types.BuildTokenFromSymbol(*lineCounter, firstByte+nextByte, types.TokenScanCode)
			}
		}
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	} else if firstByte == "!" {
		if byteScanner.Scan() {
			nextByte := byteScanner.Text()
			if nextByte == "=" {
				return types.BuildTokenFromSymbol(*lineCounter, firstByte+nextByte, types.TokenScanCode)
			}
		}
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	} else if firstByte == "/" {
		if byteScanner.Scan() {
			nextByte := byteScanner.Text()
			if nextByte == "/" {
				return types.BuildToken(*lineCounter, firstByte+nextByte, types.LineCommentScanCode)
			} else if nextByte == "*" {
				return types.BuildToken(*lineCounter, firstByte+nextByte, types.BlockCommentOpenScanCode)
			}
		}
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	} else if firstByte == "*" {
		if byteScanner.Scan() {
			nextByte := byteScanner.Text()
			if nextByte == "/" {
				return types.BuildToken(*lineCounter, firstByte+nextByte, types.BlockCommentCloseScanCode)
			}
		}
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	} else {
		return types.BuildTokenFromSymbol(*lineCounter, firstByte, types.TokenScanCode)
	}
}

func ScanErrorString(token types.Token) string {
	return "Error: Ln:" + fmt.Sprint(token.LineNumber)
}

func ScanFile(filename string) []types.Token {
	file := OpenFile(filename)
	defer CloseFile(file)
	blockCommentCounter := 0
	var tokenList []types.Token
	byteScanner := bufio.NewScanner(file)
	byteScanner.Split(bufio.ScanRunes)
	lineCounter := 1
	skipLine := false

	atNextByte := false
	for {
		token, code := ScanNextToken(byteScanner, &lineCounter, &blockCommentCounter, atNextByte)
		atNextByte = false
		if code == types.TokenScanCode {
			if !skipLine && blockCommentCounter == 0 {
				tokenList = append(tokenList, token)
			}
		} else if code == types.AtNextByteScanCode {
			if !skipLine && blockCommentCounter == 0 {
				tokenList = append(tokenList, token)
			}
			atNextByte = true
		} else if code == types.NewlineScanCode {
			lineCounter++
			if skipLine {
				skipLine = false
			}
		} else if code == types.StopWithTokenScanCode {
			if !skipLine && blockCommentCounter == 0 {
				tokenList = append(tokenList, token)
			}
			break
		} else if code == types.StopNoTokenScanCode {
			break
		} else if code == types.LineCommentScanCode {
			if blockCommentCounter == 0 {
				skipLine = true
			}
		} else if code == types.BlockCommentOpenScanCode {
			blockCommentCounter++
		} else if code == types.BlockCommentCloseScanCode {
			blockCommentCounter--
		} else if code == types.ErrorScanCode {
			if !skipLine && blockCommentCounter == 0 {
				fmt.Println(errors.New(ScanErrorString(token)))
				break
			}
		}
	}

	return tokenList
}

func PrintTokenList(tokenList []types.Token) {
	for _, value := range tokenList {
		if value.TokenType == types.IdentifierToken {
			print("Ln:", value.LineNumber, " | ", value.TokenType, " | ", value.StringValue)
		} else if value.TokenType == types.IntegerToken {
			print("Ln:", value.LineNumber, " | ", value.TokenType, " | ", value.IntValue)
		} else if value.TokenType == types.FloatToken {
			print("Ln:", value.LineNumber, " | ", value.TokenType, " | ", value.FloatValue)
		} else if value.TokenType == types.StringToken {
			print("Ln:", value.LineNumber, " | ", value.TokenType, " | ", value.StringValue)
		} else {
			print("Ln:", value.LineNumber, " | ", value.TokenType, " | ", value.StringValue)
		}
		print("\n")
	}
}
