
exports.Project = class Project {
  constructor(name) {
    this.name = name;
  }

  async create() {
    await this.createOrganizationalUnit();
    await this.createAccounts();
    await this.generateProjectRepo();
  }

  createOrganizationalUnit() { }

  createAccounts() { }

  generateProjectRepo() { }
};
