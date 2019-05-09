const { spawn } = require('child_process');
const colors = require('colors');

async function run (command, args = [], options) {
  console.log(colors.gray('Path: ', colors.green(options.cwd || process.cwd())));
  console.log(colors.gray('Command: ', colors.green(command)));
  await spawn(command, args, {
    stdio: 'inherit',
    shell: true,
    ...options
  });
}

function list(val) {
  if (val) {
    return val.split(',');
  }
}

module.exports = {
  run,
  list
};
