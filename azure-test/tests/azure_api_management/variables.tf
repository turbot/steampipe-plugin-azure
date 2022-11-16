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
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c1"
  description = "Azure subscription used for the test."
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
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

resource "azurerm_resource_group" "named_test_resource" {
 name     = var.resource_name
 location = "West Europe"
}

resource "azurerm_api_management" "named_test_resource" {
 name                = var.resource_name
 location            = azurerm_resource_group.named_test_resource.location
 resource_group_name = azurerm_resource_group.named_test_resource.name
 publisher_name      = "TurbotHQ"
 publisher_email     = "test@turbot.com"

 sku_name = "Developer_1"
}

output "region" {
  depends_on = [azurerm_api_management.named_test_resource]
  value      = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_api_management.named_test_resource]
  value      = "azure://${azurerm_api_management.named_test_resource.id}"
}

output "resource_aka_lower" {
  depends_on = [azurerm_api_management.named_test_resource]
  value      = "azure://${lower(azurerm_api_management.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  depends_on = [azurerm_api_management.named_test_resource]
  value      = azurerm_api_management.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
