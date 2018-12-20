const inquirer = require('inquirer');

exports.Project = class Project {
  constructor(name) {
    this.name = name || this.promptForName();
  }

  async promptForName() {
    const prompt = await inquirer.prompt([{
      type: 'input',
      name: 'name',
      message: 'Enter project name:'
    }]);
    return prompt.name;
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
