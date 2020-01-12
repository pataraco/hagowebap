# Highly Available Go Web App

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

### Solution 1: Dedicated Servers

This solution a pure AWS solution and uses CloudFormation, CodeCommit, CodeBuild, and CodePipeline to deploy and update the application in AWS. A VPC is used to isolate the resources. Please find [more information about the solution here](dedicated/README.md)

### Solution 2: Containerization

This solution uses Terraform, Kubernetes and Jenkins to deploy and update the application in AWS. A VPC is used to isolate the resources. Please find [more information about the solution here](containerized/README.md)

### Solution 3: Serverless

This solution uses both AWS CloudFormation and Serverless Framework to deploy AWS Lambda functions and an AWS REST API Gateway. A VPC is used to isolate the resources. Please find [more information about the solution here](serverless/README.md)
