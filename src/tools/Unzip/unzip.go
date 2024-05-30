package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipUtil struct{}

// UnzipFile 解压到本地文件中
func (z *ZipUtil) UnzipFile(zipFile, dest string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		filePath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, f.Mode())
			if err != nil {
				return err
			}
		} else {
			var dir string
			if lastIndex := strings.LastIndex(filePath, string(os.PathSeparator)); lastIndex > -1 {
				dir = filePath[:lastIndex]
				err := os.MkdirAll(dir, f.Mode())
				if err != nil {
					return err
				}
			}

			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
