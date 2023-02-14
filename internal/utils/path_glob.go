package utils

import "regexp"

var regexMap = make(map[string]*regexp.Regexp)

var regexSingleChar = regexp.MustCompile("\\?")
var regexMultipleChar = regexp.MustCompile("\\*")

func MatchGlob(file string, glob string) (bool, error) {
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
}

func convert(glob string) string {
	glob = regexSingleChar.ReplaceAllString(glob, ".")
	glob = regexMultipleChar.ReplaceAllString(glob, ".*")
	return glob
}
