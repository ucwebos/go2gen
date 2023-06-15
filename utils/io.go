package utils

import (
	"io/ioutil"
	"strings"
)

func ScanPbDir(pwd string) (pbFlies []string) {
	fileInfos, err := ioutil.ReadDir(pwd)
	if err != nil {
		return
	}
	pbFlies = make([]string, 0)
	for _, fi := range fileInfos {
		iPwd := pwd + "/" + fi.Name()
		if fi.IsDir() {
			pbFlies = append(pbFlies, ScanPbDir(iPwd)...)
		} else {
			if strings.HasSuffix(iPwd, ".proto") || strings.HasSuffix(iPwd, ".proto3") {
				pbFlies = append(pbFlies, iPwd)
			}
		}
	}
	return pbFlies
}
