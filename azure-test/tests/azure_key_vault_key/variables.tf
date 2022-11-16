
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
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c1"
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
  name                       = var.resource_name
  location                   = azurerm_resource_group.named_test_resource.location
  resource_group_name        = azurerm_resource_group.named_test_resource.name
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  tenant_id                  = data.azurerm_client_config.current.tenant_id

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "delete",
      "get",
      "purge",
      "recover"
    ]

    secret_permissions = [
      "set",
    ]
  }

  tags = {
    name = var.resource_name
  }
}

resource "azurerm_key_vault_key" "named_test_resource" {
  name         = var.resource_name
  key_vault_id = azurerm_key_vault.named_test_resource.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_key_vault_key.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_key_vault_key.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "location" {
  value = azurerm_resource_group.named_test_resource.location
}

output "location_lower" {
  value = lower(azurerm_resource_group.named_test_resource.location)
}

output "resource_id" {
  value = azurerm_key_vault_key.named_test_resource.id
}

output "key_version" {
  value = azurerm_key_vault_key.named_test_resource.version
}

output "key_uri_without_version" {
  value = azurerm_key_vault_key.named_test_resource.versionless_id
}

output "subscription_id" {
  value = var.azure_subscription
}
