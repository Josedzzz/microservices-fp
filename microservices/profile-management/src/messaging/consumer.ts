import amqp from "amqplib";
import { AppDataSource } from "../config/db";
import { Profile } from "../models/Profile";
import dotenv from "dotenv";

dotenv.config();

const RABBITMQ_URL = `amqp://${process.env.RABBITMQ_USER || "guest"}:${process.env.RABBITMQ_PASSWORD || "guest"}@${process.env.RABBITMQ_HOST || "localhost"}:${process.env.RABBITMQ_PORT || "5672"}`;
const EXCHANGE_NAME = "employees.events";
const ROUTING_KEY = "employee.created";
const QUEUE_NAME = "profile_management_queue";

export const startConsumer = async () => {
  try {
    const connection = await amqp.connect(RABBITMQ_URL);
    const channel = await connection.createChannel();

    await channel.assertExchange(EXCHANGE_NAME, "topic", { durable: true });
    const q = await channel.assertQueue(QUEUE_NAME, { durable: true });

    await channel.bindQueue(q.queue, EXCHANGE_NAME, ROUTING_KEY);

    console.log(`[*] Waiting for messages in ${q.queue}. To exit press CTRL+C`);

    channel.consume(q.queue, async (msg) => {
      if (msg !== null) {
        try {
          const content = JSON.parse(msg.content.toString());
          console.log(`[x] Received event: ${msg.fields.routingKey}`, content);

          if (msg.fields.routingKey === ROUTING_KEY) {
            await createProfileFromEvent(content);
          }

          channel.ack(msg);
        } catch (error) {
          console.error("Error processing message:", error);
          // nack if transient error, or ack if we want to skip malformed messages
          channel.nack(msg, false, false);
        }
      }
    });
  } catch (error) {
    console.error("Error starting RabbitMQ consumer:", error);
    // Retry connection logic could be added here
    setTimeout(startConsumer, 5000);
  }
};

async function createProfileFromEvent(eventData: any) {
  const profileRepository = AppDataSource.getRepository(Profile);

  // Check if profile already exists
  const existingProfile = await profileRepository.findOne({
    where: { employeeId: eventData.id.toString() }
  });

  if (existingProfile) {
    console.log(`Profile for employee ${eventData.id} already exists. Skipping.`);
    return;
  }

  const newProfile = new Profile();
  newProfile.employeeId = eventData.id.toString();
  newProfile.name = eventData.name;
  newProfile.email = eventData.email;
  // Others are default empty strings

  await profileRepository.save(newProfile);
  console.log(`Profile created for employee: ${eventData.name} (${eventData.id})`);
}
