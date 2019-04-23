#!/bin/node
const program = require('commander');
const process = require('process');
const colors = require('colors');
const utils = require('./libs/utils');
const { ProductModel } = require('./libs/products/productModel');
const createPrompt = require('./libs/prompts/createPrompt');
const addApplicationPrompt = require('./libs/prompts/addApplicationPrompt');
const addProviderPrompt = require('./libs/prompts/addProviderPrompt');
const linkPrompt = require('./libs/prompts/linkPrompt');

const productModel = new ProductModel();

console.log(colors.blue('Welcome to Strut!'));

async function main () {
  await productModel.init();

  program.version(require('./package.json').version);

  program
    .command('create [name]')
    .description('Create a new product')
    .action(async (name) => {
      await createPrompt(productModel, name);
      process.chdir(`./${productModel.name}`);
      await addApplicationPrompt(productModel);
      let provider = await addProviderPrompt(productModel);
      await linkPrompt(productModel, null, provider.name);
    });

  program
    .command('add <type> [value]')
    .description('Add an <application|provider> to the product')
    .action(async (type, value) => {
      switch (type) {
        case 'application':
          await addApplicationPrompt(productModel);
          break;
        case 'provider':
          let provider = await addProviderPrompt(productModel, value);
          await linkPrompt(productModel, null, provider.name);
          // TODO add various provider IaC to setup cloud CI/CD
          break;
        default:
          console.log(colors.red(`'${type}' is not a valid type, try --help for valid commands`));
          break;
      }
    });

  program
    .command('link <application-name> <provider-name>')
    .description('Links an application to a provider')
    .action(async (applicationName, providerName) => {
      await linkPrompt(productModel, applicationName, providerName);
    });

  program
    .command('run <cmd> [applications]')
    .description('Runs the specified command for the product applications (separated with a comma). Default will run all apps.')
    .action((cmd, applications) => {
      if (applications) {
        applications = utils.list(applications).map(a => {
          return productModel.product.applications.find(app => app.name === a);
        });
      } else {
        applications = productModel.product.applications;
      }
      applications.forEach(app => {
        if (app[cmd]) {
          utils.run(app[cmd].join(' '), [], { cwd: app.path });
        }
      });
    });

  program.parse(process.argv);
}

main();
