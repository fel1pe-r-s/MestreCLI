package cmd

import (
	"fmt"
	"mestre/internal/scaffold"
	"mestre/internal/ui"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicia um novo projeto interativo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔮 Iniciando o Wizard do Mestre...")
		model, err := ui.StartWizard()
		if err != nil {
			fmt.Println("Erro:", err)
			return
		}

		if model.ProjectType == "" {
			fmt.Println("Nenhum projeto selecionado.")
			return
		}

		// Hardcoded project name for demo purposes, in real app we ask via Input in BubbleTea
		scaffold.CreateProject(model.ProjectType, "meu-novo-projeto")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
