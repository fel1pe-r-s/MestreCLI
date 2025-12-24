package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mestre",
	Short: "Mestre CLI - A Fábrica de Software",
	Long:  `Ferramenta de automação para gerar projetos baseados na Mestre Stack.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Aqui adicionaremos os subcomandos, ex: rootCmd.AddCommand(initCmd)
}
