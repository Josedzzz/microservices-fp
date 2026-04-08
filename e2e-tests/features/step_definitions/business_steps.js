import { When, Then } from '@cucumber/cucumber';
import { expect } from 'chai';
import { waitUntil } from './polling.js';

When('I create a new employee with the following details:', async function (dataTable) {
  const data = dataTable.hashes()[0];
  const uniqueId = Date.now();
  
  // 1. Generate unique email
  const [user, domain] = data.email.split('@');
  const uniqueEmail = `${user}+${uniqueId}@${domain}`;
  this.tempData.email = uniqueEmail;

  // 2. Solution 1: Dynamic ID for negative cases
  const targetDeptId = data.departmentID === 'NON-EXISTENT' 
    ? `INVALID-DEPT-${uniqueId}` 
    : data.departmentID;

  const employeeData = { ...data, email: uniqueEmail, departmentID: targetDeptId };
  
  if (data.departmentID !== 'NON-EXISTENT') {
    // Ensure department exists for positive cases
    await this.http('POST', '/departments-service/api/departments', {
      id: targetDeptId,
      name: 'Test Dept'
    }, { saveResponse: false });
  }
  
  // 3. Main action
  await this.http('POST', '/employees-service/api/employees/', employeeData);
  
  if (this.response.status === 201) {
    this.tempData.employeeId = this.response.data.id;
  }
});

When('I delete the current employee', async function () {
  const id = this.tempData.employeeId;
  await this.http('DELETE', `/employees-service/api/employees/${id}`);
});

Then('eventually the current employee should not be able to login', async function () {
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

Then('eventually the current employee status should be {string}', async function (expectedStatus) {
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

Then('eventually a {string} notification should exist for current employee', async function (type) {
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

When('I login with the current employee credentials and password {string}', async function (password) {
  await this.http('POST', '/auth-service/api/login', {
    email: this.tempData.email,
    password
  });
});

Then('the response status code should be {int}', function (code) {
  expect(this.response.status).to.equal(code);
});

Then('the response body should contain {string}', function (text) {
  const body = typeof this.response.data === 'string' 
    ? this.response.data 
    : JSON.stringify(this.response.data);
  expect(body).to.include(text);
});
