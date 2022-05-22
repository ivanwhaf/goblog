package util

import (
	"math"
	"os"
	"strconv"
)

func FileSizeToString(size int64) string {
	var res string
	if float64(size) < (math.Pow(2, 20)) {
		res = strconv.FormatFloat(float64(size)/math.Pow(2, 10), 'f', 2, 64) + "KB"
	} else if float64(size) < (math.Pow(2, 30)) {
		res = strconv.FormatFloat(float64(size)/math.Pow(2, 20), 'f', 2, 64) + "MB"
	} else {
		res = strconv.FormatFloat(float64(size)/math.Pow(2, 20), 'f', 2, 64) + "GB"
	}
	return res
}

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
