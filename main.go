package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	var result string
	var traite string
	if len(os.Args[1:]) == 2 {
		in := os.Args[1]
		out := os.Args[2]
		stock, err := os.ReadFile(in)
		check(err, "cannot read file")
		split := strings.Split(string(stock), "\n")
		var finalrst string
		for x := 0; x < len(split); x++ {
			verifier := formatString(string(split[x]))
			finalsplit := strings.Split(string(verifier), " ")
			finalsplit = AtoAN(finalsplit)
			number := 0
			for i := 0; i < len(finalsplit); i++ {
				if strings.Contains(finalsplit[i], "(up,") && i >= 0 && i+1 < len(finalsplit) {
					number = checknum(finalsplit[i+1])
					if number == -1 {
						log.Fatal("make sure that you have a right flag ", finalsplit[i]+" ", finalsplit[i+1])
					}
					if number > len(finalsplit[:i+1])-1 {
						number = len(finalsplit[:i+1]) - 1
					}

					for j := 0; j < number && i-j-1 < len(finalsplit); j++ {
						if finalsplit[i-j-1] == "" {
							number++
						}

						if finalsplit[i-j-1] != ("") {
							finalsplit[i-j-1] = strings.ToUpper(finalsplit[i-j-1])
						}

						finalsplit[i] = ""
						if number != -1 {
							finalsplit[i+1] = ""
						}
					}
				} else if strings.Contains(finalsplit[i], "(low,") && i >= 0 && i+1 < len(finalsplit) {
					number = checknum(finalsplit[i+1])
					if number == -1 {
						log.Fatal("make sure that you have a right flag ", finalsplit[i]+" ", finalsplit[i+1])
					}
					if number > len(finalsplit[:i+1])-1 {
						number = len(finalsplit[:i+1]) - 1
					}

					for j := 0; j < number && i-j-1 < len(finalsplit); j++ {
						if finalsplit[i-j-1] == "" {
							number++
						}
						if finalsplit[i-j-1] != ("") {
							finalsplit[i-j-1] = strings.ToLower(finalsplit[i-j-1])
						}

						finalsplit[i] = ""
						if number != -1 {
							finalsplit[i+1] = ""
						}
					}
				} else if strings.Contains(finalsplit[i], "(hex") {
					if i > 0 {
						htodec, err := strconv.ParseInt(finalsplit[i-1], 16, 64)
						check(err, "You should use an hexadecimal number")
						finalsplit[i-1] = strconv.FormatInt(htodec, 10)
						finalsplit[i] = ""
					} else {
						finalsplit[i] = ""
					}
				} else if strings.Contains(finalsplit[i], "(bin") {
					if i > 0 {
						btodec, err := strconv.ParseInt(finalsplit[i-1], 2, 64)
						check(err, "You should use an binary number")
						finalsplit[i-1] = strconv.FormatInt(btodec, 10)
						finalsplit[i] = ""
					} else {
						finalsplit[i] = ""
					}
				} else if strings.Contains(finalsplit[i], "(cap,") && i >= 0 && i+1 < len(finalsplit) {
					number = checknum(finalsplit[i+1])
					if number == -1 {
						log.Fatal("make sure that you have a right flag ", finalsplit[i]+" ", finalsplit[i+1])
					}
					if number > len(finalsplit[:i+1])-1 {
						number = len(finalsplit[:i+1]) - 1
					}

					for j := 0; j < number && i-j-1 < len(finalsplit); j++ {
						if finalsplit[i-j-1] == "" {
							number++
						}
						if finalsplit[i-j-1] != ("") {
							finalsplit[i-j-1] = strings.ToLower(finalsplit[i-j-1])
							finalsplit[i-j-1] = strings.ToUpper(finalsplit[i-1][:1]) + finalsplit[i-j-1][1:]
						}

						finalsplit[i] = ""
						if number != -1 {
							finalsplit[i+1] = ""
						}

					}
				} else if strings.Contains(finalsplit[i], "(up)") && i>=0 && i+1 < len(finalsplit) {
					if finalsplit[i] != ("") {
						finalsplit[i-1] = strings.ToUpper(finalsplit[i-1])
						finalsplit[i] = ""
					}
					if finalsplit[i-1] == "" {
						number++
					} else {
						finalsplit = append(finalsplit[:i], finalsplit[i+1:]...)
					}
				} else if strings.Contains(finalsplit[i], "(low)") && i < len(finalsplit) {
					if finalsplit[i] != ("") {
						finalsplit[i-1] = strings.ToLower(finalsplit[i-1])
						finalsplit[i] = ""
					}
					if finalsplit[i-1] == "" {
						number++
					} else {
						finalsplit = append(finalsplit[:i], finalsplit[i+1:]...)
					}
				} else if strings.Contains(finalsplit[i], "(cap)") && i < len(finalsplit) {
					if finalsplit[i-1] != ("") {
						finalsplit[i-1] = strings.ToLower(finalsplit[i-1])
						finalsplit[i-1] = strings.ToUpper(finalsplit[i-1][:1]) + finalsplit[i-1][1:]
						finalsplit[i] = ""
					}
					if finalsplit[i-1] == "" {
						number++
					} else {
						finalsplit = append(finalsplit[:i], finalsplit[i+1:]...)
					}
				}
				var f string
				for i := 0; i < len(finalsplit); i++ {
					if finalsplit[i] != "" {
						f += finalsplit[i]
						if i < len(finalsplit)-1 {
							f += " "
						}
					}
				}
				input := f
				traite = GroupeofPunctuation(input)
				traite = Singlepunctuation(traite)
				testStrings := []string{traite}
				for _, str := range testStrings {
					result = SingleQuote(str)
				}
			}
			finalrst += result
			if x < len(split)-1 {
				finalrst += "\n"
			}
		}
		err = os.WriteFile(out, []byte(finalrst), 0644)
		check(err, "cannot write into the oufile")
	}
}



