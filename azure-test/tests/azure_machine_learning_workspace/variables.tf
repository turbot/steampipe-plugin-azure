variable "resource_name" {
  type        = string
  default     = "steampipe-test"
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
  # Cannot be passed as a variable
  version         = "=2.41.0"
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

resource "azurerm_application_insights" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  application_type    = "web"
}

resource "azurerm_key_vault" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = var.resource_name
  location                 = azurerm_resource_group.named_test_resource.location
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_machine_learning_workspace" "named_test_resource" {
  name                    = var.resource_name
  location                = azurerm_resource_group.named_test_resource.location
  resource_group_name     = azurerm_resource_group.named_test_resource.name
  application_insights_id = azurerm_application_insights.named_test_resource.id
  key_vault_id            = azurerm_key_vault.named_test_resource.id
  storage_account_id      = azurerm_storage_account.named_test_resource.id

  identity {
    type = "SystemAssigned"
  }

  tags = {
    "name" = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_machine_learning_workspace.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_machine_learning_workspace.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_machine_learning_workspace.named_test_resource.id
}

output "location" {
  value = azurerm_resource_group.named_test_resource.location
}

output "subscription_id" {
  value = var.azure_subscription
}
