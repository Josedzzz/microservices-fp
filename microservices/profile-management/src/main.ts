import express from "express";
import { initializeDatabase } from "./config/db";
import { startConsumer } from "./messaging/consumer";
import profileRoutes from "./routes/profileRoutes";
import dotenv from "dotenv";
import swaggerUi from "swagger-ui-express";
import YAML from "yamljs";
import path from "path";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 8085;

// Enable strict routing to distinguish between /swagger and /swagger/
app.set("strict routing", true);

app.use(express.json());

// Swagger
const swaggerDocument = YAML.load(path.join(__dirname, "../docs/swagger.yaml"));

// Redirect /swagger to /swagger/ only if exactly /swagger is requested
app.get("/swagger", (req, res) => {
  res.redirect("swagger/");
});

app.use("/swagger/", swaggerUi.serve, swaggerUi.setup(swaggerDocument, {
  swaggerOptions: {
    persistAuthorization: true
  }
}));

// Routes
app.use("/profiles", profileRoutes);

// Health check
app.get("/health", (req, res) => {
  res.json({ status: "UP", service: "profile-management" });
});

async function bootstrap() {
  // Initialize Database
  await initializeDatabase();

  // Start RabbitMQ Consumer
  await startConsumer();

  // Start Server
  app.listen(PORT, () => {
    console.log(`Profile Management Service running on port ${PORT}`);
  });
}

bootstrap().catch((err) => {
  console.error("Critical error starting the service:", err);
  process.exit(1);
});
