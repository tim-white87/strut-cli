const colors = require('colors');
const AWS = require('aws-sdk');
AWS.config.update({ region: 'us-east-1' });

exports.Project = class Project {
  constructor() {
    this.organizations = new AWS.Organizations();
  }

  async init() {
    this.organizationOuMap = await this.mapOrganizationOus();
    this.projectsOu = this.organizationOuMap[0].ous.find(ou => ou.Name === 'Projects');
    if (!this.projectsOu) {
      console.log(colors.yellow('No Projects OU detecting, generating...'));
      let projectsOu = await this.createOrganizationalUnit('Projects', this.organizationOuMap[0].root.Id);
      this.projectsOu = projectsOu.OrganizationalUnit;
    }
    this.projectsOuChildren = await this.mapChildren(this.projectsOu.Id);
  }

  async create(name) {
    let projectWithName = this.projectsOuChildren.find(ou => ou.Name === name);
    if (!projectWithName) {
      console.log(colors.green(`Creating ${name} Organizational Unit...`));
      await this.createOrganizationalUnit(name, this.projectsOu.Id);
    } else {
      console.log(colors.red(`'${name}' is already a project! Please try a unique value.`));
    }
    // await this.createAccounts();
    // await this.generateProjectRepo();
  }

  async destroy(name) {
    let projectToDestroy = this.projectsOuChildren.find(ou => ou.Name === name);
    if (!projectToDestroy) {
      console.log(colors.red(`'${name}' does not exist as a project OU`));
    } else {
      console.log(colors.red(`Destroying ${name} Organizational Unit...`));
      await this.deleteOrganizationalUnit(projectToDestroy.Id);
    }
  }

  async mapOrganizationOus() {
    let data = await this.listRoots();
    return Promise.all(data.Roots.map(async (root) => {
      let ous = await this.mapChildren(root.Id);
      return {
        root,
        ous
      };
    }));
  }

  async describeOrganization() {
    return new Promise((resolve, reject) => {
      this.organizations.describeOrganization((err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  async listRoots() {
    return new Promise((resolve, reject) => {
      this.organizations.listRoots((err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  async listChildren(options) {
    return new Promise((resolve, reject) => {
      this.organizations.listChildren(options, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  async mapChildren(ParentId) {
    let ousData = await this.listChildren({
      ChildType: 'ORGANIZATIONAL_UNIT',
      ParentId
    });
    return Promise.all(ousData.Children.map(async ou => {
      let ouData = await this.describeOrganizationalUnit(ou.Id);
      return ouData.OrganizationalUnit;
    }));
  }

  async describeOrganizationalUnit(OrganizationalUnitId) {
    return new Promise((resolve, reject) => {
      this.organizations.describeOrganizationalUnit({ OrganizationalUnitId }, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  async createOrganizationalUnit(Name, ParentId) {
    let params = { Name, ParentId };
    return new Promise((resolve, reject) => {
      this.organizations.createOrganizationalUnit(params, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  async deleteOrganizationalUnit(OrganizationalUnitId) {
    return new Promise((resolve, reject) => {
      this.organizations.deleteOrganizationalUnit({ OrganizationalUnitId }, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  createAccount(AccountName, Email) {
    return new Promise((resolve, reject) => {
      this.organizations.createAccount({ AccountName, Email }, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  generateProjectRepo() { }
};
