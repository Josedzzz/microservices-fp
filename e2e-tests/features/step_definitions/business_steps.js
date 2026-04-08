import { Then } from '@cucumber/cucumber';
import { expect } from 'chai';

// Note: "When I make a {string} request to {string}" 
// and "Then the response status code should be {int}" 
// are already defined in security_steps.js by the teammate.

Then('the response body should contain {string}', function (text) {
  const body = typeof this.response.data === 'string' 
    ? this.response.data 
    : JSON.stringify(this.response.data);
  expect(body).to.include(text);
});
