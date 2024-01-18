import { test } from '@playwright/test';

test('test', async ({ page }) => {
  const USERNAME = process.env.SEATTLE_UTILITIES_USERNAME;
  const PASSWORD = process.env.SEATTLE_UTILITIES_PASSWORD;

  if (typeof USERNAME === 'undefined' || typeof PASSWORD === 'undefined') {
    console.error('Please set env vars SEATTLE_UTILITIES_USERNAME and SEATTLE_UTILITIES_PASSWORD');
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
  await page.getByRole('link', { name: /^View Usage$/i }).click();
  await page.getByRole('link', { name: 'View Usage', exact: true }).click();
  await page.getByRole('button', { name: 'View Usage Details' }).click();
  await page.getByRole('button', { name: 'Green Button Download my data' }).click();
  await page.getByText('Export usage for a bill period').click();

  page.pause();

  const options = await page.$$eval('#period-bill-select option', (els) => {
    return els.map((option) => option.textContent);
  });
  console.log(options);

  //[
  //  'Since your last bill: Nov 03, 2023 - Jan 06, 2024',
  //  'Sep 07, 2023 - Nov 03, 2023',
  //  'Jul 11, 2023 - Sep 06, 2023',
  //  'May 09, 2023 - Jul 10, 2023',
  //  'Mar 11, 2023 - May 08, 2023',
  //  'Jan 10, 2023 - Mar 10, 2023',
  //  'Nov 04, 2022 - Jan 09, 2023',
  //  'Sep 07, 2022 - Nov 03, 2022'
  //]

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
    const selectOption = option !== null ? { label: option } : { label: undefined };
    if (selectOption.label === undefined) {
      alert('selectOption is undefined!');
    }

    await page.locator('#period-bill-select').selectOption(selectOption);
    await page.getByRole('button', { name: 'Export' }).click();
  }

  await page.waitForTimeout(10 * 1000); // 10 seconds
});
