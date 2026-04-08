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
      // Create a dummy email for USER tests
      email = "user@onboarding.com";
      password = "user123";
    }

    // Attempt login via Auth Service
    await this.http("POST", "/auth-service/api/login", {
      email: email,
      password: password,
    });

    if (this.response.status === 200 && this.response.data.access_token) {
      this.token = this.response.data.access_token;
    } else {
      // Fallback: Generate token manually for test environment
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
    const originalToken = this.token;
    this.token = null;
    await this.http(method, path);
    this.token = originalToken; // Restore for next steps
  },
);

When(
  "I make a {string} request to {string} with a valid employee",
  async function (method, path) {
    const deptId = `DEPT-${Date.now()}`;
    // Create prerequisite department
    await this.http("POST", "/departments-service/api/departments", {
      id: deptId,
      name: "Test Security Dept",
    });

    const employee = {
      name: "Security Test User",
      email: `security-${Date.now()}@example.com`,
      departmentID: deptId,
    };
    await this.http(method, path, employee);
  },
);

// Generic steps
When("I make a {string} request to {string}", async function (method, path) {
  await this.http(method, path);
});

Then("the response status code should be {int}", function (expectedStatus) {
  expect(this.response.status).to.equal(expectedStatus);
});
