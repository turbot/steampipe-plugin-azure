
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

resource "azurerm_key_vault" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  tags = {
    name = var.resource_name
  }
}

resource "azurerm_key_vault_access_policy" "named_test_resource" {
  key_vault_id = azurerm_key_vault.named_test_resource.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
  ]
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = var.resource_name
  location                 = azurerm_resource_group.named_test_resource.location
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_diagnostic_setting" "named_test_resource" {
  name               = var.resource_name
  target_resource_id = azurerm_key_vault.named_test_resource.id
  storage_account_id = azurerm_storage_account.named_test_resource.id

  log {
    category = "AuditEvent"
    enabled  = true

    retention_policy {
      enabled = true
      days    = 30
    }
  }
}

output "resource_aka" {
  value = "azure://${azurerm_key_vault.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_key_vault.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_key_vault.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "tenant_id" {
  value = data.azurerm_client_config.current.tenant_id
}

output "object_id" {
  value = data.azurerm_client_config.current.object_id
}

output "storage_account_id" {
  value = azurerm_storage_account.named_test_resource.id
}
