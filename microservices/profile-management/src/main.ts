import express from "express";
import { initializeDatabase } from "./config/db";
import { startConsumer } from "./messaging/consumer";
import profileRoutes from "./routes/profileRoutes";
import dotenv from "dotenv";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 8085;

app.use(express.json());

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
