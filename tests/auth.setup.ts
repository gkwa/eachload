import { test as setup, request } from '@playwright/test';

const authFile = 'playwright/.auth/user.json';

setup('authenticate', async ({ page }) => {
  const USERNAME = process.env.SEATTLE_UTILITIES_USERNAME;
  const PASSWORD = process.env.SEATTLE_UTILITIES_PASSWORD;

  if (typeof USERNAME === 'undefined' || typeof PASSWORD === 'undefined') {
    process.exit(1);
  }

  await page.goto('https://myutilities.seattle.gov/eportal');
  await page.getByRole('link', { name: 'Login', exact: true }).click();
  await page.getByLabel('Username *').click();
  await page.getByLabel('Username *').click();
  await page.getByLabel('Username *').fill(USERNAME);
  await page.getByPlaceholder('Password').click();
  await page.getByPlaceholder('Password').click();
  await page.getByPlaceholder('Password').fill(PASSWORD);
  await page.getByRole('button', { name: 'Login' }).click();

  await page.context().storageState({ path: authFile });
});
