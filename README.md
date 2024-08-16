# index-sync

This AWS Lambda designed to update an index from one resource and send it to a blockchain. Lambda ensures that the
index is securely and efficiently updated on the blockchain, providing a reliable way to track and manage index changes.

```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── service                     <-- Source code for a lambda function
│   ├── main.go                 <-- Index-sync function code
│   ├── config.go               <-- Configuration code
│   ├── respond.go              <-- Response structure
│   ├── secret.go               <-- Get private key from SM
└── template.yaml
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM
  CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build
target.   
Read more
about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html)

The `sam build` command is wrapped inside of the `Makefile`. To execute this simply run

```shell
make
```

### Local development

**Invoking function locally through local API Gateway**

```bash
make run-local
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your
function `http://localhost:3000/hello`

**SAM CLI** is used to emulate both Lambda and API Gateway locally and uses our `template.yaml` to understand how to
bootstrap this environment (runtime, where the source code is, etc.) - The following excerpt is what the CLI will read
in order to initialize an API and its routes:

To deploy application for the first time, run the following in your shell:

```bash
make deploy ENV=<stage/prod>
```

The command will package and deploy application to AWS.

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang
website: https://golang.org/doc/install

A quickstart way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```

## Bringing to the next level

Here are a few ideas that you can use to get more acquainted as to how this overall process works:

* Create an additional API resource (e.g. /hello/{proxy+}) and return the name requested through this new path
* Update unit test to capture that
* Package & Deploy

Next, you can use the following resources to know more about beyond hello world samples and how others structure their
Serverless applications:

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)
