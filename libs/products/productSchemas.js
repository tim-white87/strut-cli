let product = {
  name: null,
  version: '0.0.0',
  applications: [],
  providers: []
};

let application = {
  name: null,
  repository: {
    type: null,
    url: null
  },
  install: [], // install commands
  build: [], // build commands
  artifacts: [], // artifact paths,
  start: [] // start the application
};

let provider = {
  name: null,
  applications: []
};

let providerApplication = {
  name: null,
  resources: []
};

module.exports = {
  product,
  application,
  provider,
  providerApplication
};
