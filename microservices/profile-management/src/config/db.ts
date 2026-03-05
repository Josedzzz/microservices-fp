import "reflect-metadata";
import { DataSource } from "typeorm";
import { Profile } from "../models/Profile";
import dotenv from "dotenv";

dotenv.config();

export const AppDataSource = new DataSource({
  type: "postgres",
  host: process.env.DB_HOST || "localhost",
  port: parseInt(process.env.DB_PORT || "5432"),
  username: process.env.DB_USER || "postgres",
  password: process.env.DB_PASSWORD || "postgres",
  database: process.env.DB_NAME || "profiles_db",
  synchronize: true, // Should be false in production
  logging: false,
  entities: [Profile],
  migrations: [],
  subscribers: [],
});

export const initializeDatabase = async () => {
  try {
    await AppDataSource.initialize();
    console.log("Database connected successfully");
  } catch (error) {
    console.error("Error during Data Source initialization", error);
    process.exit(1);
  }
};
