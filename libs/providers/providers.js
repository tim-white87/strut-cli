const utils = require('../utils');
const { AwsModel } = require('./awsModel');

const Providers = {
  AWS: 'AWS',
  GCP: 'GCP'
};

const ProvidersMap = new Map([
  [Providers.AWS, AwsModel],
  [Providers.GCP, null]
]);

async function onProviderCommand (productModel, command, applications, providers) {
  if (applications) {
    applications = utils.list(applications).map(a => {
      return productModel.product.applications.find(app => app.name === a);
    });
  } else {
    applications = productModel.product.applications;
  }
  if (providers) {
    providers = utils.list(providers);
  }
  for (let i = 0; i < applications.length; i++) {
    let app = applications[i];
    for (let provider in app.providers) {
      if (!providers || providers.some(p => p === provider)) {
        if (app.providers[provider] && (app.providers[provider].commands || (app.providers[provider].infrastructure && app.providers[provider].infrastructure.length > 0))) {
          let Model = ProvidersMap.get(provider);
          let model = new Model(app);
          await model.init();
          await model[command]();
        }
      }
    }
  }
}

module.exports = {
  Providers,
  ProvidersMap,
  onProviderCommand
};
