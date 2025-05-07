package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Get module name from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your module name (e.g., github.com/username/project): ")
	moduleName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}
	moduleName = strings.TrimSpace(moduleName)

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Create new directory with module name
	dirName := filepath.Base(moduleName)
	if err := os.MkdirAll(dirName, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Copy necessary files and directories
	filesToCopy := []string{
		"internal",
		"utils",
		"api",
		"main.go",
		"Makefile",
		"env_example",
	}

	for _, file := range filesToCopy {
		src := filepath.Join(currentDir, file)
		dst := filepath.Join(dirName, file)

		if err := copyPath(src, dst); err != nil {
			fmt.Printf("Error copying %s: %v\n", file, err)
			os.Exit(1)
		}
	}

	// Change to the new directory
	if err := os.Chdir(dirName); err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize new module
	cmd := exec.Command("go", "mod", "init", moduleName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error initializing module: %v\n", err)
		os.Exit(1)
	}

	// Update imports in all Go files
	if err := updateImports(".", "github.com/kkwitslab/go-boilerplate", moduleName); err != nil {
		fmt.Printf("Error updating imports: %v\n", err)
		os.Exit(1)
	}

	// Copy dependencies from original go.mod
	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error getting dependencies: %v\n", err)
		os.Exit(1)
	}

	// Create .env file from env_example
	envContent, err := os.ReadFile("env_example")
	if err != nil {
		fmt.Printf("Error reading env_example: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(".env", envContent, 0644); err != nil {
		fmt.Printf("Error creating .env file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nProject setup complete! Your new project is in the '%s' directory.\n", dirName)
	fmt.Println("Next steps:")
	fmt.Println("1. cd", dirName)
	fmt.Println("2. Review and update the .env file")
	fmt.Println("3. Run 'go mod tidy' to ensure all dependencies are properly set up")
	fmt.Println("4. Run 'go run main.go' to start the server")
}

func updateImports(dir, oldModule, newModule string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace old module with new module in imports
		newContent := strings.ReplaceAll(string(content), oldModule, newModule)

		// Write updated content back to file
		return os.WriteFile(path, []byte(newContent), info.Mode())
	})
}

func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, 0644)
}
