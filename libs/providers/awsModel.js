const { BaseProviderModel } = require('./baseProviderModel');
const { Providers } = require('./providersModel');
process.env.AWS_SDK_LOAD_CONFIG = true;
const CloudFormation = require('aws-sdk/clients/cloudformation');
const cloudformation = new CloudFormation();

class AwsModel extends BaseProviderModel {
  get providerName () { return Providers.AWS; }

  async init() {
    super.init();

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
