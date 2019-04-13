#!/bin/node
const program = require('commander');
const process = require('process');
const colors = require('colors');
const { Product } = require('./libs/product');
const createPrompt = require('./libs/prompts/createPrompt');
const addApplicationPrompt = require('./libs/prompts/addApplicationPrompt');
const product = new Product();

console.log(colors.blue('Welcome to Strut!'));

program
  .version(require('./package.json').version)
  .arguments('<cmd> [value]')
  .action(async (cmd, value) => {
    await product.init();
    switch (cmd) {
      case 'create':
        await createPrompt(product, value);
        process.chdir(`./${value}`);
        await addApplicationPrompt();
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
