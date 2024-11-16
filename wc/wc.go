package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
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
	case "-m":
		count, err := countRunes(filename)
		if err != nil {
			fmt.Println("Error counting runes:", err)
			os.Exit(1)
		}
		fmt.Printf("Number of runes: %d\n", count)
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

func countRunes(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	const bufferSize = 1 << 16 // 64 KB buffer
	buf := make([]byte, bufferSize)
	reader := bufio.NewReader(file)
	count := 0

	// To handle runes that span across buffer boundaries.
	var leftover []byte
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			data := buf[:n]
			if len(leftover) > 0 {
				// Prepend leftover bytes to the current buffer.
				data = append(leftover, data...)
				leftover = nil
			}

			i := 0
			for i < len(data) {
				r, size := utf8.DecodeRune(data[i:])
				if r == utf8.RuneError && size == 1 {
					// Incomplete rune at the end of the buffer.
					break
				}
				count++
				i += size
			}

			if i < len(data) {
				// Store leftover bytes for the next read.
				leftover = append(leftover, data[i:]...)
			}
		}

		if err != nil {
			if !errors.Is(err, io.EOF) {
				return count, fmt.Errorf("error reading file: %w", err)
			}

			// If there are leftover bytes, check if they form a valid rune.
			if len(leftover) > 0 {
				if _, size := utf8.DecodeRune(leftover); size > 0 {
					count++
				}
			}
			break
		}
	}

	return count, nil
}
