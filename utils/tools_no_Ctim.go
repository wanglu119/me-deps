//go:build darwin || freebsd || netbsd || openbsd
// +build darwin freebsd netbsd openbsd

package utils

import (
	"os"
)

/**
* @Description:
* @param path
* @return int64
 */
func GetFileCreateTime(path string) int64 {
	fileInfo, _ := os.Stat(path)
	tCreate := fileInfo.ModTime().Unix()
	return tCreate
}
