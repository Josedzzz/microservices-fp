# language: es
Característica: Flujo de Offboarding de Empleados
  Como administrador del sistema
  Quiero eliminar empleados del sistema
  Para que su acceso sea revocado y ya no puedan iniciar sesión

  @offboarding @async
  Escenario: Eliminación exitosa y revocación de acceso de un empleado
    # Preparación: Crear y activar un empleado primero
    Dado que estoy autenticado con el rol "ADMIN"
    Y registro un nuevo empleado con los siguientes datos:
      | name       | email                | departmentID |
      | Jane Doe   | jane.doe@example.com | DEPT-IT      |
    Y eventualmente debe existir una notificación de tipo "SECURITY" para el empleado actual
    
    # Acción: Eliminar al empleado (Soft-delete)
    Cuando elimino al empleado actual
    Entonces el código de estado de la respuesta debe ser 204
    
    # Verificación administrativa del estado
    Y eventualmente el estado del empleado actual debe ser "RETIRED"

    # Validación Final: El acceso debe estar bloqueado
    Y eventualmente el empleado actual no debe poder iniciar sesión
