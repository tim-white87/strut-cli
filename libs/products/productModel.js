const fs = require('fs');
const colors = require('colors');
const productSchemas = require('./productSchemas');

exports.ProductModel = class ProductModel {
  async init() {
    this.product = await this.loadProduct();
  }

  async loadProduct () {
    const productJson = './strut/product.json';
    return new Promise(resolve => {
      fs.stat(productJson, err => {
        if (err && err.code === 'ENOENT') {
          resolve(this.productModel);
        } else {
          fs.readFile(productJson, 'utf8', (err, data) => {
            if (err) throw err;
            resolve(JSON.parse(data));
          });
        }
      });
    });
  }

  updateProductFile (data) {
    const productJsonPath = './strut/product.json';
    return new Promise((resolve, reject) => {
      fs.writeFile(
        productJsonPath,
        JSON.stringify(data || this.product, null, 2),
        (err) => {
          if (err) {
            reject(err);
          }
          resolve();
        }
      );
    });
  }

  async create (name) {
    this.name = name || process.cwd().split('/').pop();
    console.log(colors.yellow(`Creating product: ${colors.gray(this.name)}`));
    let dir = `./strut`;
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir);
    }
    this.product = { ...productSchemas.product, name };
    await this.updateProductFile(this.product);
    console.log(colors.green('DONE!'));
  }

  async addApplication (application) {
    console.log(colors.yellow(`Adding application to product: ${colors.gray(application.name)}`));
    if (!this.product.applications.some(a => a.name === application.name)) {
      this.product.applications.push(application);
      await this.updateProductFile();
      console.log(colors.green('DONE!'));
    } else {
      console.log(colors.red('A name with this application already exists.'));
    }
  }

  async updateApplication (application) {
    console.log(colors.yellow(`Updating ${colors.gray(application.name)}: ${colors.gray(application.name)}`));
    let existingApplication = this.product.applications.find(a => a.name === application.name);
    if (existingApplication) {
      existingApplication = { ...existingApplication, ...application };
      await this.updateProductFile();
    } else {
      console.log(colors.red('No application with this name exists'));
    }
  }
};
