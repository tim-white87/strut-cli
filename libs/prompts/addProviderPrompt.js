const colors = require('colors');
const inquirer = require('inquirer');
const productSchemas = require('../products/productSchemas');

async function addProviderPrompt (productModel, providerName) {
  console.log(colors.yellow('Lets add a provider to your product.'));
  let prompt = await inquirer.prompt([{
    type: 'list',
    name: 'providerName',
    choices: ['AWS', 'GCP', 'Azure'],
    message: 'Select a provider:'
  }]);
  let provider = {
    ...productSchemas.provider,
    name: prompt.providerName
  };
  await productModel.addProvider(provider);

  let beginAgainPrompt = await inquirer.prompt([{
    type: 'confirm',
    name: 'beginAgain',
    default: false,
    message: 'Do you want to add another provider?'
  }]);
  if (beginAgainPrompt.beginAgain) {
    await addProviderPrompt(productModel, providerName);
  }
  return provider;
};

module.exports = addProviderPrompt;
