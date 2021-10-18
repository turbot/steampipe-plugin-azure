variable "resource_name" {
  type        = string
  default     = "turbot-test-20200908-service-clusters"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "3510ae4d-530b-497d-8f30-53b9616fc6c1"
  description = "Azure environment used for the test."
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.75.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_service_fabric_cluster" "named_test_resource" {
  name                 = var.resource_name
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  location             = azurerm_resource_group.named_test_resource.location
  reliability_level    = "Bronze"
  upgrade_mode         = "Manual"
  cluster_code_version = "7.2.413.9590"
  vm_image             = "Windows"
  management_endpoint  = "https://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_service_fabric_cluster.named_test_resource]
  value      = "azure://${azurerm_service_fabric_cluster.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_service_fabric_cluster.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_service_fabric_cluster.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "reliability_level" {
  value = azurerm_service_fabric_cluster.named_test_resource.reliability_level
}

output "upgrade_mode" {
  value = azurerm_service_fabric_cluster.named_test_resource.upgrade_mode
}

output "vm_image" {
  value = azurerm_service_fabric_cluster.named_test_resource.vm_image
}

output "management_endpoint" {
  value = azurerm_service_fabric_cluster.named_test_resource.management_endpoint
}
