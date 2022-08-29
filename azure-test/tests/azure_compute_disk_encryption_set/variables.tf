
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
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c"
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

resource "azurerm_key_vault" "named_test_resource" {
  name                        = var.resource_name
  location                    = azurerm_resource_group.named_test_resource.location
  resource_group_name         = azurerm_resource_group.named_test_resource.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
  soft_delete_retention_days  = 7
}

resource "azurerm_key_vault_key" "named_test_resource" {
  name         = var.resource_name
  key_vault_id = azurerm_key_vault.named_test_resource.id
  key_type     = "RSA"
  key_size     = 2048

  depends_on = [
    azurerm_key_vault_access_policy.named_test_resource
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

resource "azurerm_disk_encryption_set" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  key_vault_key_id    = azurerm_key_vault_key.named_test_resource.id

  identity {
    type = "SystemAssigned"
  }

  tags = {
    "name" = var.resource_name
  }
}

resource "azurerm_key_vault_access_policy" "named_test_resource" {
  key_vault_id = azurerm_key_vault.named_test_resource.id

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign"
  ]
}

resource "azurerm_role_assignment" "example-disk" {
  scope                = azurerm_key_vault.named_test_resource.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = azurerm_disk_encryption_set.named_test_resource.identity.0.principal_id
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
  value = upper(var.resource_name)
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
