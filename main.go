package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	version = "1.0.0"
	janusDir = ".janus"
	templatesDir = "templates"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	arg := os.Args[1]

	switch arg {
	case "-h", "--help":
		showHelp()
	case "-v", "--version":
		fmt.Printf("Janus v%s\n", version)
	case "-o":
		if len(os.Args) > 2 && os.Args[2] == "templates" {
			listOfflineTemplates()
		} else {
			fmt.Println("Usage: janus -o templates")
		}
	default:
		domain := arg
		subdomain := ""
		if len(os.Args) > 2 {
			subdomain = os.Args[2]
		}
		scaffoldProject(domain, subdomain)
	}
}

func showHelp() {
	fmt.Printf(`Janus v%s - Cross-domain Project Scaffolder

Usage:
  janus <domain> [subdomain]     Scaffold a new project
  janus -o templates             List available offline templates
  janus -h, --help               Show this help
  janus -v, --version            Show version

Examples:
  janus nextjs                   Create nextjs-app from template
  janus python flask             Create python-flask-app from template
  janus -o templates             List all local templates

Templates are stored at: %s
`, version, getTemplatesPath())
}

func getJanusPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, janusDir)
}

func getTemplatesPath() string {
	return filepath.Join(getJanusPath(), templatesDir)
}

func listOfflineTemplates() {
	templatesPath := getTemplatesPath()

	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		fmt.Printf("No templates directory found at: %s\n", templatesPath)
		fmt.Println("Create the directory and add your templates to get started.")
		return
	}

	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		fmt.Printf("Error reading templates directory: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No templates found in:", templatesPath)
		return
	}

	fmt.Printf("Available offline templates (%s):\n", templatesPath)
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  • %s\n", entry.Name())
		}
	}
}

func scaffoldProject(domain, subdomain string) {
	templateName := domain
	if subdomain != "" {
		templateName = domain + "-" + subdomain
	}

templatePath := filepath.Join(getTemplatesPath(), domain, subdomain)

	// Check if template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		fmt.Printf("Template '%s' not found at: %s\n", templateName, templatePath)
		fmt.Println("Available templates:")
		listOfflineTemplates()
		return
	}

	// Ask user for project name
	var projectName string
	fmt.Printf("Enter your project name (default: %s-app): ", templateName)

	// Read input with better handling
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(input)

	// Use default if no input provided
	if projectName == "" {
		projectName = templateName + "-app"
	}

	// Sanitize project name (remove invalid characters)
	projectName = sanitizeProjectName(projectName)

	// Check if target directory already exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		fmt.Printf("Directory '%s' already exists. Aborting to avoid overwriting.\n", projectName)
		return
	}

	// Create target directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		fmt.Printf("Error creating target directory: %v\n", err)
		return
	}

	// Copy template to target directory
	fmt.Printf("Scaffolding project '%s' from template '%s'...\n", projectName, templateName)
	if err := copyDir(templatePath, projectName); err != nil {
		fmt.Printf("Error copying template: %v\n", err)
		// Clean up failed directory
		os.RemoveAll(projectName)
		return
	}

	fmt.Printf("✅ Project scaffolded successfully in './%s'\n", projectName)

	// Install dependencies if package manager files exist
	installDeps(projectName)

	// Auto git init
	initGit(projectName)

	fmt.Printf("\n🎉 Project ready! Next steps:\n")
	fmt.Printf("  cd %s\n", projectName)

	// Give domain-specific next steps
	if strings.Contains(templateName, "nextjs") || strings.Contains(templateName, "react") {
		fmt.Println("  npm run dev")
	} else if strings.Contains(templateName, "python") || strings.Contains(templateName, "flask") || strings.Contains(templateName, "django") {
		fmt.Println("  python main.py  # or your main file")
	} else if strings.Contains(templateName, "node") || strings.Contains(templateName, "express") {
		fmt.Println("  npm start")
	} else {
		fmt.Println("  # Start coding!")
	}
}

