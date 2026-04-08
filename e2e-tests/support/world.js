const { setWorldConstructor, World } = require('@cucumber/cucumber');
const axios = require('axios');

class CustomWorld extends World {
  constructor(options) {
    super(options);
    
    // Configuración Base
    this.baseUrl = process.env.BASE_URL || 'http://localhost:8080';
    
    // Estado de la sesión
    this.token = null;
    this.lastResponse = null;
    this.lastError = null;
    
    // Datos temporales para limpieza
    this.tempData = {
      employeeId: null,
      email: null
    };
  }

  // Helper para hacer peticiones HTTP con Axios
  async http(method, path, data = null, headers = {}) {
    const config = {
      method: method.toLowerCase(),
      url: `${this.baseUrl}${path}`,
      headers: {
        'Content-Type': 'application/json',
        ...headers
      },
      validateStatus: () => true // No lanzar error en 4xx/5xx para poder validarlos en BDD
    };

    if (this.token) {
      config.headers['Authorization'] = `Bearer ${this.token}`;
    }

    if (data) {
      config.data = data;
    }

    try {
      this.lastResponse = await axios(config);
      return this.lastResponse;
    } catch (error) {
      this.lastError = error;
      throw error;
    }
  }
}

setWorldConstructor(CustomWorld);
