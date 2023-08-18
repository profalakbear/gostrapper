package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createFolderStructure(basePath string, folders []string) error {
	for _, folder := range folders {
		path := filepath.Join(basePath, folder)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
		fmt.Printf("Created folder: %s\n", path)
	}
	return nil
}

func readStructureFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

func main() {
	p := flag.String("p", "", "Root path for creating the folder structure")
	s := flag.String("s", "", "File containing the folder structure")
	flag.Parse()

	if *p == "" || *s == "" {
		fmt.Println("Error: You must provide both the -p and -structure flags")
		return
	}

	structure, err := readStructureFromFile(*s)
	if err != nil {
		fmt.Println("Error reading structure file:", err)
		return
	}

	if err := createFolderStructure(*p, structure); err != nil {
		fmt.Println("Error:", err)
	}
}
