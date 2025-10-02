package terraform

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// VariablesGenerator handles Terraform variables generation
type VariablesGenerator struct{}

// NewVariablesGenerator creates a new variables generator
func NewVariablesGenerator() *VariablesGenerator {
	return &VariablesGenerator{}
}

// GenerateVariables generates variables.tf content
func (vg *VariablesGenerator) GenerateVariables(config *models.ProjectConfig) string {
	return fmt.Sprintf(`variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "%s"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b"]
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "min_capacity" {
  description = "Minimum number of instances"
  type        = number
  default     = 1
}

variable "max_capacity" {
  description = "Maximum number of instances"
  type        = number
  default     = 3
}

variable "desired_capacity" {
  description = "Desired number of instances"
  type        = number
  default     = 2
}`, config.Name)
}

// GenerateOutputs generates outputs.tf content
func (vg *VariablesGenerator) GenerateOutputs(config *models.ProjectConfig) string {
	return `output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

output "load_balancer_dns_name" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "load_balancer_zone_id" {
  description = "Zone ID of the load balancer"
  value       = aws_lb.main.zone_id
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}`
}

// GenerateTfVarsExample generates terraform.tfvars.example content
func (vg *VariablesGenerator) GenerateTfVarsExample(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s Terraform Variables Example
# Copy this file to terraform.tfvars and update the values

aws_region = "us-west-2"
project_name = "%s"
environment = "dev"

# VPC Configuration
vpc_cidr = "10.0.0.0/16"
availability_zones = ["us-west-2a", "us-west-2b"]

# Instance Configuration
instance_type = "t3.micro"
min_capacity = 1
max_capacity = 3
desired_capacity = 2`, config.Name, config.Name)
}
