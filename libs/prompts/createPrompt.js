const inquirer = require('inquirer');

async function createPrompt (product, value) {
  let createPrompt = await inquirer.prompt([{
    type: 'input',
    name: 'name',
    message: 'Enter new strut product:',
    when () { return !value; }
  }]);
  value = value || createPrompt.name;
  await product.create(value);
};

module.exports = createPrompt;
