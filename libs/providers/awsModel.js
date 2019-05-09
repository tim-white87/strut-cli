const { BaseProviderModel } = require('./baseProviderModel');
const { Providers } = require('./providers');
process.env.AWS_SDK_LOAD_CONFIG = true;
const CloudFormation = require('aws-sdk/clients/cloudformation');
const cloudformation = new CloudFormation();

class AwsModel extends BaseProviderModel {
  get providerName () { return Providers.AWS; }

  async init() {
    await super.init();
    await this.buildStacks();
  }

  async buildStacks () {
    await this.getStacks();

    // TODO implement cloudformation
    // TODO https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/CloudFormation.html
    let buildStackRequests = this.infrastructureData.map((resource, i) => {
      return new Promise(resolve => {
        let StackName = `${this.application.name}-${resource.name || i}-stack`;
        if (this.existingStacks.some(es => es.StackName === StackName)) {
          cloudformation.updateStack({
            StackName,
            Capabilities: ['CAPABILITY_NAMED_IAM'],
            TemplateBody: JSON.stringify(resource.fileData)
          }, (err, data) => {
            if (err) {
              console.log(err);
            } else {
              resolve(data);
            }
          });
        } else {
          cloudformation.createStack({
            StackName,
            Capabilities: ['CAPABILITY_NAMED_IAM'],
            TemplateBody: JSON.stringify(resource.fileData)
          }, (err, data) => {
            if (err) {
              console.log(err);
            } else {
              resolve(data);
            }
          });
        }
      });
    });

    let buildStacksResult = await Promise.all(buildStackRequests);
    console.log(buildStacksResult);
  }

  async getStacks () {
    let res = await new Promise(resolve => {
      cloudformation.describeStacks((err, data) => {
        if (err) {
          console.log(err);
        } else {
          resolve(data);
        }
      });
    });
    this.existingStacks = res.Stacks;
  }
};

module.exports = {
  AwsModel
};
