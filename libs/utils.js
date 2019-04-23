const { spawn } = require('child_process');
const colors = require('colors');

function run (app, args = []) {
  let command = app.start.join(' ');
  console.log(colors.gray(`${colors.green(app.name)}: ${command}`));
  const child = spawn(command, args, { cwd: app.path, stdio: 'inherit', shell: true });

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
