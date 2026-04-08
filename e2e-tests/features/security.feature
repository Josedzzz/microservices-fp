# language: es
Característica: Seguridad y Control de Acceso (RBAC)
  Como sistema de gestión
  Quiero asegurar que las reglas de acceso se cumplan estrictamente
  Para proteger los recursos del sistema

  Antecedentes:
    Dado que el sistema está desplegado y operativo

  @security @rbac
  Escenario: Acceso denegado cuando no se proporciona un token
    Cuando realizo una petición "GET" a "/employees-service/api/employees" sin un token
    Entonces el código de estado de la respuesta debe ser 401

  @security @rbac
  Escenario: Acceso prohibido cuando un USER intenta crear empleados
    Dado que estoy autenticado con el rol "USER"
    Cuando realizo una petición "POST" a "/employees-service/api/employees"
    Entonces el código de estado de la respuesta debe ser 403

  @security @rbac
  Escenario: Acceso permitido cuando un ADMIN crea un empleado
    Dado que estoy autenticado con el rol "ADMIN"
    Cuando realizo una creación de empleado con datos válidos
    Entonces el código de estado de la respuesta debe ser 201
