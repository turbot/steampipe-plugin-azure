variable "resource_name" {
  type        = string
  default     = "turbot-test-20210921-create-update"
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
  description = "Azure environment used for the test."
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.77.0"
    }
  }
}

provider "azurerm" {
  # Cannot be passed as a variable
  features {}
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

data "azuread_client_config" "current" {}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West Europe"
}

resource "azurerm_storage_sync" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  tags = {
    name = var.resource_name
  }
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_storage_sync.named_test_resource]
  value      = "azure://${azurerm_storage_sync.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_storage_sync.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_storage_sync.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
