package filepath

import (
	"os"
	"strings"
)

func FileExistOrCreate(path string) bool {
	pathClips := strings.Split(path, "/")
	dir := strings.Join(pathClips[0:len(pathClips)-1], "/")
	if !IsExist(path) {
		os.MkdirAll(dir, os.ModePerm)
		os.Create(path)
	}
	return true
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

func CreateFile(path string) {
	pathClips := strings.Split(path, "/")
	dir := strings.Join(pathClips[0:len(pathClips)-1], "/")
	os.Mkdir(dir, os.ModePerm)
	os.Create(path)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}
