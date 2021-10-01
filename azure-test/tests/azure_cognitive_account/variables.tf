variable "resource_name" {
  type        = string
  default     = "turbot-test-20200907-cognitive-account"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
  description = "Azure environment used for the test."
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.65.0"
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

resource "azurerm_cognitive_account" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  kind                = "Face"

  sku_name = "S0"

  tags = {
    name = var.resource_name
  }
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_cognitive_account.named_test_resource]
  value      = "azure://${azurerm_cognitive_account.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_cognitive_account.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_cognitive_account.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "kind" {
  value = azurerm_cognitive_account.named_test_resource.kind
}

output "sku_name" {
  value = azurerm_cognitive_account.named_test_resource.sku_name
}
