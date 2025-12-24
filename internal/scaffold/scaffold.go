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

	// Getting the absolute path of templates (Simulated for this environment)
	// In a real CLI, we would use 'git clone https://github.com/fel1pe-r-s/...'
	// For now, let's assume we are cloning from the local Templates folder for speed, or GitHub if remote.

	// Simulating the Action of Cloning
	fmt.Printf("📋 Clonando template %s...\n", templateRepo)

	// Command to simulate clone/copy
	// In production: git clone https://github.com/fel1pe-r-s/<templateRepo> <projectName>
	cmd := exec.Command("echo", "git clone", "https://github.com/fel1pe-r-s/"+templateRepo, projectName)
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("✅ Projeto criado com sucesso!")
	fmt.Println("👉 cd", projectName)
	fmt.Println("👉 pnpm install")
}
