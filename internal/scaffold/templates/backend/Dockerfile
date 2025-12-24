# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /usr/src/app

COPY package.json pnpm-lock.yaml* ./
RUN npm install -g pnpm && pnpm install

COPY . .
RUN pnpm build

# Stage 2: Runner
FROM node:20-alpine

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/dist ./dist
COPY --from=builder /usr/src/app/package.json ./
COPY --from=builder /usr/src/app/node_modules ./node_modules

EXPOSE 3333

CMD ["node", "dist/infra/http/server.js"]
