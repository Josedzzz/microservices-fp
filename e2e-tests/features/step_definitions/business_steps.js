import { When, Then } from '@cucumber/cucumber';
import { expect } from 'chai';
import { waitUntil } from './polling.js';

When('I create a new employee with the following details:', async function (dataTable) {
  const data = dataTable.hashes()[0];
  
  // Create a unique email to avoid 409 Conflict
  const uniqueId = Date.now();
  const [user, domain] = data.email.split('@');
  const uniqueEmail = `${user}+${uniqueId}@${domain}`;
  
  this.tempData.email = uniqueEmail;
  const employeeData = { ...data, email: uniqueEmail };
  
  // Ensure department exists first
  await this.http('POST', '/departments-service/api/departments', {
    id: data.departmentID,
    name: 'Test Dept'
  });

  const response = await this.http('POST', '/employees-service/api/employees/', employeeData);
  if (response.status === 201) {
    this.tempData.employeeId = response.data.id;
  }
});

When('I delete the current employee', async function () {
  const id = this.tempData.employeeId;
  await this.http('DELETE', `/employees-service/api/employees/${id}`);
});

Then('eventually the current employee should not be able to login', async function () {
  await waitUntil(async () => {
    // Attempt login
    await this.http('POST', '/auth-service/api/login', {
      email: this.tempData.email,
      password: 'newPassword123'
    });
    
    // Condition: Login must return 401 (Unauthorized)
    return this.response.status === 401;
  }, {
    maxAttempts: 10,
    intervalMs: 2000,
    timeoutMessage: `Employee ${this.tempData.email} is still able to login after 20s`
  });
});

Then('eventually a {string} notification should exist for current employee', async function (type) {
  const email = this.tempData.email;
  // Match the actual Java enum name if searching for SECURITY
  const searchType = type === 'SECURITY' || type === 'WELCOME' ? type : type.toUpperCase();

  await waitUntil(async () => {
    const response = await this.http('GET', '/notifications-service/notifications');
    if (response.status !== 200) return false;

    const notifications = response.data;
    // Look for the notification by recipient and type
    const found = notifications.find(n => n.recipient === email && n.type === searchType);
    
    if (found) {
      // Improved Regex to extract UUID token correctly
      const tokenMatch = found.message.match(/token: ([a-f0-9-]{36})/i);
      if (tokenMatch) {
        this.tempData.resetToken = tokenMatch[1];
        
        // Automated password reset using the discovered token
        const resetResponse = await this.http('POST', '/auth-service/api/reset-password', {
          email: email,
          token: this.tempData.resetToken,
          new_password: 'newPassword123'
        });
        
        return resetResponse.status === 200;
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

Then('the response body should contain {string}', function (text) {
  const body = typeof this.response.data === 'string' 
    ? this.response.data 
    : JSON.stringify(this.response.data);
  expect(body).to.include(text);
});
