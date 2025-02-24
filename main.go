package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"rajatsharma.dev/roadrunner/parser"
)

func processFile(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(src)

	fmt.Println("Processing file:", destPath)

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func(src *os.File) {
		err := dest.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(src)

	err = parser.ConvertToTS(src, dest)
	if err != nil {
		return err
	}

	fmt.Println("Processed file:", destPath)
	return nil
}

func processDirectory(srcDir, destDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", path)
			return err
		}

		if info.IsDir() {
			if path == srcDir {
				return nil
			}

			err = os.MkdirAll(filepath.Join(destDir, info.Name()), os.ModePerm)
			if err != nil {
				return err
			}
		}

		if strings.HasSuffix(info.Name(), ".js") {
			fmt.Println("Processing file:", path)
			relPath, _ := filepath.Rel(srcDir, path)
			destPath := filepath.Join(destDir, strings.TrimSuffix(relPath, ".js")+".ts")
			return processFile(path, destPath)
		}

		return nil
	})
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "<source_directory> <destination_directory>")
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current dir")
	}

	srcDir := filepath.Join(cwd, os.Args[1])
	destDir := filepath.Join(cwd, os.Args[2])

	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating directory: %s", err)
	}

	err = processDirectory(srcDir, destDir)
	if err != nil {
		log.Fatalf("Error stripping types: %s", err)
	}
}
