import { test, expect } from '@playwright/test';

test('el usuario puede resolver un problema simplex y ver el resultado', async ({ page }) => {
  // Cambia la URL si tu Vite server corre en otro puerto
  await page.goto('http://localhost:5173');

  // Cambia los coeficientes de la función objetivo
  const objectiveInputs = page.locator('.objective-inputs input[type="number"]');
  await objectiveInputs.nth(0).fill('1.1');
  await objectiveInputs.nth(1).fill('2.2');

  // Cambia los valores de la primera restricción
  const constraintInputs = page.locator('.constraints-table .constraint-row').nth(0).locator('input[type="number"]');
  await constraintInputs.nth(0).fill('3.3');
  await constraintInputs.nth(1).fill('4.4');
  await constraintInputs.nth(2).fill('5.5'); // Lado derecho

  // Cambia los valores de la segunda restricción
  const constraintInputs2 = page.locator('.constraints-table .constraint-row').nth(1).locator('input[type="number"]');
  await constraintInputs2.nth(0).fill('6.6');
  await constraintInputs2.nth(1).fill('7.7');
  await constraintInputs2.nth(2).fill('8.8'); // Lado derecho

  // Haz clic en el botón "Resolver problema"
  await page.click('button.solve-button');

  // Espera el resultado y verifica que aparezca el valor óptimo
  await expect(page.locator('.result-section')).toContainText('Valor óptimo');
  await expect(page.locator('.result-section')).toContainText('x1');
  await expect(page.locator('.result-section')).toContainText('x2');
});
