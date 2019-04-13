#!/bin/node
const program = require('commander');
const inquirer = require('inquirer');
const colors = require('colors');
const { Product } = require('./libs/product');

console.log(colors.blue('Welcome to Strut!'));

program
  .version(require('./package.json').version)
  .arguments('<cmd> [value]')
  .action(async (cmd, value) => {
    const product = new Product();

    switch (cmd) {
      case 'create':
        let create = await inquirer.prompt([{
          type: 'input',
          name: 'name',
          message: 'Enter new strut product:',
          when () { return !value; }
        }]);
        value = value || create.name;
        await product.create(value);
        break;
      // case 'provision':
      //   console.log('provision');
      //   break;
      // case 'destroy':
      //   let destroy = await inquirer.prompt([{
      //     type: 'list',
      //     name: 'name',
      //     message: 'Select project to destroy:',
      //     choices: project.projectsOuChildren.map(project => Object.assign({ name: project.Name, value: project.Name })),
      //     when () { return !value; }
      //   }, {
      //     type: 'confirm',
      //     name: 'forsure',
      //     default: false,
      //     message: colors.red('Are you sure you want to destroy this?')
      //   }]);
      //   value = value || destroy.name;
      //   if (destroy.forsure) {
      //     await project.destroy(value);
      //   }
      //   break;
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
