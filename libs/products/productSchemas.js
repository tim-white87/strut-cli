let product = {
  name: null,
  version: '0.0.0',
  applications: []
};

let application = {
  name: null,
  path: null,
  repository: {
    type: null,
    url: null
  },
  commands: {
    install: [], // install commands
    validate: [], // validate commands (i.e. test, lint, etc.)
    build: [], // build commands
    start: [] // start commands
  },
  artifacts: [], // artifact paths
  providers: {}
};

let provider = {
  name: null,
  infrastructure: []
};

let resource = {
  path: null
};

module.exports = {
  product,
  application,
  provider,
  resource
};
