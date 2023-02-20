package utils

import (
	"fmt"
	"regexp"
)

var regexMap = make(map[string]*regexp.Regexp)

var regexDot = regexp.MustCompile("\\.")
var regexSingleChar = regexp.MustCompile("\\?")
var regexMultipleChar = regexp.MustCompile("\\*")
var regexMultipleDirectory = regexp.MustCompile("\\*\\*/")

func MatchGlob(file string, glob string) (bool, error) {
	var regexGlob *regexp.Regexp
	if regexMap[glob] == nil {
		globConverted := convert(glob)
		fmt.Println("convert=", globConverted, glob)
		r, err := regexp.Compile(globConverted)
		if err != nil {
			return false, err
		}
		regexMap[glob] = r
		regexGlob = r
	} else {
		regexGlob = regexMap[glob]
	}
	return regexGlob.MatchString(file), nil
}

func convert(glob string) string {
	glob2 := ""
	for i := 0; i < len(glob); i++ {
		c := glob[i]
		s := ""
		if c == '.' {
			s = "\\."
		} else if c == '?' {
			s = "[^\\/]"
		} else if c == '*' {
			if i+2 < len(glob) && glob[i+1] == '*' && glob[i+2] == '/' {
				s = "([^\\\\/]*/)*"
				i += 2
			} else {
				s = "[^\\\\/]*"
			}
		} else {
			s = string(c)
		}
		glob2 += s
	}
	return glob2
}
