/**
 * Función de Polling reutilizable para esperar condiciones asincrónicas
 * @param {Function} condition - Una función asíncrona que retorna true si la condición se cumple
 * @param {Object} options - Configuración de reintentos
 */
async function waitUntil(condition, options = {}) {
  const {
    maxAttempts = 15,
    intervalMs = 2000,
    timeoutMessage = "La condición no se cumplió tras agotar el tiempo de espera."
  } = options;

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      const result = await condition();
      if (result) return true;
    } catch (e) {
      // Ignorar errores durante el polling (ej. 404 porque el dato aún no existe)
    }

    if (attempt < maxAttempts) {
      await new Promise(resolve => setTimeout(resolve, intervalMs));
    }
  }

  throw new Error(timeoutMessage);
}

module.exports = { waitUntil };
