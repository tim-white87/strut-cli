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
    await this.runCommands('pre_provision');
    console.log(colors.gray('Checking current CF stacks status...'));
    let isIdle = await this.checkStackStatus();
    if (this.rollbackCompleteStacks) {
      await this.removeRollbackCompleteStacks();
      isIdle = await this.checkStackStatus();
    }
    if (isIdle) {
      console.log(colors.gray('Running CF stacks...'));
      try {
        await this.buildStacks();
        // TODO add Route53 record set if stacks contained CertificateManager
      } catch (e) {
        return;
      }
    } else {
      return;
    }
    await this.checkStackStatus();
    await this.runCommands('post_provision');
  }

  async destroy() {
    let stacks = await this.getStacks();
    if (stacks.length === 0) { return; }
    console.log(colors.red('DESTROYING STACKS!'));
    let stacksResources = await this.getStacksResources(stacks);
    await this.dumpBuckets(stacksResources);
    await this.deleteStacks(stacks);
    await this.checkStackStatus();
  }

  async runCommands (type) {
    if (this.provider.commands && this.provider.commands[type]) {
      console.log(colors.gray(`Running ${colors.green(type)} commands...`));
      return run(this.provider.commands[type].join(' '), [], { cwd: this.application.path });
    }
  }

  /**
   *  https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/CloudFormation.html
   */
  async buildStacks () {
    let buildStackRequests = this.infrastructureData.map((stack, i) => {
      return new Promise((resolve, reject) => {
        let StackName = stack.name || `${this.application.name}-${i + 1}-stack`;
        const params = {
          StackName,
          Capabilities: ['CAPABILITY_NAMED_IAM'],
          TemplateBody: stack.fileData,
          Parameters: this.params
        };
        if (this.existingStacks.some(es => es.StackName === StackName)) {
          cloudformation.updateStack(params, (err, data) => {
            if (err) {
              console.log(colors.red(`UPDATE ERROR: ${colors.white(StackName)}`));
              console.log(colors.gray('Message: '), colors.yellow(err.message));
              reject(err);
            } else {
              resolve(data);
            }
          });
        } else {
          cloudformation.createStack(params, (err, data) => {
            if (err) {
              console.log(colors.red(`CREATE ERROR: ${colors.white(StackName)}`));
              console.log(colors.gray('Message: '), colors.yellow(err.message));
            } else {
              resolve(data);
            }
          });
        }
      });
    });
    for (let i = 0; i < buildStackRequests.length; i++) {
      let stackRequest = buildStackRequests[i];
      await stackRequest;
      await this.checkStackStatus();
    }
  }

  async deleteStacks(stacks) {
    const reqs = stacks.map(stack => {
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
    await Promise.all(reqs);
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

  async getStacksResources (stacks) {
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
    return res.map(s => s.StackResources);
  }

  async getStackEvents(stack, nextToken) {
    return new Promise(resolve => {
      let params = {
        StackName: stack.StackName
      };
      if (nextToken) {
        params.NextToken = nextToken;
      }
      cloudformation.describeStackEvents(params, (err, data) => {
        if (err) {
          console.log(err);
        } else {
          resolve(data);
        }
      });
    });
  }

  async dumpBuckets (stacksResources) {
    let reqs = [];
    for (let i = 0; i < stacksResources.length; i++) {
      let stackResources = stacksResources[i];
      let buckets = stackResources.filter(resource => resource.ResourceType === ResourceTypes.S3);
      if (buckets.length === 0) {
        return;
      }
      let bucketReqs = buckets.map(async bucket => {
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
      reqs = reqs.concat(bucketReqs);
    }

    return Promise.all(reqs);
  }

  async checkStackStatus (isSilent) {
    let stacks = await this.getStacks();

    if (stacks.every(stack =>
      stack.StackStatus === StackStatus.CREATE_COMPLETE ||
      stack.StackStatus === StackStatus.UPDATE_COMPLETE ||
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
        this.rollbackCompleteStacks = stacks.filter(stack => stack.StackStatus === StackStatus.ROLLBACK_COMPLETE);
      }
      return false;
    }
  }

  async removeRollbackCompleteStacks () {
    console.log(colors.yellow(`It looks like some stacks are in ${colors.red(StackStatus.ROLLBACK_COMPLETE)} state and can not be updated.`));
    console.log(colors.yellow('You need to delete these in order to update them.'));
    console.log(colors.gray('Stacks: '), this.rollbackCompleteStacks.map(stack => stack.StackName).join(', '));
    for (let i = 0; i < this.rollbackCompleteStacks.length; i++) {
      let stack = this.rollbackCompleteStacks[i];
      let stackEventsRes = await this.getStackEvents(stack);
      stackEventsRes.StackEvents.filter(event => event.ResourceStatus === StackStatus.CREATE_FAILED)
        .forEach(event => {
          console.log(`${colors.gray(event.Timestamp)}: ${colors.green(event.LogicalResourceId)}`);
          console.log(`${colors.red(event.ResourceStatus)}: ${colors.yellow(event.ResourceStatusReason)}`);
        });
      // TODO call this again recursivley if stackEventsRes.NextToken
    }
    const really = await inquirer.prompt([{
      type: 'confirm',
      message: 'Would you like strut to remove these stacks?',
      default: false,
      name: 'deleteStacks'
    }]);
    if (really.deleteStacks) {
      await this.deleteStacks(this.rollbackCompleteStacks);
      this.rollbackCompleteStacks = null;
      return this.checkStackStatus();
    }
  }
};

module.exports = {
  AwsStackStatus: StackStatus,
  AwsModel
};
