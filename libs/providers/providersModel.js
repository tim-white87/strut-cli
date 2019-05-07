const { AwsModel } = require('./awsModel');

let providerMap = new Map([
  ['AWS', AwsModel],
  ['GCP', null]
]);

module.exports = {
  providerMap
};
