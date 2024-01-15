variable "resource_name" {
  type        = string
  default     = "turbot-test-20200928-create-update"
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
      version = "=2.78.0"
    }
  }
}

provider "azurerm" {
  # Cannot be passed as a variable
  features {}
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_spring_cloud_service" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
}

resource "azurerm_spring_cloud_app" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  service_name        = azurerm_spring_cloud_service.named_test_resource.name

  identity {
    type = "SystemAssigned"
  }
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_spring_cloud_app.named_test_resource]
  value      = "azure://${azurerm_spring_cloud_app.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_spring_cloud_app.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_spring_cloud_app.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
