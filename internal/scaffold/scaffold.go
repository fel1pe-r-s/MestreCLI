package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"mestre/internal/scaffold/addons"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templatesFS embed.FS

type ProjectConfig struct {
	Name      string
	Runtime   string
	Framework string
	ORM       string
	Database  string
	UseDocker bool
	UseTurbo  bool
}

func CreateProject(projectType string, config ProjectConfig) {
	fmt.Printf("\n🚀 Gerando projeto: %s\n", config.Name)

	// 1. Determine base path in embedded FS
	var itemsPath string
	if strings.Contains(projectType, "Backend") {
		if config.Framework == "nestjs" {
			itemsPath = "templates/backend-nest"
		} else {
			itemsPath = "templates/backend"
		}
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

		// 🟢 Docker Skip Logic
		if !config.UseDocker {
			fname := filepath.Base(path)
			if fname == "Dockerfile" || fname == "compose.yaml" || fname == "docker-compose.yml" || fname == "compose.yml" {
				return nil // Skip this file
			}
		}

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read file
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}

		// If it's a template file (e.g. package.json), we might want to render variables
		if strings.HasSuffix(targetPath, "package.json") {
			// New Smart Patcher
			newContent, err := PatchPackageJSON(content, config)
			if err != nil {
				fmt.Printf("⚠️ Erro ao patchear package.json: %v\n", err)
				return os.WriteFile(targetPath, content, 0644) // Fallback to original
			}
			return os.WriteFile(targetPath, newContent, 0644)
		}

		return os.WriteFile(targetPath, content, 0644)
	})

	if err != nil {
		fmt.Printf("❌ Erro na geração: %v\n", err)
		return
	}

	// 5. Apply Addons (Prisma, Github, etc)
	addonConfig := addons.AddonConfig{
		ProjectName: config.Name,
		Framework:   config.Framework,
		ORM:         config.ORM,
		Database:    config.Database,
		Runtime:     config.Runtime,
	}
	if err := addons.ApplyAddons(config.Name, addonConfig); err != nil {
		fmt.Printf("⚠️ Erro nos addons: %v\n", err)
	}

	fmt.Println("\n✅ Projeto criado com sucesso!")
	fmt.Println("\n🏁 Próximos Passos:")
	fmt.Println("------------------------------------------------")

	// 1. Enter
	fmt.Printf("1️⃣  Entrar na pasta:\n    cd %s\n\n", config.Name)

	// 2. Install
	cmd := "pnpm"
	if config.Runtime == "bun" {
		cmd = "bun"
	}
	fmt.Printf("2️⃣  Instalar dependências:\n    %s install\n\n", cmd)

	// 3. Database
	if config.ORM == "prisma" {
		run := "npx"
		if config.Runtime == "bun" {
			run = "bunx"
		}
		fmt.Printf("3️⃣  Configurar Banco (Prisma):\n    %s prisma generate\n    %s prisma migrate dev\n\n", run, run)
	} else if config.ORM == "drizzle" {
		run := "npx"
		if config.Runtime == "bun" {
			run = "bunx"
		}
		cmdArg := "generate"
		if config.Database == "postgres" {
			cmdArg += ":pg"
		} else if config.Database == "mysql" {
			cmdArg += ":mysql"
		} else {
			cmdArg += ":sqlite"
		}
		fmt.Printf("3️⃣  Configurar Banco (Drizzle):\n    %s drizzle-kit %s\n\n", run, cmdArg)
	}

	// 4. Docker
	if config.UseDocker {
		fmt.Printf("🐳 Configurar Container:\n    docker compose up -d\n\n")
	}

	// 5. Run
	runCmd := "dev"
	if config.Framework == "nestjs" {
		runCmd = "start:dev"
	}
	fmt.Printf("🚀 Rodar o projeto:\n    %s run %s\n", cmd, runCmd)

	fmt.Println("------------------------------------------------")
}