func GroupeofPunctuation(input string) string {
	corrected := regexp.MustCompile(`([.,!?:;])([a-z|A-Z|0-9])`).ReplaceAllString(input, "$1 $2")
	return corrected
}

func Singlepunctuation(input string) string {
	corrected := regexp.MustCompile(`\s+([.,!?:;])`).ReplaceAllString(input, "$1")
	return corrected
}

func checknum(s string) int {
	size := len(s)
	flag := false
	if s[size-1] == ')' {
		for i := 0; i < len(s)-1; i++ {
			if unicode.IsNumber(rune(s[i])) {
				flag = true
			} else {
				flag = false
				break
			}
		}
		if flag {
			num, _ := strconv.Atoi(string(s[:size-1]))
			return num
		}
	}
	return -1
}

func check(e error, msg string) {
	if e != nil {
		fmt.Println(msg)
		log.Fatal(e)
	}
}

func formatString(str string) string {
	str = strings.Join(strings.Fields(str), " ")
	return str
}

func isVowel(s string) bool {
	char := s[0]
	if char == 'a' || char == 'e' || char == 'u' || char == 'o' || char == 'i' || char == 'A' || char == 'E' || char == 'U' || char == 'O' || char == 'I' || char == 'h' || char == 'H' {
		return true
	}

	return false
}

func AtoAN(s []string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == "A" || s[i] == "a" {
			if i != len(s)-1 && isVowel(s[i+1]) {
				s[i] = s[i] + "n"
			}
		}
		if s[i] == "An" || s[i] == "an" {
			if i != len(s)-1 && !isVowel(s[i+1]) {
				s[i] = string(s[i][0])
			}
		}
	}
	return s
}


func SingleQuote(s string) string {
    rs := []rune(s)
    isfirstQuote := true
    for i := 0; i < len(rs); i++ {
        if isfirstQuote && i > 0 && (i < len(rs)-1) && (rs[i] == '\'' && (rs[i-1] == ' ' || rs[i+1] == ' ')) {
            if rs[i+1] == ' ' {
                rs = append(rs[:i+1], rs[i+2:]...)
                isfirstQuote = false
                continue
            } else {
                isfirstQuote = false
                continue
            }
        } else if !isfirstQuote && i != 0 && i != len(rs)-1 && rs[i] == '\'' && (rs[i-1] == ' ' || rs[i+1] == ' ') {
            if rs[i-1] == ' ' {
                rs = append(rs[:i-1], rs[i:]...)
                isfirstQuote = true
                continue
            } else {
                isfirstQuote = true
                continue
            }
        } else if i == 0 && (rs[i] == '\'') && isfirstQuote && i != len(rs)-1 {
            if rs[i+1] == ' ' {
                rs = append(rs[:i+1], rs[i+2:]...)
                isfirstQuote = false
                continue
            } else {
                isfirstQuote = false
                continue
            }
        } else if rs[i] == '\'' && !isfirstQuote && i == len(rs)-1 && i != 0 {
            if rs[i-1] == ' ' {
                rs = append(rs[:i-1], rs[i:]...)
                isfirstQuote = true
                continue
            } else {
                isfirstQuote = true
                continue
            }
        }
    }
    return string(rs)
}