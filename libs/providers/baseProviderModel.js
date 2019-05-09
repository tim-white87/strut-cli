const path = require('path');
const { readFile } = require('../utils');

class BaseProviderModel {
  get providerName () {}
  get provider () { return this.application.providers[this.providerName]; }
  get infrastructure () { return this.provider.infrastructure; }

  constructor(application) {
    this.application = application;
  }

  async init() {
    let infrastructureFiles = await Promise.all(this.infrastructure.map(
      resource => {
        return readFile(path.join(this.application.path, resource.path));
      }));
    this.infrastructureData = this.infrastructure.map((resource, i) => {
      return { ...resource, fileData: infrastructureFiles[i] };
    });
  }
}

module.exports = {
  BaseProviderModel
};
