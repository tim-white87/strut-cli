const path = require('path');
const { readFile } = require('../utils');
const { ProviderModel } = require('./baseProviderModel');
process.env.AWS_SDK_LOAD_CONFIG = true;
const CloudFormation = require('aws-sdk/clients/cloudformation');
const cloudformation = new CloudFormation();

class AwsModel extends ProviderModel {
  async init() {
    // TODO abstract this and put in base model
    this.providerName = this.application.providers.AWS.name || 'AWS';
    this.infrastructure = this.application.providers.AWS.infrastructure;
    this.infrastructureFiles = await Promise.all(this.infrastructure.map(
      resource => {
        return readFile(path.join(this.application.path, resource.path));
      }));
    this.infrastructureData = this.infrastructure.map((resource, i) => {
      return { ...resource, fileData: this.infrastructureFiles[i] };
    });
    console.log(this.infrastructureData);
    cloudformation.describeStacks((err, data) => {
      console.log(err, data);
    });
    // TODO implement cloudformation
    // TODO https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/CloudFormation.html
  }
};

module.exports = {
  AwsModel
};
