const CloudFormation = require('aws-sdk/clients/cloudformation');
const { ProviderModel } = require('./baseProviderModel');
const cloudformation = new CloudFormation();

class AwsModel extends ProviderModel {
  async load() {
    console.log('providing aws');
    console.log(this.application);
  }
};

module.exports = {
  AwsModel
};
