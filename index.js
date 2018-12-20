#!/bin/node
// const spawn = require('child_process').spawnSync;
// const fs = require('fs');
const program = require('commander');
const inquirer = require('inquirer');
const colors = require('colors');

function list(val) {
  return val.split(',');
}

program
  .version('0.0.1', '-v, --version')
  .command('init [name]', 'initialize an AWS Project')
  .command('provision', 'Provisions the project\'s accounts with the defined infrastructure')
  .parse(process.argv);

const questions = [{
  // type: 'list',
  // name: 'tool',
  // message: 'What tool would you like to run?',
  // choices: tools.toolList.map(tool => Object.assign({ name: tool.name, value: tool.name })),
  // when () {
  //   return program._.length === 0;
  // }
}];

async function main () {
  console.log(colors.blue('Welcome to Strut!'));
  // const answers = await inquirer.prompt(questions);
  // const selectedTool = program._[0] || answers.tool;
  // const tool = tools.toolList.find(t => t.name === selectedTool);
  // tools.runTool(tool, project.project);
}

main();
