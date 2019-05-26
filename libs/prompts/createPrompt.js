const inquirer = require('inquirer');

async function createPrompt (productModel, name) {
  let createPrompt = await inquirer.prompt([{
    type: 'input',
    name: 'name',
    default: process.cwd().split('/').pop(),
    message: 'Enter new strut product:',
    when () { return !name; }
  }, {
    type: 'list',
    name: 'fileType',
    default: 'yml',
    choices: [{
      name: 'YAML',
      value: 'yml'
    }, {
      name: 'JSON',
      value: 'json'
    }]
  }]);
  name = name || createPrompt.name;
  await productModel.create(name, createPrompt.fileType === 'yml');
};

module.exports = createPrompt;
