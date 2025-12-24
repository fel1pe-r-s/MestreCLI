package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CreateProject(projectType string, projectName string) {
	fmt.Printf("🚀 Iniciando criação do projeto: %s\n", projectName)
	fmt.Printf("📦 Padrão escolhido: %s\n", projectType)

	var templateRepo string

	if strings.Contains(projectType, "Backend") {
		templateRepo = "MestreJS_Backend"
	} else if strings.Contains(projectType, "Universal") {
		templateRepo = "MestreJS_Universal"
	} else {
		templateRepo = "MestreJS_Monorepo"
	}

	// Uses 'gh' to respect user's auth (SSH or HTTPS) automatically
	repoID := "fel1pe-r-s/" + templateRepo
	fmt.Printf("📋 Clonando template: %s\n", repoID)

	cmd := exec.Command("gh", "repo", "clone", repoID, projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("❌ Erro ao clonar: %v\n", err)
		return
	}

	fmt.Println("✅ Projeto criado com sucesso!")
	fmt.Println("👉 cd", projectName)
	fmt.Println("👉 pnpm install")
	fmt.Println("👉 mestre list (em breve)")
}
