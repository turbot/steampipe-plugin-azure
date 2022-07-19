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
  default     = "3510ae4d-530b-497d-8f30-53c0616fc6c1"
  description = "Azure environment used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

data "azuread_client_config" "current" {}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West Europe"
}

resource "azurerm_lb" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  tags = {
    name = var.resource_name
  }
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_lb.named_test_resource]
  value      = "azure://${azurerm_lb.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_lb.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_lb.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
