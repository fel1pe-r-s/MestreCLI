import { fastify } from "fastify";
// import express from "express"; // CLI will uncomment this if needed

async function bootstrap() {
  // Configuração Padrão: Fastify
  const app = fastify();
  
  app.get("/health", async () => {
    return { status: "ok", framework: "fastify" };
  });

  try {
    await app.listen({ port: 3333 });
    console.log("🚀 Server running on http://localhost:3333 (Fastify)");
  } catch (err) {
    app.log.error(err);
    process.exit(1);
  }
}

bootstrap();
