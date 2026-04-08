import { setWorldConstructor } from "@cucumber/cucumber";
import axios from "axios";

class CustomWorld {
  constructor() {
    this.baseUrl = process.env.BASE_URL || "http://localhost:8000";
    this.token = null;
    this.response = null;
  }

  /**
   * Helper method to perform HTTP requests
   * @param {string} method - HTTP method (GET, POST, DELETE, etc.)
   * @param {string} path - URL path relative to the base URL
   * @param {object} data - Optional request body
   * @returns {Promise<object>} - Axios response object
   */
  async http(method, path, data = null) {
    const url = `${this.baseUrl}${path}`;
    const headers = {};

    // Add Authorization header if a token is available
    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    try {
      this.response = await axios({
        method,
        url,
        data,
        headers,
        validateStatus: () => true, // Ensure Axios doesn't throw on 4xx/5xx status codes
      });
    } catch (error) {
      this.response = error.response;
    }
    return this.response;
  }
}

setWorldConstructor(CustomWorld);
