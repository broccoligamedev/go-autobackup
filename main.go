package main

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		panic(errors.New("wrong number of arguments: <directory to zip> <directory to save to> <frequency in minutes>"))
	}
	thingToZip := os.Args[1]
	saveDirectory := os.Args[2]
	frequencyArg, err := strconv.Atoi(os.Args[3])
	if err != nil || frequencyArg <= 0 {
		panic(errors.New("frequency should be an integer"))
	}
	frequency := time.Duration(frequencyArg) * time.Minute
	// note(ryan): this time format isn't random gibberish. Go just uses a weird formatting system.
	// we are assuming these will be unique so we don't worry about overwriting
	log.Println("starting")
	log.Println("saving \"" + thingToZip + "\" to directory \"" + saveDirectory + "\" every " + os.Args[3] + " minute(s)")
	timer := time.NewTimer(frequency)
	for {
		<-timer.C
		timer.Reset(frequency)
		timestamp := time.Now().Format("02012006030405")
		outputName := saveDirectory + "autobackup" + timestamp + ".zip"
		err := zipFiles(thingToZip, outputName)
		if err != nil {
			log.Println("can't save file: " + err.Error())
			continue
		}
		log.Println("saved file to: " + outputName)
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
		log.Println("Found ", path)
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
