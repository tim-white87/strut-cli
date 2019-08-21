#!/bin/node
const program = require('commander');
const process = require('process');
const fs = require('fs');
const colors = require('colors');
const utils = require('./libs/utils');
const { ProductModel } = require('./libs/products/productModel');
const { onProviderCommand } = require('./libs/providers/providers');
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

  program
    .version(require('./package.json').version)
    .option('-c, --cloudformationParams <params>', 'Params to pass to Cloudformation in Key:Value format separated by commas', (val) => {
      if (val.indexOf(',') > -1) {
        return val.split(',');
      }
      return [val];
    });

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
          break;
        default:
          console.log(colors.red(`'${type}' is not a valid type, try --help for valid commands`));
          break;
      }
    });

  program
    .command('run <cmd> [applications]')
    .description('Runs the specified command <install|build|start> for the product applications (separated with a comma). Default will run all apps.')
    .action(async (cmd, applications) => {
      console.log(colors.blue(`run: ${colors.gray(cmd)}`));
      if (applications) {
        applications = utils.list(applications).map(a => {
          return productModel.product.applications.find(app => app.name === a);
        });
      } else {
        applications = productModel.product.applications;
      }
      for (let i = 0; i < applications.length; i++) {
        let app = applications[i];
        if (app.localConfig && app.localConfig.commands &&
        app.localConfig.commands[cmd] && app.localConfig.commands[cmd].length > 0) {
          for (let j = 0; j < app.localConfig.commands[cmd].length; j++) {
            utils.run(app.localConfig.commands[cmd][j], [], { cwd: app.path });
          }
        } else {
          console.log(colors.yellow(`No local ${colors.gray(cmd)} command defined for ${colors.gray(app.name)}`));
        }
      };
    });

  program
    .command('provision [applications] [providers]')
    .description('Provisions the defined infrastructure for the applications to the specified provider. Defaults to all applications deployed to all providers')
    .action(async (applications, providers) => {
      await onProviderCommand(productModel, 'provision', applications, providers, program.cloudformationParams);
    });

  program
    .command('destroy [applications] [providers]')
    .description('Destroys the defined infrastructure for the applications to the specified provider. Defaults to all applications destroyed for all providers. Careful with this one dude, it will kill your shit.')
    .action(async (applications, providers) => {
      productModel.product.applications.reverse();
      await onProviderCommand(productModel, 'destroy', applications, providers, program.cloudformationParams);
    });

  program
    .command('clone [applications]')
    .description('Clones the specified applications to their respective local config paths')
    .action(applications => {
      productModel.product.applications.forEach(async app => {
        await utils.run(`git clone ${app.repository.url}`);
      });
    });

  program.parse(process.argv);
}

main();
