package scaffold_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mestre/internal/scaffold"
)

func TestMatrixGenerator(t *testing.T) {
	tests := []struct {
		name   string
		pType  string // Project Type
		config scaffold.ProjectConfig
		checks func(t *testing.T, dir string)
	}{
		{
			name:  "NestJS_Drizzle_Postgres",
			pType: "Backend API",
			config: scaffold.ProjectConfig{
				Name: "test-nest-drizzle", Runtime: "node", Framework: "nestjs", ORM: "drizzle", Database: "postgres", UseDocker: true,
			},
			checks: func(t *testing.T, dir string) {
				assertFileExists(t, dir, "src/main.ts")       // Nest
				assertFileExists(t, dir, "drizzle.config.ts") // Drizzle
				assertFileExists(t, dir, "compose.yaml")      // Docker
				assertPackageJsonDeps(t, dir, []string{"@nestjs/core", "drizzle-orm"}, []string{"drizzle-kit"})
			},
		},
		{
			name:  "Fastify_Prisma_Mongo_Bun",
			pType: "Backend API",
			config: scaffold.ProjectConfig{
				Name: "test-fastify-prisma", Runtime: "bun", Framework: "fastify", ORM: "prisma", Database: "mongo", UseDocker: false,
			},
			checks: func(t *testing.T, dir string) {
				assertFileExists(t, dir, "src/infra/http/server.ts") // Fastify/CleanArch
				assertFileExists(t, dir, "prisma/schema.prisma")     // Prisma
				assertFileNotExists(t, dir, "compose.yaml")          // No Docker
				assertPackageJsonDeps(t, dir, []string{"fastify", "@prisma/client"}, []string{"prisma"})

				// Verify Mongo Provider
				content, _ := os.ReadFile(filepath.Join(dir, "prisma/schema.prisma"))
				if !strings.Contains(string(content), `provider = "mongodb"`) {
					t.Errorf("Prisma Schema should have mongodb provider")
				}
			},
		},
		{
			name:  "NestJS_NoORM_SQLite",
			pType: "Backend API",
			config: scaffold.ProjectConfig{
				Name: "test-nest-clean", Runtime: "node", Framework: "nestjs", ORM: "none", Database: "sqlite", UseDocker: true,
			},
			checks: func(t *testing.T, dir string) {
				assertFileExists(t, dir, "src/main.ts")
				assertFileNotExists(t, dir, "prisma/schema.prisma")
				assertFileNotExists(t, dir, "drizzle.config.ts")
				assertPackageJsonDeps(t, dir, []string{"@nestjs/core"}, nil)
			},
		},
		{
			name:  "Universal_Turbo",
			pType: "Universal App (Web/Mobile/Desktop)",
			config: scaffold.ProjectConfig{
				Name: "test-universal", Runtime: "node", UseDocker: true, UseTurbo: true,
			},
			checks: func(t *testing.T, dir string) {
				assertFileExists(t, dir, "turbo.json")
				// apps/web/package.json might be missing if template is empty, handled by fix
				assertFileExists(t, dir, "apps/web/package.json")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			os.RemoveAll(tt.config.Name)
			defer os.RemoveAll(tt.config.Name)

			// Run
			scaffold.CreateProject(tt.pType, tt.config) // Using "Backend" as default type for most, adjusted logic inside checks implies structure

			// Verify
			tt.checks(t, tt.config.Name)
		})
	}
}

// Helpers
func assertFileExists(t *testing.T, base, path string) {
	if _, err := os.Stat(filepath.Join(base, path)); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it does not", path)
	}
}

func assertFileNotExists(t *testing.T, base, path string) {
	if _, err := os.Stat(filepath.Join(base, path)); err == nil {
		t.Errorf("Expected file %s to NOT exist, but it does", path)
	}
}

func assertPackageJsonDeps(t *testing.T, base string, deps []string, devDeps []string) {
	pkgBytes, err := os.ReadFile(filepath.Join(base, "package.json"))
	if err != nil {
		t.Fatalf("Could not read package.json: %v", err)
	}

	var pkg map[string]interface{}
	json.Unmarshal(pkgBytes, &pkg)

	pDeps, _ := pkg["dependencies"].(map[string]interface{})
	pDevDeps, _ := pkg["devDependencies"].(map[string]interface{})

	for _, d := range deps {
		if _, ok := pDeps[d]; !ok {
			t.Errorf("Missing dependency: %s", d)
		}
	}
	for _, d := range devDeps {
		if _, ok := pDevDeps[d]; !ok {
			t.Errorf("Missing devDependency: %s", d)
		}
	}
}
