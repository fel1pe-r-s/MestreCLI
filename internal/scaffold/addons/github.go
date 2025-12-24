package addons

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetupGithubActions(root string, config AddonConfig) error {
	fmt.Println("   🔹 Configurando GitHub Actions (CI)...")

	workflowsDir := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		return err
	}

	// Minimal CI workflow
	ciContent := `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node
      uses: actions/setup-node@v4
      with:
        node-version: 20
`
	// Add Packet Manager specifics
	if config.Runtime == "bun" {
		ciContent += `
    - name: Setup Bun
      uses: oven-sh/setup-bun@v1
      with:
        bun-version: latest

    - name: Install dependencies
      run: bun install

    - name: Run Tests
      run: bun test
`
	} else {
		ciContent += `
    - name: Setup PNPM
      uses: pnpm/action-setup@v2
      with:
        version: 8

    - name: Install dependencies
      run: pnpm install

    - name: Build
      run: pnpm build

    - name: Test
      run: pnpm test
`
	}

	return writeFile(filepath.Join(workflowsDir, "ci.yml"), ciContent)
}
