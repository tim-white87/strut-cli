const colors = require('colors');
const inquirer = require('inquirer');

async function addProviderPrompt (product, providerName) {
  console.log(colors.yellow('Lets add a provider to your product.'));
  let prompt = await inquirer.prompt([{
    type: 'list',
    name: 'providerName',
    choices: ['AWS', 'GCP', 'Azure'],
    message: 'Select a provider:'
  }]);
  let provider = {
    ...product.providerModel,
    name: prompt.providerName
  };
  await product.addProvider(provider);

  let beginAgainPrompt = await inquirer.prompt([{
    type: 'confirm',
    name: 'beginAgain',
    default: false,
    message: 'Do you want to add another provider?'
  }]);
  if (beginAgainPrompt.beginAgain) {
    await addProviderPrompt(product, providerName);
  }
};

module.exports = addProviderPrompt;
