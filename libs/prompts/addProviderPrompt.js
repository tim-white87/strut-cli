const colors = require('colors');
const inquirer = require('inquirer');
const productSchemas = require('../products/productSchemas');

async function addProviderPrompt (productModel, providerName, applicationName) {
  console.log(colors.yellow('Lets add a provider to your product.'));
  let prompt = await inquirer.prompt([{
    type: 'list',
    name: 'providerName',
    choices: ['AWS', 'GCP', 'Azure'],
    message: 'Select a provider:',
    when () { return !providerName; }
  }, {
    type: 'list',
    name: 'application',
    choices: productModel.product.applications.map(a => {
      return { name: a.name, value: a };
    }),
    message: 'Select application:',
    when () {
      if (applicationName && !productModel.product.applications.some(a => a.name === applicationName)) {
        console.log(colors.red('Application not listed in product'));
        return true;
      }
      return !applicationName;
    }
  }]);
  let application = productModel.product.applications.find(a => a.name === applicationName) || prompt.application;
  providerName = providerName || prompt.providerName;
  console.log(application);
  application.providers[providerName] = {
    ...productSchemas.provider,
    name: providerName
  };
  await productModel.updateApplication(application);

  let beginAgainPrompt = await inquirer.prompt([{
    type: 'confirm',
    name: 'beginAgain',
    default: false,
    message: 'Do you want to add another provider?'
  }]);
  if (beginAgainPrompt.beginAgain) {
    await addProviderPrompt(productModel);
  }
};

module.exports = addProviderPrompt;
