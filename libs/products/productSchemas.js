let product = {
  name: null,
  version: '0.0.0',
  applications: []
};

let localConifg = {
  path: null,
  commands: {
    install: [], // install commands
    validate: [], // validate commands (i.e. test, lint, etc.)
    build: [], // build commands
    start: [], // start commands
    deploy: [] // deploy commands
  },
  artifacts: [] // artifact paths
};

let application = {
  name: null,
  version: null,
  repository: {
    type: null,
    url: null
  },
  localConifg,
  providers: {}
};

let commands = {
  pre_provision: null, // runs before infrastructure provisioning
  post_provision: null // runs after infrastructure provisioning
};

let provider = {
  name: null,
  infrastructure: [],
  commands
};

let resource = {
  name: null,
  path: null,
  body: null,
  commands
};

module.exports = {
  product,
  application,
  provider,
  resource
};
