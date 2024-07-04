# Index sync chain

### Overview

This service designed to update an index from one resource and send it to a blockchain. This service
ensures that the index is securely and efficiently updated on the blockchain, providing a reliable way to track and
manage index changes.

The project includes a deployment process that builds and pushes the Docker image 
to AWS Elastic Container Registry (ECR) and updates an AWS Lambda function. The Lambda function is 
scheduled to run every hour using AWS CloudWatch Events.

## Prerequisites

- Go 1.22.4
- Docker
- AWS CLI
- AWS account with appropriate permissions
- GitHub repository for CI/CD


### Config File

Create a config file or use environment variables to provide the necessary configuration:

```json
{
  "LiteConnectionsURLTestnet": "<your_testnet_url>",
  "LiteConnectionsURLMainnet": "<your_mainnet_url>",
  "MasterContractAddress": "<your_master_contract_address>",
  "TONStakingContractAddress": "<your_staking_contract_address>"
}
```

### GitHub Actions Workflow

This project uses GitHub Actions for CI/CD. 
The GitHub Actions workflow is defined in `.github/workflows/on_push_main.yml`
- check out the code
- Set up Go
- Build the Go application
- Configure AWS credentials
- Log in to Amazon ECR
- Build, tag, and push the Docker image to ECR
- Update the AWS Lambda function with the new image URI


### AWS Lambda Deployment

The Lambda function is updated with the new Docker image by the GitHub Actions workflow.
The Lambda function is configured to run every hour using AWS CloudWatch Events.