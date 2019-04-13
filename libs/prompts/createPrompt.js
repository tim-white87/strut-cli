const inquirer = require('inquirer');

async function createPrompt (product, name) {
  let createPrompt = await inquirer.prompt([{
    type: 'input',
    name: 'name',
    message: 'Enter new strut product:',
    when () { return !name; }
  }]);
  name = name || createPrompt.name;
  await product.create(name);
};

module.exports = createPrompt;
