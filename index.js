#!/bin/node
const program = require('commander');
const inquirer = require('inquirer');
const process = require('process');
const colors = require('colors');
const { Product } = require('./libs/product');
const product = new Product();

console.log(colors.blue('Welcome to Strut!'));

program
  .version(require('./package.json').version)
  .arguments('<cmd> [value]')
  .action(async (cmd, value) => {
    await product.init();
    switch (cmd) {
      case 'create':
        await create(value);
        process.chdir(`./${value}`);
        await addApplication();
        break;
      default:
        console.log(colors.red(`'${cmd}' command does not exist, try --help for valid commands`));
        break;
    }
  })
  .on('--help', function () {
    console.log('');
    console.log('Commands:');
    console.log('  create [name]    Create a new product definition');
  })
  .parse(process.argv);

async function create (value) {
  let createPrompt = await inquirer.prompt([{
    type: 'input',
    name: 'name',
    message: 'Enter new strut product:',
    when () { return !value; }
  }]);
  value = value || createPrompt.name;
  await product.create(value);
};

async function addApplication () {
  console.log(colors.yellow('Lets add an application to your product.'));
  let applicationPrompt = await inquirer.prompt([{
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
    name: applicationPrompt.name,
    path: applicationPrompt.path,
    repository: {
      type: applicationPrompt.repoType,
      url: applicationPrompt.repoUrl
    }
  };
  await product.addApplication(application);
};
