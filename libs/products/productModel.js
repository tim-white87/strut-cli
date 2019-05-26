const fs = require('fs');
const yaml = require('js-yaml');
const colors = require('colors');
const productSchemas = require('./productSchemas');

exports.ProductModel = class ProductModel {
  get fileName () { return 'strut'; }

  async init() {
    this.product = this.loadStrut();
  }

  loadStrut () {
    const yamlFile = `${this.fileName}.yml`;
    const jsonFile = `${this.fileName}.json`;
    let strutYaml;
    let strutJson;
    try {
      strutYaml = yaml.safeLoad(fs.readFileSync(yamlFile, 'utf8'));
    } catch (e) {}
    if (strutYaml) {
      this.isYaml = true;
      return strutYaml;
    }
    try {
      strutJson = JSON.parse(fs.readFileSync(jsonFile, 'utf8'));
    } catch (e) {}
    if (strutJson) {
      this.isYaml = false;
      return strutJson;
    } else {
      return this.productModel;
    }
  }

  updateProductFile (data) {
    return new Promise((resolve, reject) => {
      fs.writeFile(
        `${this.fileName}.${this.isYaml ? 'yml' : 'json'}`,
        this.isYaml ? yaml.safeDump(data || this.product) : JSON.stringify(data || this.product, null, 2),
        (err) => {
          if (err) {
            reject(err);
          }
          resolve();
        }
      );
    });
  }

  async create (name, isYaml) {
    this.name = name || process.cwd().split('/').pop();
    this.isYaml = isYaml;
    console.log(colors.yellow(`Creating product: ${colors.gray(this.name)}`));
    let dir = `./`;
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
};
