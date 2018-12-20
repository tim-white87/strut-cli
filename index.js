#!/bin/node
// const spawn = require('child_process').spawnSync;
// const fs = require('fs');
const program = require('commander');
const colors = require('colors');

console.log(colors.blue('Welcome to Strut!'));

program
  .version('0.0.1', '-v, --version')
  .arguments('<cmd> [value]') // 'initialize an AWS Project'
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
    console.log('  init [name]    Inits a Project as an OU in AWS');
    console.log('  provision      Provisions the accounts in the project from the infrastructure.json');
  })
  .parse(process.argv);

//   function list(val) {
//   return val.split(',');
// }
