const fs = require('fs');
const colors = require('colors');

exports.Product = class Product {
  get productModel() {
    return {
      name: null,
      version: '0.0.0',
      applications: [],
      providers: []
    };
  };

  get applicationModel () {
    return {
      name: null,
      repository: {
        type: null,
        url: null
      },
      path: null, // optionally point at local folder for app
      install: [], // install commands
      build: [], // build commands
      artifacts: [] // artifact paths
    };
  };

  async create (name) {
    console.log(colors.yellow(`Creating product: ${colors.gray(name)}`));
    this.name = name || 'myproduct';
    this.dir = `./${name}`;
    if (!fs.existsSync(this.dir)) {
      fs.mkdirSync(this.dir);
    }
    await this.updateProductFile({ ...this.productModel, name });
    console.log(colors.green('DONE!'));
  }

  async updateProductFile (data) {
    let stream = new Promise((resolve, reject) => {
      fs.writeFile(
        `${this.dir}/product.json`,
        JSON.stringify(data, null, 2),
        (err) => {
          if (err) {
            reject(err);
          }
          resolve();
        }
      );
    });
    return stream;
  }
};
