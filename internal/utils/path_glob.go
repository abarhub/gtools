package utils

import (
	"github.com/bmatcuk/doublestar/v4"
	"path/filepath"
	"regexp"
)

var regexMap = make(map[string]*regexp.Regexp)

var regexDot = regexp.MustCompile("\\.")
var regexSingleChar = regexp.MustCompile("\\?")
var regexMultipleChar = regexp.MustCompile("\\*")
var regexMultipleDirectory = regexp.MustCompile("\\*\\*/")

func MatchGlob(file string, glob string) (bool, error) {
	if false {
		var regexGlob *regexp.Regexp
		if regexMap[glob] == nil {
			globConverted := convert(glob)
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
	} else {
		glob2 := filepath.ToSlash(glob)
		file2 := filepath.ToSlash(file)
		return doublestar.PathMatch(glob2, file2)
	}
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
