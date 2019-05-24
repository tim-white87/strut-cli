const colors = require('colors');
const inquirer = require('inquirer');
const productSchemas = require('../products/productSchemas');

async function addApplicationPrompt (productModel) {
  console.log(colors.yellow('Lets add an application to your product.'));
  let prompt = await inquirer.prompt([ {
    type: 'input',
    name: 'path',
    message: 'Please provide the local path to your application:'
  },
  // TODO make repo optional
  {
    type: 'list',
    name: 'repoType',
    choices: ['git', 'svn', 'mercurial']
  }, {
    type: 'input',
    name: 'repoUrl',
    message: 'Provide the remote URL to your code:'
  }, {
    type: 'input',
    name: 'name',
    default (answers) {
      return answers.repoUrl.split('/').pop().split('.')[0];
    }
  }]);
  let application = {
    ...productSchemas.application,
    name: prompt.name,
    path: prompt.path,
    repository: {
      type: prompt.repoType,
      url: prompt.repoUrl
    }
  };
  await productModel.addApplication(application);
  let beginAgainPrompt = await inquirer.prompt([{
    type: 'confirm',
    name: 'beginAgain',
    default: false,
    message: 'Do you want to add another applciation?'
  }]);
  if (beginAgainPrompt.beginAgain) {
    await addApplicationPrompt(productModel);
  }
};

module.exports = addApplicationPrompt;
