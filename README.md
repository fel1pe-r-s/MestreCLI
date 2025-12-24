# 🐹 Mestre CLI

A ferramenta oficial de automação para o ecossistema **MestreJS**. Gere projetos padronizados em segundos.

## 🚀 Instalação

### Pré-requisitos
*   Go 1.21+

### Compilando Localmente
```bash
git clone https://github.com/fel1pe-r-s/MestreCLI.git
cd MestreCLI
go build -o mestre
```

### Adicionando ao Path (Linux/Mac)
Para usar o comando `mestre` em qualquer terminal:
```bash
sudo mv mestre /usr/local/bin/
```

## 📖 Como Usar

### Iniciar um Novo Projeto
O comando principal é o `init`, que abre um assistente interativo (Wizard).

```bash
mestre init
```

Você será guiado pelas seguintes escolhas:
1.  **Backend API**: Gera uma API com Clean Architecture (escolha entre Express/Fastify).
2.  **Universal App**: Gera um Monorepo com Next.js (Web), Expo (Mobile) e Tauri (Desktop).
3.  **Monorepo Genérico**: Estrutura base para múltiplos pacotes/apps.

## 🏗️ Templates Disponíveis

O CLI baixa automaticamente os templates oficiais:
*   [MestreJS_Backend](https://github.com/fel1pe-r-s/MestreJS_Backend)
*   [MestreJS_Universal](https://github.com/fel1pe-r-s/MestreJS_Universal)
*   [MestreJS_Monorepo](https://github.com/fel1pe-r-s/MestreJS_Monorepo)

## 🛠️ Desenvolvimento

```bash
# Rodar sem compilar
go run main.go init
```
