package scaffold_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mestre/internal/scaffold"
)

func TestEndToEndGenerator(t *testing.T) {
	// config simulating: NestJS + Drizzle + Postgres + Bun
	config := scaffold.ProjectConfig{
		Name:      "test-e2e-project",
		Runtime:   "bun",
		Framework: "nestjs",
		ORM:       "drizzle",
		Database:  "postgres",
		UseDocker: true,
	}

	// Clean up before
	os.RemoveAll(config.Name)
	defer os.RemoveAll(config.Name)

	// Execute Generator
	scaffold.CreateProject("Backend API", config)

	// --- VERIFICATIONS ---

	// 1. Check Directory Structure (NestJS)
	if _, err := os.Stat(filepath.Join(config.Name, "src", "main.ts")); os.IsNotExist(err) {
		t.Errorf("NestJS main.ts not found. Template selection failed.")
	}

	// 2. Check Package.json Injection (Patcher)
	pkgBytes, _ := os.ReadFile(filepath.Join(config.Name, "package.json"))
	var pkg map[string]interface{}
	json.Unmarshal(pkgBytes, &pkg)

	deps := pkg["dependencies"].(map[string]interface{})
	devDeps := pkg["devDependencies"].(map[string]interface{})

	if _, ok := deps["@nestjs/core"]; !ok {
		t.Errorf("Patcher failed: @nestjs/core missing")
	}
	if _, ok := deps["drizzle-orm"]; !ok {
		t.Errorf("Patcher failed: drizzle-orm missing")
	}
	if _, ok := devDeps["drizzle-kit"]; !ok {
		t.Errorf("Patcher failed: drizzle-kit missing in devDeps")
	}
	// Bun check
	scripts := pkg["scripts"].(map[string]interface{})
	if !strings.Contains(scripts["start:dev"].(string), "bun --watch") {
		t.Errorf("Runtime adjust failed: bun --watch not found in scripts")
	}

	// 3. Check Addon: Drizzle Config
	if _, err := os.Stat(filepath.Join(config.Name, "drizzle.config.ts")); os.IsNotExist(err) {
		t.Errorf("Addon failed: drizzle.config.ts missing")
	}

	// 4. Check Addon: Github CI
	ciBytes, _ := os.ReadFile(filepath.Join(config.Name, ".github", "workflows", "ci.yml"))
	if !strings.Contains(string(ciBytes), "oven-sh/setup-bun") {
		t.Errorf("Addon failed: CI does not use Bun action")
	}
}
