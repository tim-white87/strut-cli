#!/bin/node
const spawn = require('child_process').spawnSync;
const fs = require('fs');
const program = require('commander');
const colors = require('colors');

function list(val) {
  return val.split(',');
}

program
  .parse(process.argv);

console.log('test');