const { Providers } = require('./providers');
const { AwsModel } = require('./awsModel');

let ProvidersMap = new Map([
  [Providers.AWS, AwsModel],
  [Providers.GCP, null]
]);

module.exports = {
  ProvidersMap
};
