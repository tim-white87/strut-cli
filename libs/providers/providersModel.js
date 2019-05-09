const { AwsModel } = require('./awsModel');

const Providers = {
  AWS: 'AWS',
  GCP: 'GCP'
};

let providerMap = new Map([
  [Providers.AWS, AwsModel],
  [Providers.GCP, null]
]);

module.exports = {
  providers,
  providerMap
};
