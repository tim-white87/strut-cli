const colors = require('colors');
const inquirer = require('inquirer');

async function addApplicationPrompt (product) {
  console.log(colors.yellow('Lets add an application to your product.'));
  let prompt = await inquirer.prompt([{
    type: 'confirm',
    name: 'isLocal',
    message: 'Do you have a repository setup?'
  }, {
    type: 'input',
    name: 'path',
    when () { return this.isLocal; }
  }, {
    type: 'list',
    name: 'repoType',
    choices: ['git', 'svn', 'mercurial'],
    when () { return !this.isLocal; }
  }, {
    type: 'input',
    name: 'repoUrl',
    when () { return !this.isLocal; }
  }, {
    type: 'input',
    name: 'name',
    default (answers) {
      return answers.repoUrl.split('/').pop().split('.')[0];
    },
    when () { return !this.isLocal; }
  }]);
  let application = {
    ...product.applicationModel,
    name: prompt.name,
    path: prompt.path,
    repository: {
      type: prompt.repoType,
      url: prompt.repoUrl
    }
  };
  await product.addApplication(application);
  let beginAgainPrompt = await inquirer.prompt([{
    type: 'confirm',
    name: 'beginAgain',
    default: false,
    message: 'Do you want to add another applciation?'
  }]);
  if (beginAgainPrompt.beginAgain) {
    await addApplicationPrompt(product);
  }
};

module.exports = addApplicationPrompt;
