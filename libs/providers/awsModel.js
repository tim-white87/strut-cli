const colors = require('colors');
const { run } = require('../utils');
const { BaseProviderModel } = require('./baseProviderModel');
const { Providers } = require('./providers');
process.env.AWS_SDK_LOAD_CONFIG = true;
const CloudFormation = require('aws-sdk/clients/cloudformation');
const cloudformation = new CloudFormation();

const StackStatus = {
  CREATE_IN_PROGRESS: 'CREATE_IN_PROGRESS',
  CREATE_FAILED: 'CREATE_FAILED',
  CREATE_COMPLETE: 'CREATE_COMPLETE',
  ROLLBACK_IN_PROGRESS: 'ROLLBACK_IN_PROGRESS',
  ROLLBACK_FAILED: 'ROLLBACK_FAILED',
  ROLLBACK_COMPLETE: 'ROLLBACK_COMPLETE',
  DELETE_IN_PROGRESS: 'DELETE_IN_PROGRESS',
  DELETE_FAILED: 'DELETE_FAILED',
  DELETE_COMPLETE: 'DELETE_COMPLETE',
  UPDATE_IN_PROGRESS: 'UPDATE_IN_PROGRESS',
  UPDATE_COMPLETE_CLEANUP_IN_PROGRESS: 'UPDATE_COMPLETE_CLEANUP_IN_PROGRESS',
  UPDATE_COMPLETE: 'UPDATE_COMPLETE',
  UPDATE_ROLLBACK_IN_PROGRESS: 'UPDATE_ROLLBACK_IN_PROGRESS',
  UPDATE_ROLLBACK_FAILED: 'UPDATE_ROLLBACK_FAILED',
  UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS: 'UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS',
  UPDATE_ROLLBACK_COMPLETE: 'UPDATE_ROLLBACK_COMPLETE',
  REVIEW_IN_PROGRESS: 'REVIEW_IN_PROGRESS'
};

class AwsModel extends BaseProviderModel {
  get providerName () { return Providers.AWS; }

  async init() {
    await super.init();
    // TODO await this.runPreProvisonCommands();
    await this.getStacks();
    await this.buildStacks();
    let success = await this.checkStackStatus();
    if (success) {
      await this.runPostProvisionCommands();
    }
  }

  async runPostProvisionCommands () {
    if (this.provider.commands && this.provider.commands.post_provision) {
      return run(this.provider.commands.post_provision.join(' '), [], { cwd: this.application.path });
    }
  }

  /**
   *  https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/CloudFormation.html
   */
  async buildStacks () {
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

    await Promise.all(buildStackRequests);
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
    return res.Stacks;
  }

  async checkStackStatus () {
    let stacks = await this.getStacks();
    stacks.forEach(stack => {
      console.log(`${colors.green(stack.StackName)}: ${colors.yellow(stack.StackStatus)}`);
    });
    if (stacks.some(stack => stack.StackStatus.indexOf('PROGRESS') > -1)) {
      await new Promise(resolve => {
        setTimeout(() => {
          resolve(this.checkStackStatus());
        }, 1000);
      });
    } else if (stacks.every(stack => stack.StackStatus.indexOf('COMPLETE') > -1)) {
      return true;
    } else {
      return false;
    }
  }
};

module.exports = {
  AwsStackStatus: StackStatus,
  AwsModel
};
