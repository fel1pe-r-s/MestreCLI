# 🏆 Mestre Stack - Padrões de Tecnologia

Com base na análise de **160+ arquivos de configuração** do ecossistema MestreJS, definimos os seguintes padrões para os Templates Automáticos:

## 🚀 Core (Universal)
*   **Linguagem**: TypeScript (Dominante em 85% dos projetos)
*   **Runtime**: Node.js (Padrão) ou Bun (Opcional)
*   **Gerenciador de Pacotes**: PNPM

## 🌐 Frontend (Web & Mobile)
*   **Framework Web**: Next.js (App Router) ou Vite (SPAs simples)
*   **Mobile**: Expo (React Native)
*   **CSS Style**: TailwindCSS + Tailwind Merge
*   **Ícones**: Lucide React
*   **HTTP Client**: Axios
*   **Validação**: Zod (Presente em 20 projetos)
*   **Datas**: DayJS

## ⚙️ Backend (API)
*   **ORM**: Prisma (Padrão de ouro para SQL)
*   **Frameworks**: Fastify (Performance) ou Express (Legado/Simples)
*   **Config**: Dotenv + Zod (Validação de Variáveis de Ambiente)
*   **Testes**: Vitest (Mais rápido que Jest, compatível com Vite)

## 🐳 Infraestrutura
*   **Container**: Docker + Compose
*   **Imagem Base**: `node:alpine` ou `oven/bun:alpine`

---
> *Este documento guia as decisões arquiteturais do MestreCLI.*
