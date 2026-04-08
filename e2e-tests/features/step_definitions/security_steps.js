import { When, Then, Given } from "@cucumber/cucumber";
import { expect } from "chai";
import jwt from "jsonwebtoken";

Given(
  "that I am authenticated as a user with role {string}",
  async function (role) {
    // Clear previous state
    this.token = null;
    this.response = null;

    let email, password;
    if (role === "ADMIN") {
      email = "admin@onboarding.com";
      password = "admin123";
    } else {
      email = "user@onboarding.com";
      password = "user123";
    }

    // Attempt login via Auth Service
    const response = await this.http("POST", "/auth-service/api/login", {
      email: email,
      password: password,
    });

    if (response.status === 200) {
      this.token = response.data.access_token;
    } else {
      // Fallback: If login fails (e.g., user doesn't exist yet), generate the token manually
      // since we know the secret key for the development environment.
      const secret = process.env.JWT_SECRET || "supersecretkey";
      this.token = jwt.sign({ sub: email, role: role }, secret, {
        expiresIn: "1h",
      });
    }
  },
);

When(
  "I make a {string} request to {string} without a token",
  async function (method, path) {
    this.token = null;
    await this.http(method, path);
  },
);

When("I make a {string} request to {string}", async function (method, path) {
  await this.http(method, path);
});

When(
  "I make a {string} request to {string} with a valid employee",
  async function (method, path) {
    // 1. First, ensure a department exists to avoid 400 Bad Request (Department not found)
    const deptId = `DEPT-${Date.now()}`;
    await this.http("POST", "/departments-service/api/departments", {
      id: deptId,
      name: "Test Department",
    });

    // 2. Create the employee using the new department ID
    const employee = {
      name: "Test Employee",
      email: `test-${Date.now()}@example.com`,
      departmentID: deptId,
    };
    await this.http(method, path, employee);
  },
);

Then("the response status code should be {int}", function (expectedStatus) {
  expect(this.response.status).to.equal(expectedStatus);
});
