const fs = require('fs');

exports.Product = class Product {
  get productModel() {
    return {
      name: null,
      version: '0.0.0',
      applications: [],
      providers: []
    };
  };

  async create (name) {
    this.name = name || 'myproduct';
    this.dir = `./${name}`;
    if (!fs.existsSync(this.dir)) {
      fs.mkdirSync(this.dir);
    }
    await this.updateProductFile({ ...this.productModel, name });
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
