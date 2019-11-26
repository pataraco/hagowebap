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

### Solution 1 (without containerization):

  - Infrastructure (Cloud)
    - Provider: AWS
    - Provisioner: Terraform or CloudForrmation
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
  - CI/CD
    - CI: git (GitHub or CodeCommit)
      - manage the application and infrastructure code.
      - With separate environment branches, e.g. dev, staging, prod
      - Test and promote from dev -> prod.
      - Feature branches used to test and submit PR’s.
      - PR’s merged on to specific branches and ultimately to prod/master for deployments
    - CD: Jenkins/Jenkinsfile (or CodeBuild/CodeDeploy)
      - configured to watch the code repo pushes/merges to prod/master branch
        - Or GitHub hooks to start Jenkins builds
      - Performs tests and deploy.
      - Tests can be automated or manual prompting for authorization to deploy.
      - Deployment by updating EC2 launch templates/configs and performing autoscaling rolling update

### Solution 2 (with containerization - specifically kubernetes):

Similar to above design above, with the following exceptions

  - Use EKS to deploy kubernetes master nodes
  - k8s: Use “cluster-autoscaler” and horizontal pod autoscaler (hpa) to create the AWS auto scaling group and dynamically scale pods and worker nodes
  - k8s: Use “alb-ingress-controller” k8s deployment to deploy the AWS ALB which will front the k8s worker nodes
  - Jenkins to pull/test code and build/push images to a container repo (ECR/Docker Hub). Then if/when ready/approved to deploy, update the images of the web app k8s deployments and specify a rolling deployment.

### Solution 3 (with containerization - AWS Fargate):

  - Mostly the same as Solution 2, just replace all kubernetes components with ECS/Fargate
  - Create the ALB as in Solution 1 to front the Fargate instances
  - Dynamic scaling of worker nodes handled by Fargate
  - Jenkins would update ECS task definitions for deployments
