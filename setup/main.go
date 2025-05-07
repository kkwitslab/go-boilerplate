package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

func main() {
	// Get module name from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s🚀 %sEnter your module name (e.g., github.com/username/project): %s", colorCyan, colorYellow, colorReset)
	moduleName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("%s❌ Error reading input: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	moduleName = strings.TrimSpace(moduleName)

	// Get directory name from user
	fmt.Printf("%s📁 %sEnter directory name for the project (press Enter to use module name): %s", colorCyan, colorYellow, colorReset)
	dirName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("%s❌ Error reading input: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	dirName = strings.TrimSpace(dirName)
	if dirName == "" {
		dirName = filepath.Base(moduleName)
	}

	// Get the module's source directory
	_, currentFile, _, _ := runtime.Caller(0)
	sourceDir := filepath.Dir(filepath.Dir(currentFile))

	// Create new directory with specified name
	fmt.Printf("%s📂 %sCreating project directory...%s\n", colorBlue, colorYellow, colorReset)
	if err := os.MkdirAll(dirName, 0755); err != nil {
		fmt.Printf("%s❌ Error creating directory: %v%s\n", colorRed, err, colorReset)
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

	fmt.Printf("%s📋 %sCopying project files...%s\n", colorBlue, colorYellow, colorReset)
	for _, file := range filesToCopy {
		src := filepath.Join(sourceDir, file)
		dst := filepath.Join(dirName, file)

		if err := copyPath(src, dst); err != nil {
			fmt.Printf("%s❌ Error copying %s: %v%s\n", colorRed, file, err, colorReset)
			os.Exit(1)
		}
		fmt.Printf("%s  ✓ %s%s\n", colorGreen, file, colorReset)
	}

	// Change to the new directory
	if err := os.Chdir(dirName); err != nil {
		fmt.Printf("%s❌ Error changing directory: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Initialize new module
	fmt.Printf("%s🔧 %sInitializing Go module...%s\n", colorBlue, colorYellow, colorReset)
	cmd := exec.Command("go", "mod", "init", moduleName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%s❌ Error initializing module: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Update imports in all Go files
	fmt.Printf("%s🔄 %sUpdating imports...%s\n", colorBlue, colorYellow, colorReset)
	if err := updateImports(".", "github.com/kkwitslab/go-boilerplate", moduleName); err != nil {
		fmt.Printf("%s❌ Error updating imports: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Copy dependencies from original go.mod
	fmt.Printf("%s📦 %sSetting up dependencies...%s\n", colorBlue, colorYellow, colorReset)
	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		fmt.Printf("%s❌ Error getting dependencies: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Create .env file from env_example
	fmt.Printf("%s⚙️ %sCreating environment file...%s\n", colorBlue, colorYellow, colorReset)
	envContent, err := os.ReadFile("env_example")
	if err != nil {
		fmt.Printf("%s❌ Error reading env_example: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	if err := os.WriteFile(".env", envContent, 0644); err != nil {
		fmt.Printf("%s❌ Error creating .env file: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("\n%s✨ %sProject setup complete! Your new project is in the '%s' directory.%s\n", colorGreen, colorYellow, dirName, colorReset)
	fmt.Printf("\n%s📝 %sNext steps:%s\n", colorPurple, colorYellow, colorReset)
	fmt.Printf("%s  1. %scd %s%s\n", colorCyan, colorYellow, dirName, colorReset)
	fmt.Printf("%s  2. %sReview and update the .env file%s\n", colorCyan, colorYellow, colorReset)
	fmt.Printf("%s  3. %sRun 'go mod tidy' to ensure all dependencies are properly set up%s\n", colorCyan, colorYellow, colorReset)
	fmt.Printf("%s  4. %sRun 'go run main.go' to start the server%s\n", colorCyan, colorYellow, colorReset)
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
