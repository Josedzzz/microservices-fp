import { When, Then } from '@cucumber/cucumber';
import { expect } from 'chai';
import { waitUntil } from './polling.js';

When('registro un nuevo empleado con los siguientes datos:', async function (dataTable) {
  const data = dataTable.hashes()[0];
  const uniqueId = Date.now();
  const [user, domain] = data.email.split('@');
  const uniqueEmail = `${user}+${uniqueId}@${domain}`;
  this.tempData.email = uniqueEmail;

  const targetDeptId = data.departmentID === 'NON-EXISTENT' 
    ? `INVALID-DEPT-${uniqueId}` 
    : data.departmentID;

  const employeeData = { ...data, email: uniqueEmail, departmentID: targetDeptId };
  
  if (data.departmentID !== 'NON-EXISTENT') {
    await this.http('POST', '/departments-service/api/departments', {
      id: targetDeptId,
      name: 'Test Dept'
    }, { saveResponse: false });
  }
  
  const response = await this.http('POST', '/employees-service/api/employees/', employeeData);
  if (response.status === 201) {
    this.tempData.employeeId = response.data.id;
  }
});

When('elimino al empleado actual', async function () {
  const id = this.tempData.employeeId;
  await this.http('DELETE', `/employees-service/api/employees/${id}`);
});

Then('eventualmente el empleado actual no debe poder iniciar sesión', async function () {
  await waitUntil(async () => {
    await this.http('POST', '/auth-service/api/login', {
      email: this.tempData.email,
      password: 'newPassword123'
    });
    return this.response.status === 401;
  }, {
    maxAttempts: 10,
    intervalMs: 2000,
    timeoutMessage: `Employee ${this.tempData.email} is still able to login after 20s`
  });
});

Then('eventualmente el estado del empleado actual debe ser {string}', async function (expectedStatus) {
  const id = this.tempData.employeeId;
  await waitUntil(async () => {
    await this.http('GET', `/employees-service/api/employees/${id}`, null, { saveResponse: true });
    return this.response.status === 200 && this.response.data.status === expectedStatus;
  }, {
    maxAttempts: 10,
    intervalMs: 2000,
    timeoutMessage: `Employee ${id} status is not ${expectedStatus} after 20s`
  });
});

Then('eventualmente debe existir una notificación de tipo {string} para el empleado actual', async function (type) {
  const email = this.tempData.email;
  const searchType = type === 'SECURITY' || type === 'WELCOME' ? type : type.toUpperCase();
  
  await waitUntil(async () => {
    const response = await this.http('GET', '/notifications-service/notifications', null, { saveResponse: false });
    if (response.status !== 200) return false;
    
    const notifications = response.data;
    const found = notifications.find(n => n.recipient === email && n.type === searchType);
    
    if (found) {
      const tokenMatch = found.message.match(/token: ([a-f0-9-]{36})/i);
      if (tokenMatch) {
        this.tempData.resetToken = tokenMatch[1];
        await this.http('POST', '/auth-service/api/reset-password', {
          email: email,
          token: this.tempData.resetToken,
          new_password: 'newPassword123'
        }, { saveResponse: false });
        return true;
      }
    }
    return false;
  }, {
    maxAttempts: 15,
    intervalMs: 2000,
    timeoutMessage: `Could not find or process ${searchType} notification for ${email} after 30s`
  });
});

When('intento hacer login con las credenciales del empleado actual y contraseña {string}', async function (password) {
  await this.http('POST', '/auth-service/api/login', {
    email: this.tempData.email,
    password
  });
});

Then('el código de estado de la respuesta debe ser {int}', function (code) {
  expect(this.response.status).to.equal(code);
});

Then('el cuerpo de la respuesta debe contener {string}', function (text) {
  const body = typeof this.response.data === 'string' 
    ? this.response.data 
    : JSON.stringify(this.response.data);
  expect(body).to.include(text);
});