func copyDir(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
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
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Copy file permissions
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}

	return nil
}

func installDeps(projectDir string) {
	// Check for package.json (npm)
	packageJsonPath := filepath.Join(projectDir, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		fmt.Println("📦 Found package.json, installing npm dependencies...")
		if err := runCommand(projectDir, "npm", "install"); err != nil {
			fmt.Printf("⚠️  Warning: npm install failed: %v\n", err)
		} else {
			fmt.Println("✅ npm dependencies installed successfully")
		}
		return
	}

	// Check for requirements.txt (pip)
	requirementsPath := filepath.Join(projectDir, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		fmt.Println("📦 Found requirements.txt, installing pip dependencies...")
		if err := runCommand(projectDir, "pip", "install", "-r", "requirements.txt"); err != nil {
			// Try pip3 if pip fails
			if err := runCommand(projectDir, "pip3", "install", "-r", "requirements.txt"); err != nil {
				fmt.Printf("⚠️  Warning: pip install failed: %v\n", err)
			} else {
				fmt.Println("✅ pip dependencies installed successfully")
			}
		} else {
			fmt.Println("✅ pip dependencies installed successfully")
		}
		return
	}

	// Check for go.mod (Go modules)
	goModPath := filepath.Join(projectDir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		fmt.Println("📦 Found go.mod, installing Go dependencies...")
		if err := runCommand(projectDir, "go", "mod", "tidy"); err != nil {
			fmt.Printf("⚠️  Warning: go mod tidy failed: %v\n", err)
		} else {
			fmt.Println("✅ Go dependencies installed successfully")
		}
		return
	}

	// Check for Cargo.toml (Rust)
	cargoPath := filepath.Join(projectDir, "Cargo.toml")
	if _, err := os.Stat(cargoPath); err == nil {
		fmt.Println("📦 Found Cargo.toml, installing Rust dependencies...")
		if err := runCommand(projectDir, "cargo", "fetch"); err != nil {
			fmt.Printf("⚠️  Warning: cargo fetch failed: %v\n", err)
		} else {
			fmt.Println("✅ Rust dependencies fetched successfully")
		}
		return
	}

	fmt.Println("📦 No recognized dependency files found, skipping dependency installation")
}

func initGit(projectDir string) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("⚠️  Git not found, skipping git initialization")
		return
	}

	// Check if already a git repository
	gitDir := filepath.Join(projectDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		fmt.Println("📁 Directory is already a git repository")
		return
	}

	fmt.Println("🔧 Initializing git repository...")
	if err := runCommand(projectDir, "git", "init"); err != nil {
		fmt.Printf("⚠️  Warning: git init failed: %v\n", err)
		return
	}

	fmt.Println("✅ Git repository initialized")

	// Add initial commit
	fmt.Println("📝 Creating initial commit...")
	if err := runCommand(projectDir, "git", "add", "."); err != nil {
		fmt.Printf("⚠️  Warning: git add failed: %v\n", err)
		return
	}

	if err := runCommand(projectDir, "git", "commit", "-m", "Initial commit from Janus"); err != nil {
		fmt.Printf("⚠️  Warning: git commit failed: %v\n", err)
		return
	}

	fmt.Println("✅ Initial commit created")
}

func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	// On Windows, we need to handle the PATH differently
	if runtime.GOOS == "windows" {
		cmd.Env = os.Environ()
	}

	return cmd.Run()
}

func sanitizeProjectName(name string) string {
	// Remove invalid characters for directory names
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := reg.ReplaceAllString(name, "")

	// Replace spaces with hyphens
	sanitized = strings.ReplaceAll(sanitized, " ", "-")

	// Convert to lowercase
	sanitized = strings.ToLower(sanitized)

	// Remove leading/trailing hyphens and dots
	sanitized = strings.Trim(sanitized, "-.")

	// Ensure it's not empty
	if sanitized == "" {
		sanitized = "my-project"
	}

	return sanitized
}
