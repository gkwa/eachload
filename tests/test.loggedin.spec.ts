import { test, expect } from '@playwright/test';

test.use({ storageState: 'playwright/.auth/user.json' });

test.beforeEach(async ({ page }) => {
  await page.goto('/');
});

test('test', async ({ page }) => {
  // await page.getByRole('button', { name: 'Login' }).click();
  // await page.goto('https://myutilities.seattle.gov/eportal/#/account');
  await expect(page.getByRole('link', { name: 'Logout' })).toBeVisible();
  await page.getByRole('link', { name: /^View Usage$/i }).click();
  await page.getByRole('link', { name: 'View Usage', exact: true }).click();
  await page.getByRole('button', { name: 'View Usage Details' }).click();
  await page.getByRole('button', { name: 'Green Button Download my data' }).click();
  await page.getByText('Export usage for a bill period').click();

  const options = await page.$$eval('#period-bill-select option', (els) => {
    return els.map((option) => option.textContent);
  });
  console.log(options);

  page.on('download', async (download) => {
    const downloadPath = await download.path();
    console.log('Downloaded file path:', downloadPath);

    const os = require('os');
    const fs = require('fs').promises;
    const path = require('path');

    const destinationDirectory = './data';
    console.log(destinationDirectory);

    fs.mkdir(destinationDirectory, { recursive: true })
      .then(() => {
        const filename = path.basename(downloadPath);
        const destinationPath = path.resolve(path.join(destinationDirectory, filename));
        console.log('moving file %s to %s', downloadPath, destinationPath);
        return fs.rename(downloadPath, destinationPath);
      })
      .then(() => {
        console.log('File moved successfully');
      })
      .catch((error) => {
        console.error('Error moving file:', error);
      });
  });

  for (const option of options) {
    await page.locator('#period-bill-select').selectOption({ label: option });
    await page.getByRole('button', { name: 'Export' }).click();
  }

  await page.waitForTimeout(10 * 60 * 1000); // enuf time to download export
});
