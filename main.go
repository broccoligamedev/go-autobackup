package main

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		panic(errors.New("wrong number of arguments: <directory to zip> <directory to save to> <frequency in minutes>"))
	}
	thingToZip := os.Args[1]
	saveDirectory := os.Args[2]
	// note(ryan): the .zip specification requires directories use forward slashes, but windows uses backslashes.
	// this creates an issue where extracting files on linux leads to file names like "these\should\be\directories\test.txt"
	// instead of actual directories.
	saveDirectory = strings.ReplaceAll(saveDirectory, "\\", "/")
	// note(ryan): if our last character isn't a forward slash then we need to insert you to get the correct path
	if saveDirectory[len(saveDirectory)-1] != '/' {
		saveDirectory = saveDirectory + "/"
	}
	frequencyArg, err := strconv.Atoi(os.Args[3])
	if err != nil || frequencyArg <= 0 {
		panic(errors.New("frequency should be an integer greater than 0"))
	}
	frequency := time.Duration(frequencyArg) * time.Minute
	// note(ryan): this time format isn't random gibberish. Go just uses a weird formatting system.
	// we are assuming these will be unique so we don't worry about overwriting
	log.Println("starting")
	log.Println("saving \"" + thingToZip + "\" to directory \"" + saveDirectory + "\" every " + os.Args[3] + " minute(s)")
	timer := time.NewTimer(frequency)
	fileCount := 0
	errorCount := 0
	for {
		timestamp := time.Now().Format("02012006030405")
		outputName := saveDirectory + "autobackup" + timestamp + ".zip"
		err := zipFiles(thingToZip, outputName)
		fileCount++
		if err != nil {
			errorCount++
			continue
		}
		log.Println("saved file to: " + outputName)
		log.Println(strconv.Itoa(errorCount) + " error(s) out of " + strconv.Itoa(fileCount) + " files")
		<-timer.C
		timer.Reset(frequency)
	}
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
		header.Name = strings.ReplaceAll(path, "\\", "/")
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
