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
			fmt.Println("Operação cancelada.")
			return
		}

		config := scaffold.ProjectConfig{
			Name:       model.ProjectName,
			Runtime:    model.Runtime,
			Framework:  model.Framework,
			ORM:        model.ORM,
			Database:   model.Database,
			ApiPattern: model.ApiPattern,
			UseDocker:  model.UseDocker,
			UseTurbo:   model.UseTurbo,
		}

		scaffold.CreateProject(model.ProjectType, config)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
