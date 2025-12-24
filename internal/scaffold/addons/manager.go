package addons

import (
	"fmt"
	"os"
)

type AddonConfig struct {
	ProjectName string
	Framework   string
	ORM         string
	Database    string
	Runtime     string
}

// ApplyAddons runs all necessary addons based on config
func ApplyAddons(rootPath string, config AddonConfig) error {
	fmt.Println("\n⚙️  Configurando Addons...")

	if config.ORM == "prisma" {
		if err := SetupPrisma(rootPath, config); err != nil {
			return err
		}
	} else if config.ORM == "drizzle" {
		if err := SetupDrizzle(rootPath, config); err != nil {
			return err
		}
	}

	if err := SetupGithubActions(rootPath, config); err != nil {
		return err
	}

	return nil
}

func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
