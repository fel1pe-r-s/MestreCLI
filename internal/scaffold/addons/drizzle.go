package addons

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetupDrizzle(root string, config AddonConfig) error {
	fmt.Println("   🔹 Configurando Drizzle ORM...")

	// 1. Determine Driver/URL
	driver := "pg" // default
	url := "postgresql://user:password@localhost:5432/mydb"

	if config.Database == "postgres" {
		driver = "pg"
		url = "postgresql://user:password@localhost:5432/mydb"
	} else if config.Database == "sqlite" {
		driver = "better-sqlite"
		url = "sqlite.db"
	} else if config.Database == "mysql" {
		driver = "mysql2"
		url = "mysql://user:password@localhost:3306/mydb"
	}

	// 2. Generate drizzle.config.ts
	// Note: Drizzle config format varies slightly, keeping it standard V5 style
	configContent := fmt.Sprintf(`import type { Config } from "drizzle-kit";

export default {
  schema: "./src/db/schema.ts",
  out: "./drizzle",
  driver: "%s",
  dbCredentials: {
    connectionString: process.env.DATABASE_URL!,
  },
} satisfies Config;
`, driver)

	if err := writeFile(filepath.Join(root, "drizzle.config.ts"), configContent); err != nil {
		return err
	}

	// 3. Create Schema Example
	schemaDir := filepath.Join(root, "src", "db")
	if err := os.MkdirAll(schemaDir, 0755); err != nil {
		return err
	}

	schemaContent := `import { pgTable, serial, text, varchar } from "drizzle-orm/pg-core";

export const users = pgTable('users', {
  id: serial('id').primaryKey(),
  fullName: text('full_name'),
  phone: varchar('phone', { length: 256 }),
});
`
	if config.Database == "sqlite" {
		schemaContent = `import { sqliteTable, text, integer } from 'drizzle-orm/sqlite-core';

export const users = sqliteTable('users', {
  id: integer('id').primaryKey(), // auto-increment
  fullName: text('full_name'),
});`
	}

	if err := writeFile(filepath.Join(schemaDir, "schema.ts"), schemaContent); err != nil {
		return err
	}

	// 4. Update .env
	envPath := filepath.Join(root, ".env")
	envContent := fmt.Sprintf("\n# Drizzle Configuration\nDATABASE_URL=\"%s\"\n", url)

	f, err := os.OpenFile(envPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(envContent); err != nil {
		return err
	}

	return nil
}
