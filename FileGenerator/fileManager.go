package FileGenerator

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetFileStr(name string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Cannot get caller information")
	}

	// Get the directory of this file
	dir := filepath.Dir(filename) + "/GoFiles/"
	bts, err := os.ReadFile(dir + name)
	if err != nil {
		panic(err)
	}

	return string(bts)
}

func GetComplexFileStr(name string, args map[string]string) string {
	fileStr := GetFileStr(name)

	for key, value := range args {
		fileStr = strings.Replace(fileStr, "$"+key+"$", value, -1)
	}

	return fileStr
}
