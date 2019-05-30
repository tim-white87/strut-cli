const colors = require('colors/safe');
const inquirer = require('inquirer');
const { run } = require('../utils');
const { BaseProviderModel } = require('./baseProviderModel');
process.env.AWS_SDK_LOAD_CONFIG = true;
const CloudFormation = require('aws-sdk/clients/cloudformation');
const cloudformation = new CloudFormation();
const S3 = require('aws-sdk/clients/s3');
const s3 = new S3();

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

const ResourceTypes = {
  CLOUDFRONT: 'AWS::CloudFront::Distribution',
  S3: 'AWS::S3::Bucket'
};

class AwsModel extends BaseProviderModel {
  get providerName () { return 'AWS'; } // TODO why can't i use Providers?

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
    await this.runPostProvisionCommands();
  }

  async destroy() {
    let stacks = await this.getStacks();
    if (stacks.length === 0) { return; }
    console.log(colors.red('DESTROYING STACKS!'));
    let stackResources = await this.getStackResources(stacks);
    await this.dumpBuckets(stackResources);
    await this.destroyStacks(stacks);
    await this.checkStackStatus();
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
    let buildStackRequests = this.infrastructureData.map((stack, i) => {
      return new Promise(resolve => {
        let StackName = stack.name || `${this.application.name}-${i + 1}-stack`;
        const params = {
          StackName,
          Capabilities: ['CAPABILITY_NAMED_IAM'],
          TemplateBody: stack.fileData
        };
        if (this.existingStacks.some(es => es.StackName === StackName)) {
          cloudformation.updateStack(params, (err, data) => {
            if (err) {
              console.log(err);
            } else {
              resolve(data);
            }
          });
        } else {
          cloudformation.createStack(params, (err, data) => {
            if (err) {
              console.log(err);
            } else {
              resolve(data);
            }
          });
        }
      });
    });
    // let stacks;
    // if (this.application.parallel) {
    //   stacks = await Promise.all(buildStackRequests);
    //   return this.checkStackStatus();
    // }
    for (let i = 0; i < buildStackRequests.length; i++) {
      let stackRequest = buildStackRequests[i];
      await stackRequest;
      await this.checkStackStatus();
    }
  }

  async destroyStacks(stacks) {
    const destroyStackRequests = stacks.map(stack => {
      return new Promise(resolve => {
        cloudformation.deleteStack({
          StackName: stack.StackName
        }, (err, data) => {
          if (err) {
            console.log(err);
          } else {
            resolve(data);
          }
        });
      });
    });
    await Promise.all(destroyStackRequests);
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
    return res.Stacks.filter(stack => stack.StackName.indexOf(this.application.name) > -1);
  }

  async getStackResources (stacks) {
    const reqs = stacks.map(stack => {
      return new Promise(resolve => {
        cloudformation.describeStackResources({
          StackName: stack.StackName
        }, (err, data) => {
          if (err) {
            console.log(err);
          } else {
            resolve(data);
          }
        });
      });
    });
    let res = await Promise.all(reqs);
    return res[0].StackResources;
  }

  async dumpBuckets (stackResources) {
    let buckets = stackResources.filter(resource => resource.ResourceType === ResourceTypes.S3);
    if (buckets.length === 0) {
      return;
    }
    let reqs = buckets.map(async bucket => {
      let objectsRes = await new Promise(resolve => {
        s3.listObjectsV2({ Bucket: bucket.PhysicalResourceId }, (err, data) => {
          if (err) {
            console.log(err);
          } else {
            resolve(data);
          }
        });
      });
      if (objectsRes.Contents.length === 0) {
        return;
      }

      console.log(colors.red(`WARNING: You have stuff in your S3 bucket: ${colors.gray(bucket.PhysicalResourceId)}`));
      console.log(colors.yellow('Deleting your S3 Resources will fail if these items are not removed, however strut strongly recommendeds you back these up.'));
      let really = await inquirer.prompt([{
        type: 'confirm',
        message: 'Should strut empty your buckets?',
        name: 'deleteObjects',
        default: false
      }]);
      if (!really.deleteObjects) {
        return;
      }

      return new Promise(resolve => {
        s3.deleteObjects({
          Bucket: bucket.PhysicalResourceId,
          Delete: {
            Objects: objectsRes.Contents.map(o => {
              return { Key: o.Key };
            })
          }
        }, (err, data) => {
          if (err) {
            console.log(err);
          } else {
            resolve(data);
          }
        });
      });
    });
    return Promise.all(reqs);
  }

  async checkStackStatus (isSilent) {
    let stacks = await this.getStacks();

    if (stacks.every(stack =>
      stack.StackStatus === StackStatus.CREATE_COMPLETE ||
      stack.StackStatus === StackStatus.UPDATE_COMPLETE ||
      stack.StackStatus === StackStatus.ROLLBACK_COMPLETE ||
      stack.StackStatus === StackStatus.UPDATE_ROLLBACK_COMPLETE ||
      stack.StackStatus === StackStatus.DELETE_COMPLETE)) {
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
