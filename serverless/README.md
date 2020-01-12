# Serverless - Highly Available Go Web App

**NOTE**: WIP - this project is **NOT** complete yet

## Description

Creates a Go web service running in a public cloud. This service exposes a
RESTful API that returns JSON formatted information. The architecture ensures
that the CI/CD pipeline can roll out updates frequently with zero failed
requests.

## Requirements

The service must be able to be:

  - Updated regularly with no dropped requests
  - Highly available
  - Easily scaled (preferably autoscaled)

## Design Overview

### Solution: Serverless

This solution uses both AWS CloudFormation and Serverless Framework.

  - Infrastructure (Cloud)
    - Provider: AWS
    - Provisioners:
      - CloudFormation
        - Common core infrastructure, e.g. VPC and networking
	- CodeBuild/CodePipeline (CI/CD)
	- ACM certificate
      - Serverless Framework
        - Deploys AWS REST API Gateway
	- Associated AWS Lambda function(s)
	- API custom domain name
	- Route 53 DNS record(s)
    - Networking:
      - VPC: isolate/protect infrastructure
      - Public subnet: NAT GWs
      - Private subnet: Lambda functions
      - NAT GW: public internet access for lambdas
      - Security groups: control connectivity access
    - Compute:
      - Lambda
        - VPC Private subnets (enables access to VPC resources)
        - Auto scaled - load, deployments
        - IAM role to access specific resources
        - Triggered from REST API Gateway
    - ACM to manage SSL/TLS certificates for the REST API
    - Route 53 to configure DNS with alias record to the API Gateway
  - VCS/CI/CD
    - VCS: git and GitHub
      - manage the application and infrastructure code.
      - With separate environment branches, e.g. dev, staging, prod
      - Test and promote from dev -> prod.
      - Feature branches used to test and submit PR’s.
      - PR’s merged on to specific branches and ultimately to prod/master for deployments
    - CI/CD: CodeBuildCodePipeline
      - configured to be triggered by repo hooks from pushes/merges to  branches
      - Performs tests and deploy.
      - Tests can be automated or manual prompting for authorization to deploy.
      - Deployment using Serverless Framework

## Set up

- clone repo
- install requirements
- deploy common core infrastructure via `Stacker` (`CloudFormation`)
- log into AWS console and connect CodePipeline and CodeBuild to Github
- CodePipeline will deploy serverless resources

## Updates

- CodePipeline will get triggered from pushes/merges to `dev` branch which will deploy updates using CodeBuild
- submit changes by merging feature branch off of dev branch or push push directly to `dev`
