import { After } from '@cucumber/cucumber';

// Teardown logic to clean up after each scenario
After(async function (scenario) {
  // If an employee was created during the test, try to delete it to leave a clean state
  if (this.tempData.employeeId && this.token) {
    try {
      console.log(`[Teardown] Cleaning up employee: ${this.tempData.employeeId}`);
      await this.http('DELETE', `/employees-service/api/employees/${this.tempData.employeeId}`, null, { saveResponse: false });
    } catch (e) {
      // Ignore errors during teardown
    }
  }
});
