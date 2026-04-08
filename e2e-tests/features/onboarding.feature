# language: es
Característica: Flujo de Onboarding de Empleados
  Como administrador del sistema
  Quiero registrar nuevos empleados
  Para que eventualmente puedan acceder al sistema con sus credenciales

  @onboarding @async
  Escenario: Onboarding completo y exitoso de un nuevo empleado
    Dado que estoy autenticado con el rol "ADMIN"
    Cuando registro un nuevo empleado con los siguientes datos:
      | name       | email                | departmentID |
      | John BDD   | john.bdd@example.com | DEPT-IT      |
    Entonces el código de estado de la respuesta debe ser 201
    
    # Parte asíncrona: Esperando a RabbitMQ
    Y eventualmente debe existir una notificación de tipo "SECURITY" para el empleado actual
    
    # Validación final: Acceso
    Cuando intento hacer login con las credenciales del empleado actual y contraseña "newPassword123"
    Entonces el código de estado de la respuesta debe ser 200
    Y el cuerpo de la respuesta debe contener "access_token"

  @onboarding @negative
  Escenario: Error al crear un empleado con departamento inexistente
    Dado que estoy autenticado con el rol "ADMIN"
    Cuando registro un nuevo empleado con los siguientes datos:
      | name       | email             | departmentID |
      | Ghost User | ghost@example.com | NON-EXISTENT |
    Entonces el código de estado de la respuesta debe ser 400
    Y el cuerpo de la respuesta debe contener "department ID does not exists"
