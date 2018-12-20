#!/bin/node
const program = require('commander');
const inquirer = require('inquirer');
const colors = require('colors');
const { Project } = require('./libs/project');

console.log(colors.blue('Welcome to Strut!'));

program
  .version(require('./package.json').version)
  .arguments('<cmd> [value]')
  .action(async (cmd, value) => {
    switch (cmd) {
      case 'init':
        value = value || await inquirer.prompt([{
          type: 'input',
          name: 'name',
          message: 'Enter project name:'
        }]);
        const project = new Project(value);
        await project.create();
        break;
      case 'provision':
        console.log('provision');
        break;
      case 'destroy':
        console.log('destroy');
        break;
      default:
        console.log(colors.red(`'${cmd}' command does not exist, try --help for valid commands`));
        break;
    }
  })
  .on('--help', function () {
    console.log('');
    console.log('Commands:');
    console.log('  init [name]    Initialize a project in an OU in AWS');
    console.log('  provision      Provisions the accounts in the project from the infrastructure.json');
    console.log('  destroy        Destroys the project and all associated accounts including the OU');
  })
  .parse(process.argv);
