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
  default     = "cdffd708-7da0-4cea-abeb-0a4c334d7f64"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

data "azurerm_client_config" "current" {}

data "null_data_source" "resource" {
  inputs = {
    scope = "azure:///subscriptions/${data.azurerm_client_config.current.subscription_id}"
  }
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West US"
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = var.resource_name
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  location                 = azurerm_resource_group.named_test_resource.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_diagnostic_setting" "named_test_resource" {
  name               = var.resource_name
  target_resource_id = "/subscriptions/${var.azure_subscription}"
  storage_account_id = azurerm_storage_account.named_test_resource.id

  log {
    category = "Alert"
    enabled  = true

    retention_policy {
      enabled = false
    }
  }
}

output "resource_aka" {
  value = "azure://subscriptions/${var.azure_subscription}/providers/microsoft.insights/diagnosticSettings/${var.resource_name}"
}

output "resource_aka_lower" {
  value = "azure://subscriptions/${var.azure_subscription}/providers/microsoft.insights/diagnosticsettings/${var.resource_name}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = "subscriptions/${var.azure_subscription}/providers/microsoft.insights/diagnosticSettings/${var.resource_name}"
}

output "subscription_id" {
  value = var.azure_subscription
}
