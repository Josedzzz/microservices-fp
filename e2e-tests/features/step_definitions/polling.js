/**
 * Función de Polling reutilizable para esperar condiciones asincrónicas
 */
export async function waitUntil(condition, options = {}) {
  const {
    maxAttempts = 15,
    intervalMs = 2000,
    timeoutMessage = "La condición no se cumplió tras agotar el tiempo de espera."
  } = options;

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      const result = await condition();
      if (result) return true;
    } catch (e) {}

    if (attempt < maxAttempts) {
      await new Promise(resolve => setTimeout(resolve, intervalMs));
    }
  }

  throw new Error(timeoutMessage);
}
