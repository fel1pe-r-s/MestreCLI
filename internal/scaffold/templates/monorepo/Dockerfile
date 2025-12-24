FROM node:20-alpine AS base

# 1. Prune (Select only the app we are building)
FROM base AS pruner
RUN apk add --no-cache libc6-compat
WORKDIR /app
RUN npm install -g turbo
COPY . .
# Prune for 'web' app (default for universal)
RUN turbo prune --scope=web --docker

# 2. Install Dependencies (from pruned lockfile)
FROM base AS installer
RUN apk add --no-cache libc6-compat
WORKDIR /app

# First install dependencies (for caching)
COPY .gitignore .gitignore
COPY --from=pruner /app/out/json/ .
COPY --from=pruner /app/out/pnpm-lock.yaml ./pnpm-lock.yaml
RUN npm install -g pnpm
RUN pnpm install

# Then copy source code
COPY --from=pruner /app/out/full/ .
COPY turbo.json turbo.json

# 3. Build
FROM installer AS builder
WORKDIR /app
RUN npm install -g pnpm
RUN pnpm turbo run build --filter=web...

# 4. Runner
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
# Create non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
USER nextjs

# Copy built assets
COPY --from=builder /app/apps/web/public ./apps/web/public
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/static ./apps/web/.next/static

EXPOSE 3000
CMD ["node", "apps/web/server.js"]
