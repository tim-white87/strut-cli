#!/bin/node
const program = require('commander');
const process = require('process');
const fs = require('fs');
const colors = require('colors');
const utils = require('./libs/utils');
const { ProductModel } = require('./libs/products/productModel');
const { ProvidersMap } = require('./libs/providers/providersMap');
const createPrompt = require('./libs/prompts/createPrompt');
const addApplicationPrompt = require('./libs/prompts/addApplicationPrompt');
const addProviderPrompt = require('./libs/prompts/addProviderPrompt');

const productModel = new ProductModel();

console.log(colors.blue('Welcome to Strut!'));

async function main () {
  const strutDir = 'strut';
  if (fs.existsSync(strutDir)) {
    process.chdir(strutDir);
  }

  await productModel.init();

  program.version(require('./package.json').version);

  program
    .command('create [name]')
    .description('Create a new strut product')
    .action(async (name) => {
      await createPrompt(productModel, name);
      await addApplicationPrompt(productModel);
      await addProviderPrompt(productModel);
    });

  program
    .command('add <type> [name] [value]')
    .description('Add an <application|provider> to the product')
    .action(async (type, name, value) => {
      switch (type) {
        case 'application':
          await addApplicationPrompt(productModel, name);
          break;
        case 'provider':
          await addProviderPrompt(productModel, name, value);
          // TODO add various provider IaC to setup cloud CI/CD
          break;
        default:
          console.log(colors.red(`'${type}' is not a valid type, try --help for valid commands`));
          break;
      }
    });

  program
    .command('run <cmd> [applications]')
    .description('Runs the specified command <install|build|start> for the product applications (separated with a comma). Default will run all apps.')
    .action((cmd, applications) => {
      console.log(colors.blue(`run: ${colors.gray(cmd)}`));
      if (applications) {
        applications = utils.list(applications).map(a => {
          return productModel.product.applications.find(app => app.name === a);
        });
      } else {
        applications = productModel.product.applications;
      }
      applications.forEach(app => {
        if (app.commands[cmd] && app.commands[cmd].length > 0) {
          utils.run(app.commands[cmd].join(' '), [], { cwd: app.path });
        } else {
          console.log(colors.yellow(`No ${colors.gray(cmd)} command defined for ${colors.gray(app.name)}`));
        }
      });
    });

  program
    .command('deploy [applications] [providers]')
    .description('Deploys the applications to the specified provider. Defaults to all applications deployed to all providers')
    .action((applications, providers) => {
      if (applications) {
        applications = utils.list(applications).map(a => {
          return productModel.product.applications.find(app => app.name === a);
        });
      } else {
        applications = productModel.product.applications;
      }
      if (providers) {
        providers = utils.list(providers);
      }
      applications.forEach(async app => {
        for (let provider in app.providers) {
          if (!providers || providers.some(p => p === provider)) {
            if (app.providers[provider] && app.providers[provider].infrastructure &&
              app.providers[provider].infrastructure.length > 0) {
              let Model = ProvidersMap.get(provider);
              let model = new Model(app);
              await model.init();
            }
          }
        }
      });
    });

  program.parse(process.argv);
}

main();
