package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wc [option] [file]")
		os.Exit(1)
	}

	option := os.Args[1]
	filename := os.Args[2]
	if filename == "" {
		fmt.Println("File path is required.")
		os.Exit(1)
	}

	switch option {
	case "-w":
		count, err := countWords(filename)
		if err != nil {
			fmt.Println("Error counting words:", err)
			os.Exit(1)
		}
		fmt.Printf("Number of words: %d\n", count)
	case "-c":
		count := countBytes(filename)
		fmt.Printf("Number of bytes: %d\n", count)
	case "-l":
		count := countLines(filename)
		fmt.Printf("Number of lines: %d\n", count)
	default:
	}
}

func countBytes(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		os.Exit(1)
	}

	if fileInfo.Mode().IsRegular() {
		return int(fileInfo.Size())
	}

	reader := bufio.NewReader(file)
	buf := make([]byte, 1<<16) // 64 KB buffer
	count := 0
	for {
		n, err := reader.Read(buf)
		count += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
	}

	return count
}

func countLines(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	count := 0

	buf := make([]byte, 1<<16) // 64 KB buffer
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			count += bytes.Count(buf[:n], []byte{'\n'})
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
	}

	return count
}

func countWords(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count int
	inWord := false

	buf := make([]byte, 1<<16) // 64 KB buffer
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			for i := range n {
				if unicode.IsSpace(rune(buf[i])) {
					if inWord {
						count++
						inWord = false
					}
				} else {
					inWord = true
				}
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				if inWord {
					count++
				}
				break
			}
			return count, fmt.Errorf("error reading file: %w", err)
		}
	}

	return count, nil
}
