
const { ProviderModel } = require('./baseProviderModel');

class AwsModel extends ProviderModel {
  async load() {
    console.log('providing aws');
  }
};

module.exports = {
  AwsModel
};
