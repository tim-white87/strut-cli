const { spawn } = require('child_process');
const colors = require('colors');
const fs = require('fs');

async function run (command, args = [], options) {
  console.log(colors.gray('Path: ', colors.green(options.cwd || process.cwd())));
  console.log(colors.gray('Command: ', colors.green(command)));
  let cmd = await spawn(command, args, {
    stdio: 'inherit',
    shell: true,
    ...options
  });
  return new Promise(resolve => {
    cmd.on('close', code => {
      resolve(code);
    });
  });
}

function list(val) {
  if (val) {
    return val.split(',');
  }
}

function readFile (path) {
  return new Promise((resolve, reject) => {
    fs.stat(path, err => {
      if (err && err.code === 'ENOENT') {
        reject(err);
      } else {
        fs.readFile(path, 'utf8', (err, data) => {
          if (err) throw err;
          try {
            resolve(data);
          } catch (e) {
            console.log(e);
          }
        });
      }
    });
  });
}

module.exports = {
  run,
  list,
  readFile
};
