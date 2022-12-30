package utils

import (
	"bufio"
	"encoding/json"
	"os"
)

func ReadUrls(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	return urls, nil
}

func WriteJson(data interface{}, filename string, indent bool) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	e := json.NewEncoder(outFile)
	if indent {
		e.SetIndent("", "\t")
	}
	return e.Encode(data)
}

func WritePlain(data []string, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	for _, str := range data {
		outFile.WriteString(str + "\n")
	}
	return nil
}
