// de-duplication project main.go
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func trimChar(r rune) bool {
	return r == ' ' || r == '*'
}

func writeFile(fileName string, m map[string]bool) error {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	for str, _ := range m {
		buf.WriteString(str)
		buf.WriteString("\r\n")
	}

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(buf.String())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func readFile(fileName string, m map[string]bool) error {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		str := strings.TrimFunc(string(line), trimChar)
		if _, ok := m[str]; !ok {
			m[str] = true
		} else {
			fmt.Println("duplicate value : ", str)
		}
	}
	return nil
}

func main() {
	var input, outputFile, force string
	flag.StringVar(&input, "i", "", "input files")
	flag.StringVar(&outputFile, "o", "", "output file")
	flag.StringVar(&force, "f", "no", "force")
	flag.Parse()

	if input == "" {
		fmt.Println("input files is empty")
		return
	}
	if outputFile == "" {
		fmt.Println("output file is empty")
		return
	}

	if !strings.EqualFold(force, "y") && !strings.EqualFold(force, "yes") {
		if isFileExist(outputFile) {
			fmt.Printf("output file \"%s\" has existed, do you want to override it ? (y/n) : ", outputFile)
			var command string
			fmt.Scanln(&command)
			if !strings.EqualFold(command, "y") {
				fmt.Println("exit!")
				return
			}
		}
	}

	m := make(map[string]bool)

	inputFiles := strings.Split(input, "+")
	fmt.Printf("input file count : %d\n", len(inputFiles))
	for _, v := range inputFiles {
		fmt.Printf("-| %s : ", v)
		if !isFileExist(v) {
			fmt.Println("file not exist")
			continue
		} else {
			fmt.Println("")
		}
		readFile(v, m)
	}

	writeFile(outputFile, m)
	fmt.Println("output file :", outputFile)
}
