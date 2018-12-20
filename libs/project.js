const colors = require('colors');
const AWS = require('aws-sdk');
AWS.config.update({ region: 'us-east-1' });

exports.Project = class Project {
  constructor(name) {
    this.name = name;
    this.organizations = new AWS.Organizations();
  }

  async create() {
    this.organizationOuMap = await this.mapOrganizationOus();
    this.projectsOu = this.organizationOuMap[0].ous.find(ou => ou.Name === 'Projects');
    if (!this.projectsOu) {
      console.log(colors.yellow('No Projects OU detecting, generating...'));
      let projectsOu = await this.createOrganizationalUnit('Projects', this.organizationOuMap[0].root.Id);
      this.projectsOu = projectsOu.OrganizationalUnit;
    }
    console.log('projects OU', this.projectsOu);
    // await this.createOrganizationalUnit();
    // await this.createAccounts();
    // await this.generateProjectRepo();
  }

  async mapOrganizationOus() {
    let data = await this.listRoots();
    return Promise.all(data.Roots.map(async (root) => {
      let ousData = await this.listChildren({
        ChildType: 'ORGANIZATIONAL_UNIT',
        ParentId: root.Id
      });
      let ous = await Promise.all(ousData.Children.map(async ou => {
        let ouData = await this.describeOrganizationalUnit(ou.Id);
        return ouData.OrganizationalUnit;
      }));
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

  createAccounts() { }

  generateProjectRepo() { }
};
