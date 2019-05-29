const colors = require('colors/safe');
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

  async init () {
    await super.init();
  }

  async provision() {
    this.stacks = null;
    // TODO await this.runPreProvisonCommands();
    console.log(colors.gray('Checking current CF stacks status...'));
    let isIdle = await this.checkStackStatus();
    if (isIdle) {
      console.log(colors.gray('Running CF stacks...'));
      await this.buildStacks();
    } else {
      return;
    }
    isIdle = await this.checkStackStatus();
    if (isIdle) {
      await this.runPostProvisionCommands();
    }
  }

  async runPostProvisionCommands () {
    if (this.provider.commands && this.provider.commands.post_provision) {
      console.log(colors.gray('Running post provision commands...'));
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
            TemplateBody: resource.fileData
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
            TemplateBody: resource.fileData
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

  async checkStackStatus (isSilent) {
    let stackRes = await this.getStacks();
    let stacks = stackRes.filter(stack => stack.StackName.indexOf(this.application.name) > -1);

    if (stacks.every(stack =>
      stack.StackStatus === StackStatus.CREATE_COMPLETE ||
      stack.StackStatus === StackStatus.UPDATE_COMPLETE ||
      stack.StackStatus === StackStatus.ROLLBACK_COMPLETE ||
      stack.StackStatus === StackStatus.UPDATE_ROLLBACK_COMPLETE)) {
      if (!isSilent) {
        console.log(colors.green('Stacks status:'));
        stacks.forEach(stack => {
          console.log(`${colors.white(stack.StackName)}: ${colors.green(stack.StackStatus)}`);
        });
      }
      return true;
    } else if (stacks.some(stack => stack.StackStatus.indexOf('PROGRESS') > -1)) {
      let isUpdated = !this.stacks || !this.stacks.every(stack => stacks.find(s => s.StackName === stack.StackName).StackStatus === stack.StackStatus);
      if (isUpdated) {
        this.stacks = stacks;
        if (!isSilent) {
          stacks.forEach(stack => {
            console.log(`${colors.white(stack.StackName)}: ${colors.yellow(stack.StackStatus)}`);
            if (stack.StackStatusReason) {
              console.log(colors.gray(stack.StackStatusReason));
            }
          });
        }
      }

      await new Promise(resolve => {
        setTimeout(() => {
          resolve(this.checkStackStatus(isSilent));
        }, 1000);
      });
    } else {
      if (!isSilent) {
        console.log(colors.red('Stacks status:'));
        stacks.forEach(stack => {
          console.log(`${colors.white(stack.StackName)}: ${colors.red(stack.StackStatus)}`);
        });
      }
      return false;
    }
  }
};

module.exports = {
  AwsStackStatus: StackStatus,
  AwsModel
};
