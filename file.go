package main

import (
	"io"
	"os"
	"strings"
)

func DeepCopyFile(src, dst string) error {
	//create directory
	parts := strings.Split(dst, string(os.PathSeparator))
	path := parts[0]
	for i := 1; i < len(parts)-1; i++ {
		path = path + string(os.PathSeparator) + parts[i]
		if !Exists(path) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return CopyFile(src, dst)
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// copy file
// if file is exists, overwrite it.
// if directory is not exists, create it.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	out.Truncate(0)
	out.Seek(0, 0)

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
