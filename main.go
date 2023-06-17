package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var totalCount int

	directory := flag.String("d", "", "the directory to process")
	escapeDirs := flag.String("esc", "", "the escape directories to skip processing (comma-separated)")

	flag.Parse()

	if *directory == "" {
		fmt.Println("Directory argument is required.")
		return
	}

	escapeMap := make(map[string]bool)
	if *escapeDirs != "" {
		escapeList := parseEscapeDirs(*escapeDirs)
		for _, dir := range escapeList {
			escapeMap[dir] = true
		}
	}

	err := filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return nil
		}

		if !info.IsDir() {
			count, err := countRows(path)
			if err != nil {
				log.Printf("Error counting rows in file %q: %v\n", path, err)
				return nil
			}

			fmt.Printf("File: %s | Rows: %d\n", path, count)
			totalCount += count
		} else {
			dirName := info.Name()
			if escapeMap[dirName] {
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total count: %d \n", totalCount)
}

func countRows(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineCount := 0

	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		lineCount++
	}

	return lineCount, nil
}

func parseEscapeDirs(escapeDirs string) []string {
	dirList := strings.Split(escapeDirs, ",")
	dirs := make([]string, 0)

	for _, dir := range dirList {
		trimmedDir := strings.TrimSpace(dir)
		if trimmedDir != "" {
			dirs = append(dirs, trimmedDir)
		}
	}

	return dirs
}
