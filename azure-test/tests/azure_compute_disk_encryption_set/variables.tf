
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
  # Cannot be passed as a variable
  version         = "=2.41.0"
  features {}
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

resource "azurerm_key_vault" "named_test_resource" {
  name                        = var.resource_name
  location                    = azurerm_resource_group.named_test_resource.location
  resource_group_name         = azurerm_resource_group.named_test_resource.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  enabled_for_disk_encryption = true
  soft_delete_enabled         = true
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true
  sku_name                    = "standard"
}

resource "azurerm_key_vault_key" "named_test_resource" {
  name         = var.resource_name
  key_vault_id = azurerm_key_vault.named_test_resource.id
  key_type     = "RSA"
  key_size     = 2048

  depends_on = [
    azurerm_key_vault_access_policy.local_user,
  ]
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "local_user" {
  key_vault_id = azurerm_key_vault.named_test_resource.id
  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
  key_permissions = [
    "create",
    "get",
    "list",
    "delete",
  ]
}

resource "azurerm_disk_encryption_set" "named_test_resource" {
  name = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  key_vault_key_id = azurerm_key_vault_key.named_test_resource.id

  identity {
    type = "SystemAssigned"
  }

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_disk_encryption_set.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_disk_encryption_set.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_name_upper" {
  value = "${upper(var.resource_name)}"
}

output "resource_id" {
  value = azurerm_disk_encryption_set.named_test_resource.id
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

output "vault_id" {
  value = azurerm_key_vault.named_test_resource.id
}

output "key_id" {
  value = azurerm_key_vault_key.named_test_resource.id
}
