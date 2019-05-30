const path = require('path');
const { readFile } = require('../utils');

class BaseProviderModel {
  get providerName () {}
  get provider () { return this.application.providers[this.providerName]; }
  get infrastructure () { return this.provider.infrastructure; }

  constructor(application, rawParams) {
    this.application = application;
    this.params = this.parseRawParams(rawParams);
  }

  async init() {
    let infrastructureFiles = await Promise.all(this.infrastructure.map(
      resource => {
        return readFile(path.join(this.application.path || './', resource.path));
      }));
    this.infrastructureData = this.infrastructure.map((resource, i) => {
      return { ...resource, fileData: infrastructureFiles[i] };
    });
  }

  parseRawParams(rawParams) {
    if (rawParams) {
      return rawParams.map(param => {
        let key = param.split(':')[0];
        let value = param.split(':').pop();
        return { ParameterKey: key, ParameterValue: value };
      });
    }
  }
}

module.exports = {
  BaseProviderModel
};
