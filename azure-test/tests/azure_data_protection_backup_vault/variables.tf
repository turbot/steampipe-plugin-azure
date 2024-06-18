variable "resource_name" {
  type        = string
  default     = "steampipe-test-03042024"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "00000000-0000-0000-0000-000000000000"
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

resource "azurerm_data_protection_backup_vault" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

output "resource_aka" {
  value = "azure://${azurerm_data_protection_backup_vault.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_data_protection_backup_vault.named_test_resource.id)}"
}

output "id" {
  value = azurerm_data_protection_backup_vault.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}
