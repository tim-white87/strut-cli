const colors = require('colors');
const inquirer = require('inquirer');
const productSchemas = require('../products/productSchemas');
const utils = require('../utils');

async function addProviderPrompt (productModel, providerName, applicationNames) {
  console.log(colors.yellow('Lets add a provider to your product.'));
  let prompt = await inquirer.prompt([{
    type: 'list',
    name: 'providerName',
    choices: ['AWS', 'GCP', 'Azure'],
    message: 'Select a provider:',
    when () { return !providerName; }
  }, {
    type: 'checkbox',
    name: 'applicationNames',
    choices (input) {
      let choices = productModel.product.applications
        .filter(a => !a.providers || !a.providers[input.providerName])
        .map(a => {
          return { name: a.name, value: a.name, checked: true };
        });
      return choices;
    },
    message: 'Select application(s):',
    when () {
      if (applicationNames) {
        let appNames = utils.list(applicationNames);
        for (let i = 0; i < appNames; i++) {
          if (!productModel.product.applications.some(a => a.name === appNames[i])) {
            console.log(colors.red(`Application: ${appNames[i].name} not listed in product`));
            return true;
          }
        }
      }
      return !applicationNames;
    }
  }]);
  applicationNames = utils.list(applicationNames) || prompt.applicationNames;
  providerName = providerName || prompt.providerName;
  productModel.product.applications.forEach(app => {
    if (applicationNames.some(name => name === app.name)) {
      app.providers = app.providers || {};
      if (!app.providers[providerName]) {
        app.providers[providerName] = {
          ...productSchemas.provider,
          name: providerName
        };
      }
    }
  });

  await productModel.updateProductFile();

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
