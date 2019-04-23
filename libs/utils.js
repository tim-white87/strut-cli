const { spawn } = require('child_process');
const colors = require('colors');

function run (command, args = [], options) {
  console.log(colors.gray('Path: ', colors.green(options.cwd || process.cwd())));
  console.log(colors.gray('Command: ', colors.green(command)));
  const child = spawn(command, args, {
    stdio: 'inherit',
    shell: true,
    ...options
  });

  child.on('error', function (err) {
    console.log(err);
  });

  child.stdout.on('data', (data) => {
    console.log(data);
  });

  child.stderr.on('data', (data) => {
    console.log(`stderr: ${data}`);
  });

  child.on('close', (code) => {
    console.log(`child process exited with code ${code}`);
  });
}

function list(val) {
  return val.split(',');
}

module.exports = {
  run,
  list
};
