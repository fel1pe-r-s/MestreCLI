package addons

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetupHono(root string) error {
	fmt.Println("   🔥 Configurando Hono (Edge Adapter)...")

	// 1. Create API Route Adapter
	apiRoutePath := filepath.Join(root, "app", "api", "[[...route]]")
	if err := os.MkdirAll(apiRoutePath, 0755); err != nil {
		return err
	}

	routeContent := `import { Hono } from 'hono'
import { handle } from 'hono/vercel'

export const runtime = 'edge'

const app = new Hono().basePath('/api')

app.get('/hello', (c) => {
  return c.json({
    message: 'Hello from Hono on Edge!',
  })
})

export const GET = handle(app)
export const POST = handle(app)
`
	if err := writeFile(filepath.Join(apiRoutePath, "route.ts"), routeContent); err != nil {
		return err
	}

	// 2. We need to add 'hono' to package.json.
	// This is currently handled by the Patcher if 'ApiPattern' is set,
	// or we can force it here if Patcher logic is complex.
	// For simplicity, let's assume Patcher handles the dependency injection based on config.

	return nil
}
