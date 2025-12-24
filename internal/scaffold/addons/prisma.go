package addons

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetupPrisma(root string, config AddonConfig) error {
	fmt.Println("   🔹 Configurando Prisma ORM...")

	// 1. Determine Provider and URL
	provider := "postgresql"
	url := "postgresql://johndoe:randompassword@localhost:5432/mydb?schema=public"

	if config.Database == "postgres" {
		provider = "postgresql"
		url = "postgresql://user:password@localhost:5432/mydb?schema=public"
	} else if config.Database == "sqlite" {
		provider = "sqlite"
		url = "file:./dev.db"
	} else if config.Database == "mongo" {
		provider = "mongodb"
		url = "mongodb+srv://root:randompassword@cluster0.ab1cd.mongodb.net/mydb?retryWrites=true&w=majority"
	}

	// 2. Generate Schema
	schemaContent := fmt.Sprintf(`// This is your Prisma schema file,
// learn more about it in the docs: https://pris.ly/d/prisma-schema

generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "%s"
  url      = env("DATABASE_URL")
}

model User {
  id    Int     @id @default(autoincrement())
  email String  @unique
  name  String?
}
`, provider)

	// MongoDB uses ObjectID, so we need a cleaner schema for it
	if config.Database == "mongo" {
		schemaContent = `generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "mongodb"
  url      = env("DATABASE_URL")
}

model User {
  id    String @id @default(auto()) @map("_id") @db.ObjectId
  email String @unique
  name  String?
}
`
	}

	// 3. Write Files
	prismaDir := filepath.Join(root, "prisma")
	if err := os.MkdirAll(prismaDir, 0755); err != nil {
		return err
	}

	if err := writeFile(filepath.Join(prismaDir, "schema.prisma"), schemaContent); err != nil {
		return err
	}

	// 4. Update .env
	envPath := filepath.Join(root, ".env")
	envContent := fmt.Sprintf("\n# Prisma Configuration\nDATABASE_URL=\"%s\"\n", url)

	// Append or Create
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
