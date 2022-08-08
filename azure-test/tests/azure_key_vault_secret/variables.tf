
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

    secret_permissions = [
      "set",
      "get",
      "list",
      "delete",
      "purge",
      "recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "named_test_resource" {
  name         = var.resource_name
  value        = "steampipe"
  content_type = "text"
  key_vault_id = azurerm_key_vault.named_test_resource.id

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_key_vault_secret.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_key_vault_secret.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_key_vault_secret.named_test_resource.id
}

output "secret_version" {
  value = azurerm_key_vault_secret.named_test_resource.version
}

output "secret_uri_without_version" {
  value = azurerm_key_vault_secret.named_test_resource.versionless_id
}

output "location" {
  value = lower(azurerm_resource_group.named_test_resource.location)
}

output "subscription_id" {
  value = var.azure_subscription
}
