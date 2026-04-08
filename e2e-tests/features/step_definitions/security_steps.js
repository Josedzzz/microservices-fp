import { When, Then, Given } from "@cucumber/cucumber";
import jwt from "jsonwebtoken";

Given(
  "que estoy autenticado con el rol {string}",
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
  "realizo una petición {string} a {string} sin un token",
  async function (method, path) {
    const originalToken = this.token;
    this.token = null;
    await this.http(method, path);
    this.token = originalToken;
  },
);

When(
  "realizo una creación de empleado con datos válidos",
  async function () {
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
    await this.http("POST", "/employees-service/api/employees/", employee);
  },
);

When("realizo una petición {string} a {string}", async function (method, path) {
  await this.http(method, path);
});

Given("que el sistema está desplegado y operativo", async function () {
  await this.http("GET", "/", null, { saveResponse: false });
});
