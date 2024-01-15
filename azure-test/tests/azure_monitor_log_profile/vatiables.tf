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
  default     = "3510ae4d-530b-497d-8f30-53b9616fc6c1"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
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
  location = "East US"
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = var.resource_name
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  location                 = azurerm_resource_group.named_test_resource.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_eventhub_namespace" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_monitor_log_profile" "named_test_resource" {
  name = var.resource_name

  categories = [
    "Action",
    "Delete",
    "Write",
  ]

  locations = [
    "eastus",
    "global",
  ]

  # RootManageSharedAccessKey is created by default with listen, send, manage permissions
  servicebus_rule_id = "${azurerm_eventhub_namespace.named_test_resource.id}/authorizationrules/RootManageSharedAccessKey"
  storage_account_id = azurerm_storage_account.named_test_resource.id

  retention_policy {
    enabled = true
    days    = 7
  }
}

output "resource_aka" {
  value = "azure://${azurerm_monitor_log_profile.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_monitor_log_profile.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id_lower" {
  value = lower(azurerm_monitor_log_profile.named_test_resource.id)
}

output "storage_account_id" {
  value = azurerm_storage_account.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "location" {
  value = azurerm_resource_group.named_test_resource.location
}

