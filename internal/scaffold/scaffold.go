package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templatesFS embed.FS

type ProjectConfig struct {
	Name      string
	Runtime   string
	UseDocker bool
	UseTurbo  bool
}

func CreateProject(projectType string, config ProjectConfig) {
	fmt.Printf("\n🚀 Gerando projeto: %s\n", config.Name)

	// 1. Determine base path in embedded FS
	var itemsPath string
	if strings.Contains(projectType, "Backend") {
		itemsPath = "templates/backend"
	} else if strings.Contains(projectType, "Universal") {
		itemsPath = "templates/universal"
	} else {
		itemsPath = "templates/monorepo"
	}

	// 2. Create Root Directory
	if err := os.MkdirAll(config.Name, 0755); err != nil {
		fmt.Printf("❌ Erro ao criar pasta: %v\n", err)
		return
	}

	// 3. Walk and Copy/Render
	err := fs.WalkDir(templatesFS, itemsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(itemsPath, path)
		targetPath := filepath.Join(config.Name, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read file
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}

		// If it's a template file (e.g. package.json), we might want to render variables
		// For now, let's just copy everything. Advanced interpolation can be added.
		if strings.HasSuffix(targetPath, "package.json") {
			// Simple string replacement for now to support Project Name
			sContent := string(content)
			sContent = strings.ReplaceAll(sContent, "mestrejs-backend-template", config.Name)
			if config.Runtime == "bun" {
				sContent = strings.ReplaceAll(sContent, "tsx watch", "bun --watch")
				sContent = strings.ReplaceAll(sContent, "npm", "bun")
			}
			return os.WriteFile(targetPath, []byte(sContent), 0644)
		}

		return os.WriteFile(targetPath, content, 0644)
	})

	if err != nil {
		fmt.Printf("❌ Erro na geração: %v\n", err)
		return
	}

	// 4. Create Optional Files (Docker)
	if config.UseDocker {
		dockerContent := `version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
`
		os.WriteFile(filepath.Join(config.Name, "compose.yaml"), []byte(dockerContent), 0644)
		fmt.Println("🐳 Docker Compose criado.")
	}

	fmt.Println("\n✅ Projeto criado com sucesso!")
	fmt.Printf("👉 cd %s\n", config.Name)
	if config.Runtime == "bun" {
		fmt.Println("👉 bun install")
	} else {
		fmt.Println("👉 pnpm install")
	}
}
