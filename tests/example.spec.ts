import { test, expect } from '@playwright/test';
import fs from 'fs/promises';
import path from 'path';

test('has title', async ({ page }) => {
  const USERNAME = process.env.SEATTLE_UTILITIES_USERNAME;
  const PASSWORD = process.env.SEATTLE_UTILITIES_PASSWORD;
  // console.log(USERNAME);
  // console.log(PASSWORD);

  if (typeof USERNAME === 'undefined' || typeof PASSWORD === 'undefined') {
    process.exit(1);
  }

  await page.goto('https://myutilities.seattle.gov/eportal');
  await page
    .locator('a')
    .filter({ hasText: /^Login $/ })
    .click();
  await page.type('#userName', USERNAME);
  await page.locator('input[type=password]').type(PASSWORD);
  await page.getByRole('button', { name: 'Login' }).click();
  await page.getByRole('link', { name: /^View Usage$/i }).click();
  await page.waitForTimeout(2 * 1000);
  await page.getByRole('link', { name: /^View Usage$/i }).click();
  await page.getByRole('button', { name: 'View Usage Details' }).click();
  await page.getByRole('button', { name: 'ENERGY USE' }).click();
  await page.getByLabel('Change view').selectOption({ label: 'Day view' });
  await page
    .locator('span')
    .filter({ hasText: /Download my data/ })
    .click();
  await page.keyboard.press('PageDown');
  await page.keyboard.press('PageDown');
  await page.keyboard.press('PageDown');
  await page
    .locator('label')
    .filter({ hasText: /Export usage for a bill period/ })
    .click();
  await page.locator('label').filter({ hasText: /CSV/ }).click();
  await page.getByRole('button', { name: 'EXPORT' }).click();

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

  await page.waitForTimeout((1 / 2) * 60 * 1000); // enuf time to download export
});
