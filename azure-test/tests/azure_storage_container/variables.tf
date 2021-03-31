
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
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  version         = "=1.36.0"
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
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
  value = "azure://${azurerm_storage_container.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_storage_container.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_storage_container.named_test_resource.id
}
