const inquirer = require('inquirer');

async function linkPrompt (productModel, applicationName, providerName) {
  let application = productModel.product.applications.find(a => a.name === applicationName);
  let provider = productModel.product.providers.find(p => p.name === providerName);
  let prompt = await inquirer.prompt([{
    type: 'list',
    name: 'application',
    choices: productModel.product.applications.map(a => {
      return { name: a.name, value: a };
    }),
    message: 'Select application:',
    when () { return !application; }
  }, {
    type: 'list',
    name: 'provider',
    choices: productModel.product.providers.map(p => {
      return { name: p.name, value: p };
    }),
    message: 'Select provider:',
    when () { return !provider; }
  }]);
  application = application || prompt.application;
  provider = provider || prompt.provider;
  await productModel.link(application, provider);
};

module.exports = linkPrompt;
