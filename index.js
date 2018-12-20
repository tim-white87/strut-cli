#!/bin/node
const program = require('commander');
const colors = require('colors');

console.log(colors.blue('Welcome to Strut!'));

program
  .version(require('./package.json').version)
  .arguments('<cmd> [value]')
  .action((cmd, value) => {
    switch (cmd) {
      case 'init':
        console.log('init', value);
        break;
      case 'provision':
        console.log('provision');
        break;
      default:
        break;
    }
  })
  .on('--help', function () {
    console.log('');
    console.log('Commands:');
    console.log('  init [name]    Initialize a project in an OU in AWS');
    console.log('  provision      Provisions the accounts in the project from the infrastructure.json');
  })
  .parse(process.argv);

//   function list(val) {
//   return val.split(',');
// }
