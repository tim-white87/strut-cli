# Strut CLI

> Strut is an IaC (Infrastructure as Code) devops utility that orchestrates applications and cloud providers.

Want to run containers across multiple clouds? We got you. Need a way to orchestrate 10 microservice repos and a SPA front end on AWS? No problem. AWS prices got you in a bind? Simply add a new provider. We will get your whole flow setup using the tools native to the provider. Thinking about moving to the cloud? We can get you off that metal and floating faster than you ever thought possible.

Developers are empowered to define their application(s) stack in a way that makes sense for their team's product. The provider setup is designed to provide a cloud agnostic provisioning and deployment system that is consistent for supported providers. By enabling dynamic cloud providers plugins, we intend to support any cloud provider.

Our goal is for developers to focus on useful business code while still retaining the ability to define their infrastructure.

## WIP

Please be aware this is under active development and only supports AWS at the moment. Other providers will be added as time permits.


## Installation

MacOs
```
brew install strut-cli
```
TODO add other operating system installs

## Usage
```
strut create <app name>
```

### CLI Options:

  * -h, --help               output usage information

See above help for details

## Contributing

```
$ make install // From project clone location
```


## License

[ISC](https://github.com/cecotw/strut-cli/blob/master/LICENSE)