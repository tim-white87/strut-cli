const colors = require('colors');
const inquirer = require('inquirer');

async function addApplicationPrompt (product) {
  console.log(colors.yellow('Lets add an application to your product.'));
  let addApplicationPrompt = await inquirer.prompt([{
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
    name: addApplicationPrompt.name,
    path: addApplicationPrompt.path,
    repository: {
      type: addApplicationPrompt.repoType,
      url: addApplicationPrompt.repoUrl
    }
  };
  await product.addApplication(application);
};

module.exports = addApplicationPrompt;
