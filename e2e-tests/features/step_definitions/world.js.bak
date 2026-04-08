import { setWorldConstructor, World } from '@cucumber/cucumber';
import axios from 'axios';

class CustomWorld extends World {
  constructor(options) {
    super(options);
    
    // Base configuration
    this.baseUrl = process.env.BASE_URL || 'http://localhost:8080';
    
    // Session state
    this.token = null;
    this.response = null; // Matches teammate's naming convention
    this.lastError = null;
    
    // Temporary data for cleanup
    this.tempData = {
      employeeId: null,
      email: null
    };
  }

  // Helper to perform HTTP requests
  async http(method, path, data = null, headers = {}) {
    const config = {
      method: method.toLowerCase(),
      url: `${this.baseUrl}${path}`,
      headers: {
        'Content-Type': 'application/json',
        ...headers
      },
      validateStatus: () => true // Don't throw on 4xx/5xx
    };

    if (this.token) {
      config.headers['Authorization'] = `Bearer ${this.token}`;
    }

    if (data) {
      config.data = data;
    }

    try {
      this.response = await axios(config);
      return this.response;
    } catch (error) {
      this.response = error.response;
      this.lastError = error;
      return this.response;
    }
  }
}

setWorldConstructor(CustomWorld);
