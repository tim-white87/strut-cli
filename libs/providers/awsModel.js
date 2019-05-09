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

    console.log(this.existingStacks);

    // TODO implement cloudformation
    // TODO https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/CloudFormation.html
    // let createRequests = this.infrastructureData.map((resource, i) => {
    //   return new Promise(resolve => {
    //     cloudformation.updateStack({
    //       StackName: `${this.application.name}-${resource.name || i}-stack`,
    //       TemplateBody: JSON.stringify(resource.fileData)
    //     }, (err, data) => {
    //       if (err) {
    //         console.log(err);
    //       } else {
    //         resolve(data);
    //       }
    //     });
    //   });
    // });

    // let createResult = await Promise.all(createRequests);
    // console.log(createResult);
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
