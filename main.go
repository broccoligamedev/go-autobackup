package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	var directoryToZip string
	if len(os.Args) > 1 {
		directoryToZip = os.Args[1]
	} else {
		directoryToZip = "default"
	}
	output := "backup.zip"
	err := zipFiles(directoryToZip, output)
	if err != nil {
		panic(err)
	}
	fmt.Println("Zipped file: ", output)
}

func zipFiles(directoryToZip, output string) error {
	newZipFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	filepath.Walk(directoryToZip, func(path string, info os.FileInfo, err error) error {
		fmt.Println("Found ", path)
		if info.IsDir() {
			return nil
		}
		fileToZip, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToZip.Close()
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = path
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}
