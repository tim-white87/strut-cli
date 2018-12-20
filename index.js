#!/bin/node
// const spawn = require('child_process').spawnSync;
// const fs = require('fs');
const program = require('commander');
const colors = require('colors');

function list(val) {
  return val.split(',');
}

program
  .version('0.0.1', '-v, --version')
  .arguments('<cmd>') // 'initialize an AWS Project'
  .action(cmd => {
    console.log(cmd);
  })
  // .command('provision', 'Provisions the project\'s accounts with the defined infrastructure')
  .parse(process.argv);

console.log(colors.blue('Welcome to Strut!'));

