import { setWorldConstructor, World, setDefaultTimeout } from '@cucumber/cucumber';
import axios from 'axios';

// Polling operations time! 
setDefaultTimeout(60 * 1000);

class CustomWorld extends World {
  constructor(options) {
    super(options);
    
    // Base configuration
    this.baseUrl = process.env.BASE_URL || 'http://localhost:8080';
    
    // Session state
    this.token = null;
    this.response = null; 
    this.lastError = null;
    
    // Temporary data for cleanup
    this.tempData = {
      employeeId: null,
      email: null,
      resetToken: null
    };
  }

  /**
   * Helper to perform HTTP requests
   * @param {string} method 
   * @param {string} path 
   * @param {object} data 
   * @param {object} options - { headers: {}, saveResponse: true }
   */
  async http(method, path, data = null, options = { saveResponse: true }) {
    const config = {
      method: method.toLowerCase(),
      url: `${this.baseUrl}${path}`,
      headers: {
        'Content-Type': 'application/json',
        ...(options.headers || {})
      },
      validateStatus: () => true 
    };

    if (this.token) {
      config.headers['Authorization'] = `Bearer ${this.token}`;
    }

    if (data) {
      config.data = data;
    }

    try {
      const res = await axios(config);
      if (options.saveResponse !== false) {
        this.response = res;
      }
      return res;
    } catch (error) {
      const errRes = error.response || { 
        status: 500, 
        data: { error: `Connection failed: ${error.message}` } 
      };
      if (options.saveResponse !== false) {
        this.response = errRes;
      }
      this.lastError = error;
      return errRes;
    }
  }
}

setWorldConstructor(CustomWorld);
