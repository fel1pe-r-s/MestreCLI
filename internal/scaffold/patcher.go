package scaffold

import (
	"encoding/json"
	"strings"
)

// PatchPackageJSON reads a package.json, modifies it based on ProjectConfig, and returns the new bytes.
// It uses map[string]interface{} to preserve all fields in the JSON.
func PatchPackageJSON(content []byte, config ProjectConfig) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	// 1. Update Project Name
	data["name"] = config.Name

	// 2. Ensure dependencies/devDependencies maps exist
	deps, _ := data["dependencies"].(map[string]interface{})
	if deps == nil {
		deps = make(map[string]interface{})
	}

	devDeps, _ := data["devDependencies"].(map[string]interface{})
	if devDeps == nil {
		devDeps = make(map[string]interface{})
	}

	// Helpers to add dependencies
	addDep := func(name, version string) { deps[name] = version }
	addDev := func(name, version string) { devDeps[name] = version }

	// 3. Inject Dependencies based on Config

	// ORM Strategy
	if config.ORM == "prisma" {
		addDep("@prisma/client", "^5.7.0")
		addDev("prisma", "^5.7.0")
	} else if config.ORM == "drizzle" {
		addDep("drizzle-orm", "^0.29.0")
		addDep("postgres", "^3.4.0") // Drizzle needs a driver
		addDev("drizzle-kit", "^0.20.0")
	}

	// Database Drivers (if generic or needed)
	if config.Database == "postgres" && config.ORM != "drizzle" {
		// If Drizzle is used, we added 'postgres' above.
		// If Prisma is used, it has its own engine but 'pg' is sometimes useful.
		// Let's add 'pg' for standard usage if not using Drizzle's driver.
		if config.ORM == "none" || config.Framework == "fastify" {
			addDep("pg", "^8.11.0")
			addDev("@types/pg", "^8.10.0")
		}
	} else if config.Database == "mongo" {
		addDep("mongodb", "^6.3.0")
	} else if config.Database == "sqlite" {
		addDep("sqlite3", "^5.1.0")
	}

	// Frameworks (Advanced Injection)
	if config.Framework == "nestjs" {
		// Core NestJS deps
		addDep("@nestjs/core", "^10.0.0")
		addDep("@nestjs/common", "^10.0.0")
		addDep("@nestjs/platform-express", "^10.0.0") // Default platform
		addDep("reflect-metadata", "^0.1.13")
		addDep("rxjs", "^7.8.0")

		// NestJS often needs these too
		addDev("@nestjs/cli", "^10.0.0")
		addDev("@nestjs/schematics", "^10.0.0")
	}

	// Re-assign maps back to data
	data["dependencies"] = deps
	data["devDependencies"] = devDeps

	// 4. Runtime Adjustments (Bun)
	if config.Runtime == "bun" {
		scripts, _ := data["scripts"].(map[string]interface{})
		if scripts != nil {
			for k, v := range scripts {
				sVal := v.(string)
				sVal = strings.ReplaceAll(sVal, "npm", "bun")
				// Specific tweaks for our templates
				sVal = strings.ReplaceAll(sVal, "tsx watch", "bun --watch")
				sVal = strings.ReplaceAll(sVal, "node ", "bun ")
				scripts[k] = sVal
			}
			data["scripts"] = scripts
		}
	}

	// Return indented JSON
	return json.MarshalIndent(data, "", "  ")
}
