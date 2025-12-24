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

	// Executing Real Git Clone
	repoURL := "https://github.com/fel1pe-r-s/" + templateRepo + ".git"
	fmt.Printf("📋 Clonando de: %s\n", repoURL)

	cmd := exec.Command("git", "clone", repoURL, projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr // Show git progress
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
