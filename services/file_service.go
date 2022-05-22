package services

import (
	"goblog/util"
	"io/ioutil"
)

func GetFiles(path string) []map[string]string {
	var resMap []map[string]string
	fileInfos, _ := ioutil.ReadDir(path)

	for _, fi := range fileInfos {
		if !fi.IsDir() {
			name := fi.Name()
			size := fi.Size()
			m := map[string]string{
				"name": name,
				"size": util.FileSizeToString(size),
			}
			resMap = append(resMap, m)
		}
	}
	return resMap
}
