package app

import (
	"os"
	"strings"
)

func (a *App) ScanDir(pwd string) (dirList []string) {
	fileInfos, err := os.ReadDir(pwd)
	if err != nil {
		return
	}
	dirList = make([]string, 0)
	for _, fi := range fileInfos {
		if fi.IsDir() {
			iPwd := pwd + "/" + fi.Name()
			if fi.Name() != "converter" && a.isWatchDir(iPwd) {
				dirList = append(dirList, iPwd)
			} else { //多层结构
				dirList = append(dirList, a.ScanDir(iPwd)...)
			}
		}
	}
	return dirList
}

func (a *App) isWatchDir(pwd string) bool {
	fileInfos, err := os.ReadDir(pwd)
	if err != nil {
		return false
	}
	for _, fi := range fileInfos {
		if strings.HasSuffix(fi.Name(), ".go") && !fi.IsDir() {
			return true
		}
	}
	return false
}
