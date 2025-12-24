# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies (cached)
COPY package.json pnpm-lock.yaml* ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Copy source and build
COPY . .
RUN pnpm build

# Stage 2: Runner
FROM node:20-alpine AS runner

WORKDIR /app

# Don't run as root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 mestre
USER mestre

# Copy only necessary files
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./
COPY --from=builder /app/node_modules ./node_modules
# Copy Prisma Client if it exists (generated in builder)
# COPY --from=builder /app/node_modules/.prisma ./node_modules/.prisma

ENV NODE_ENV production
EXPOSE 3333

CMD ["node", "dist/infra/http/server.js"]
