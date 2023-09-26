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

func initGoMod(basePath, modulePath string) error {
	// Change working directory to base path
	if err := os.Chdir(basePath); err != nil {
		return err
	}

	var cmd *exec.Cmd
	if modulePath != "" {
		// Run 'go mod init' command with the provided module path
		cmd = exec.Command("go", "mod", "init", modulePath)
	} else {
		// Run 'go mod init' command without specifying module path
		cmd = exec.Command("go", "mod", "init")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running 'go mod init': %v\nOutput: %s", err, output)
	}
	return nil
}

func writePackageDeclaration(filePath string, packageName string) error {
	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Open the file for writing (truncating existing content)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write package declaration and the original content
	_, err = file.WriteString("package " + packageName + "\n")
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func cleanupOnError(basePath string) {
	fmt.Println("Cleaning up...")
	if err := os.RemoveAll(basePath); err != nil {
		fmt.Println("Error cleaning up:", err)
	}
}

func readStructureFromFile(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

func getStructurePathFromFlagOrEnv() (string, error) {
	// Check if the -s flag is provided
	sFlag := flag.String("s", "", "File containing the folder structure")
	flag.Parse()
	if *sFlag != "" {
		return *sFlag, nil
	}

	// If -s flag is not provided, check the GOSTRAP_STRUCT environment variable
	structureEnv := os.Getenv("GOSTRAP_STRUCT")
	if structureEnv != "" {
		return structureEnv, nil
	}

	return "", fmt.Errorf("structure file not provided (-s flag or GOSTRAP_STRUCT environment variable)")
}

func main() {
	p := flag.String("p", "", "Root path for creating the folder structure")
	m := flag.String("m", "", "Optional go mod module path")
	flag.Parse()

	if *p == "" {
		fmt.Println("Error: You must provide -p flag")
		return
	}

	absFilePath, err := filepath.Abs(*p)
	if err != nil {
		return
	}

	// Get the structure file path from flag or environment variable
	structurePath, err := getStructurePathFromFlagOrEnv()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	absStructurePath, err := filepath.Abs(structurePath)
	if err != nil {
		return
	}

	fmt.Println("Absolute target path", absFilePath)
	fmt.Println("Absolute structure file path", absStructurePath)

	structure, err := readStructureFromFile(absStructurePath)
	if err != nil {
		fmt.Println("Error reading structure file:", err)
		return
	}

	if err := createFolderStructure(absFilePath, structure); err != nil {
		fmt.Println("Error creating folder structure:", err)
		cleanupOnError(absFilePath)
		return
	}

	// Initialize go.mod and go.sum
	if err := initGoMod(absFilePath, *m); err != nil {
		fmt.Println("Error initializing go.mod:", err)
		cleanupOnError(absFilePath)
		return
	}

	// Write package declarations for .go files
	err = filepath.Walk(absFilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			relPath, err := filepath.Rel(absFilePath, path)
			if err != nil {
				return err
			}

			packageName := ""
			if filepath.Base(path) == "main.go" {
				packageName = "main"
			} else {
				dir := filepath.Dir(relPath)
				packageName = filepath.Base(dir)
			}

			if err := writePackageDeclaration(path, packageName); err != nil {
				fmt.Printf("Error writing package declaration for %s: %v\n", path, err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error writing package declarations:", err)
		cleanupOnError(absFilePath)
		return
	}
}
