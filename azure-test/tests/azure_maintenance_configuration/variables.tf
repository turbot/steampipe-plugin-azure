
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c"
  description = "Azure subscription used for the test."
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.73.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West US"
}

resource "azurerm_maintenance_configuration" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  scope               = "SQLDB"

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_maintenance_configuration.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_maintenance_configuration.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = lower(azurerm_maintenance_configuration.named_test_resource.id)
}

output "location" {
  value = lower(azurerm_resource_group.named_test_resource.location)
}

output "subscription_id" {
  value = var.azure_subscription
}
