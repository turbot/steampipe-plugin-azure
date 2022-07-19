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
  account_kind             = "StorageV2"
  access_tier              = "Cool"
  account_replication_type = "LRS"

  tags = {
    name = var.resource_name
  }
}
resource "azurerm_storage_container" "named_test_resource" {
  name                  = var.resource_name
  storage_account_name  = azurerm_storage_account.named_test_resource.name
  container_access_type = "private"
}

output "resource_aka" {
  value = "azure:///subscriptions/${var.azure_subscription}/resourceGroups/${var.resource_name}/providers/Microsoft.Storage/storageAccounts/${var.resource_name}/blobServices/default/containers/${var.resource_name}"
}

output "resource_aka_lower" {
  value = "azure:///subscriptions/${var.azure_subscription}/resourcegroups/${var.resource_name}/providers/microsoft.storage/storageaccounts/${var.resource_name}/blobservices/default/containers/${var.resource_name}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = "/subscriptions/${var.azure_subscription}/resourceGroups/${var.resource_name}/providers/Microsoft.Storage/storageAccounts/${var.resource_name}/blobServices/default/containers/${var.resource_name}"
}

output "subscription_id" {
  value = var.azure_subscription
}
