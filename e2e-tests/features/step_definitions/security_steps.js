import { When, Given } from "@cucumber/cucumber";
import jwt from "jsonwebtoken";

Given(
  "that I am authenticated as a user with role {string}",
  async function (role) {
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

    await this.http("POST", "/auth-service/api/login", {
      email: email,
      password: password,
    }, { saveResponse: false });

    if (this.response && this.response.status === 200 && this.response.data.access_token) {
      this.token = this.response.data.access_token;
    } else {
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
    this.token = originalToken;
  },
);

When(
  "I make a {string} request to {string} with a valid employee",
  async function (method, path) {
    const deptId = `DEPT-${Date.now()}`;
    await this.http("POST", "/departments-service/api/departments", {
      id: deptId,
      name: "Test Security Dept",
    }, { saveResponse: false });

    const employee = {
      name: "Security Test User",
      email: `security-${Date.now()}@example.com`,
      departmentID: deptId,
    };
    await this.http(method, path, employee);
  },
);

When("I make a {string} request to {string}", async function (method, path) {
  await this.http(method, path);
});
