# language: es
Característica: Verificación inicial del sistema (Humo)
  Como desarrollador de pruebas
  Quiero verificar que el Gateway responde correctamente
  Para asegurar que la infraestructura de pruebas está lista

  @smoke
  Escenario: El Gateway está operativo y muestra el directorio
    Cuando realizo una petición "GET" a "/"
    Entonces el código de estado de la respuesta debe ser 200
    Y el cuerpo de la respuesta debe contener "Service Directory"
