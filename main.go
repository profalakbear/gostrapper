package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createFolderStructure(basePath string, entries []string) error {
	for _, entry := range entries {
		path := filepath.Join(basePath, entry)
		if path == basePath {
			continue // Skip creating the base path itself
		}

		if strings.HasSuffix(entry, "/") {
			// Create directory
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
			fmt.Printf("Created folder: %s\n", path)
		} else {
			// Create file
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
			if _, err := os.Create(path); err != nil {
				return err
			}
			fmt.Printf("Created file: %s\n", path)
		}
	}
	return nil
}

func initGoMod(basePath string) error {
	// Change working directory to base path
	if err := os.Chdir(basePath); err != nil {
		return err
	}

	// Run 'go mod init' command
	cmd := exec.Command("go", "mod", "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running 'go mod init': %v\nOutput: %s", err, output)
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
		fmt.Println("Error: You must provide both the -p and -s flags")
		return
	}

	structure, err := readStructureFromFile(*s)
	if err != nil {
		fmt.Println("Error reading structure file:", err)
		return
	}

	if err := createFolderStructure(*p, structure); err != nil {
		fmt.Println("Error creating folder structure:", err)
		return
	}

	// Initialize go.mod and go.sum
	if err := initGoMod(*p); err != nil {
		fmt.Println("Error initializing go.mod:", err)
	}
}
