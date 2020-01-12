# Dedicated Web Servers - Highly Available Go Web App

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

### Solution: Dedicated Web Servers (AWS EC2)

This is a pure AWS solution.

  - Infrastructure (Cloud)
    - Provider: AWS
    - Provisioner:
      - CloudFormation
        - Common core infrastructure, e.g. VPC and networking
	- EC2 instances, autoscaling and ALB/Listeners/Target Groups
	- CodeCommit/CodeBuild/CodePipeline (VCS/CI/CD)
	- ACM certificate
	- Route 53 DNS record(s)
    - Networking:
      - VPC: isolate/protect infrastructure
      - Public subnet: ALB and Jenkins server
      - Private subnet: Web Server(s)
      - NAT GW: public internet access for web servers
      - Security groups: control connectivity access
    - Compute:
      - EC2 - web server(s)
        - Private subnet
        - Auto scaled - outages, load, deployments
        - Security group - access from ALB
      - Auto Scaling groups:
        - For Web Servers
        - Used to roll out app updates with update EC2 launch templates/configs
        - Scaling policy configured to dynamically scale out/in based on load (or other desired metric)
        - Registered with an ALB and ELB health check - to automatically replace unhealthy web servers
      - Launch template/config
        - Specify server requirements/specs
        - Place on private subnets
        - Attach security groups
        - Configure user data script to install necessary packages and pull/start correct/tagged application version from CI repo (GitHub/CodeCommit)
      - EC2 - Jenkins server (CD)  (or can use CodeBuild/CodeDeploy)
        - Public subnet for GitHub integration
        - IAM role to give access to specify resources needed for deployments (as described below). Eliminates the need to use AWS keys.
        - Security group to only allow traffic from GitHub and Office
        - Can also use as a build server, or use plugins to dynamically launch build servers (spot instances or ECS/Fargate and containers)
      - Load Balancing:
        - ALB on a public subnet, with http/https listeners and target groups, allowing/providing:
        - HA during application updates/rollouts
        - SSL/TLS termination (ACM certs)
        - HTTP redirection to HTTPS
        - Authentication
        - Host/Path based routing
        - Public access to application
        - Serves/Fronts the web servers (created by the auto scaling group)
        - Security group allows HTTP/HTTPS traffic from all
    - ACM to manage SSL/TLS certificates for the web app
    - Route 53 to configure DNS with alias record to the ALB
  - VCS/CI/CD
    - VCS: git and AWS CodeCommit
      - manage the application and infrastructure code.
      - With separate environment branches, e.g. dev, staging, prod
      - Test and promote from dev -> prod.
      - Feature branches used to test and submit PR’s.
      - PR’s merged on to specific branches and ultimately to prod for deployments
    - CI/CD: CodeBuild/CodeDeploy/CodePipeline
      - configured to watch the code repo pushes/merges to prod/master branch via hooks to start CodePipeline
      - Performs tests and deploy.
      - Tests can be automated or manual prompting for authorization to deploy.
      - Deployment by updating EC2 launch templates/configs and performing autoscaling rolling update
