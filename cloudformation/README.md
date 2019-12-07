# CloudFormation Infrastructure as Code (IaC)

**THIS README IS NOT COMPLETE!**

## Overview

The backend infrastructure is deployed and configured using AWS CloudFormation
and `stacker`

## Setup (`stacker`)

`stacker` uses 3 main files to deploy CloudFormation templates

1. Environments file

   Contains specific environment settings that are used to assign values to the
   CloudFormation stack parameters. These variables are referenced in the stack
   definitions YAML file.

2. Stack definitions YAML file

   Contains the list of CloudFormation stacks to deploy and the values to
   provide to the CloudFormation templates' parameters. These are numbered to
   signify order of deployment preference.

3. CloudFormation Templates

   The actual CloudFormation templates.

## Deploying

For each `environment` you'll need a corresponding environment file.
The naming convention of the environment file is
`$ENVIRONMENT.$AWS_REGION.env`.

### 1. Environment Settings

Make any changes to the necessary variables in the specific environment file
corresponding to the environment that you want to deploy/update, e.g.

* environment
* namespace
* region
* instance sizes
* stacker_bucket_name

### 2. CloudFormation Stack Parameters

Make any necessary changes to the parameters you want to specify or use the
default values of in the stack definition YAML files. You'll need to look at
the actual CloudFormation templates to see which parameters can be specified
and what their default values are. In some case, the default values have been
specified in the `environment` files.

### 3. Deploy

Use `stacker` to build/deploy the stacks individually, the syntax is:

`stacker build -i -r $REGION $ENVIRONMENT_FILE $YAML_FILE`
("-i" for interactive mode)

So, for example to build the `landing-zone` stack for `dev` in `us-west-2`,
run:

`stacker build -i -r us-west-2 dev-us-west-2.env 01-landing-zone.yaml`

### Destroying

You can destroy the stacks in the AWS console or with the `stacker destroy`
command, the syntax is:

`stacker destroy -f|--force -r $REGION $ENVIRONMENT_FILE $YAML_FILE`

So, for example to destroy the `landing-zone` stack for `dev` in `us-west-2`,
run:

`stacker destroy -f -r us-west-2 dev-us-west-2.env 01-landing-zone.yaml`
    
You should destroy stacks in reverse order as they were created due to
interdependencies.  You can look at the progress/status in the AWS console
CloudFormation section.  If failures occur look in the events to determine the
issue, most likely resources can not be deleted because they are referred to by
others.  You may have to delete these manually.  Then run the `stacker destroy`
command again or delete the stack using the AWS console or CLI.

## Requirements

1. [stacker](https://github.com/remind101/stacker) v1.7.0+
   [stacker (pip)](https://pypi.python.org/pypi/stacker) -
   See [Stacker Documentation](http://stacker.readthedocs.io/en/latest),
   [Introducing Stacker](http://engineering.remind.com/introduction-to-stacker)
2. [GitPython](https://pypi.python.org/pypi/GitPython) (this should be
   automatically installed by Stacker)

## Cloudformation Templates S3 Bucket Save Location

Defined in environment settings file(s) by `stacker_bucket_name`

* e.g. PROJECT-ENV-cloudformation

## AWS Services
* [AWS CloudFormation](https://aws.amazon.com/cloudformation/)
* [AWS Identity and Access Management (IAM)](https://aws.amazon.com/iam)
* [Amazon Route 53](https://aws.amazon.com/route53)
* [Amazon Simple Storage Service (Amazon S3)](https://aws.amazon.com/s3)
* [Amazon Virtual Private Cloud (Amazon VPC)](https://aws.amazon.com/vpc)

## CloudFormation Resources
* [AWS::CloudFormation::Stack](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stack.html)
* [AWS::EC2::EIP](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-eip.html)
* [AWS::EC2::InternetGateway](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-internetgateway.html)
* [AWS::EC2::NatGateway](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-natgateway.html)
* [AWS::EC2::Route](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route.html)
* [AWS::EC2::RouteTable](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route-table.html)
* [AWS::EC2::SecurityGroup](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-security-group.html)
* [AWS::EC2::Subnet](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-subnet.html)
* [AWS::EC2::SubnetRouteTableAssociation](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-subnet-route-table-assoc.html)
* [AWS::EC2::VPC](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpc.html)
* [AWS::EC2::VPCGatewayAttachment](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpc-gateway-attachment.html)
* [AWS::IAM::ManagedPolicy](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-managedpolicy.html)
* [AWS::IAM::Role](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-role.html)
* [AWS::Route53::RecordSet](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-recordset.html)
* [AWS::S3::Bucket](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket.html)

## Updates

This infrastructure automation is meant for initial environment configuraiton
and overall infrastructure updates.  Infrastructure updates should/could be
made by updating the CloudFormation templates, `stacker` YAML files and
environment settings. API Gateway and Lambda function updates should be made
via a CI/CD implementation, e.g. AWS CodePipeline, CodeBuild and CodeDeploy
and/or Jenkins.

## Important Notes

* **CAUTION**: Some Cognito "updates" require a new stack completely - so it
would need to be treated like a DB with either migrations or spinning up and
backfilling. When you run `stacker` with the interactive option (`-i`) you can
view the changes to be made and verify if resources will be replaced or not.

## Troubleshooting

Typically `stacker` will fail for the following common reasons:

* Missing/Undefined variables in the environment file (`*.env`) referenced in
  the configuration (`YAML`) file
* Missing/Undefined/Mismatched variables in the configuration (`YAML`) file
  referenced in the CloudFormation template file
* CloudFormation errors

Pay attention to and look closely to the `stacker` Python error/traceback. It
hould be fairly obvious where the issue is. And, if it is a CloudFormation
error, it should list which stack failed. Go into the AWS console and into
CloudFormation and look at the events of the failed stack to get more insight
into the issue.