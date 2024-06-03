variable "resource_name" {
  type        = string
  default     = "turbot-test-log-analytics-workspace"
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

provider "azurerm" {
  # Cannot be passed as a variable
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

data "azuread_client_config" "current" {}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West US"
}

resource "azurerm_log_analytics_workspace" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  depends_on = [azurerm_log_analytics_workspace.named_test_resource]
  value      = "azure://${azurerm_log_analytics_workspace.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_log_analytics_workspace.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = lower(azurerm_log_analytics_workspace.named_test_resource.id)
}
