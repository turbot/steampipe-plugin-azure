variable "resource_name" {
  type        = string
  default     = "turbot-test-20200927-create-update"
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
  location = "West Europe"
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = var.resource_name
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  location                 = azurerm_resource_group.named_test_resource.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "named_test_resource" {
  name               = var.resource_name
  storage_account_id = azurerm_storage_account.named_test_resource.id
}

resource "azurerm_synapse_workspace" "named_test_resource" {
  name                                 = var.resource_name
  resource_group_name                  = azurerm_resource_group.named_test_resource.name
  location                             = azurerm_resource_group.named_test_resource.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.named_test_resource.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  tags = {
    Name = var.resource_name
  }
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_synapse_workspace.named_test_resource]
  value      = "azure://${azurerm_synapse_workspace.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_synapse_workspace.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_synapse_workspace.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "sql_administrator_login" {
 value = azurerm_synapse_workspace.named_test_resource.sql_administrator_login
}
