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
  depends_on                 = [azurerm_resource_group.named_test_resource]
  name                       = var.resource_name
  location                   = azurerm_resource_group.named_test_resource.location
  resource_group_name        = azurerm_resource_group.named_test_resource.name
  sku_name                   = "standard"
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
}

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [azurerm_key_vault.named_test_resource]
  provisioner "local-exec" {
    command = "az keyvault delete --name ${var.resource_name} --resource-group ${var.resource_name}"
  }
  provisioner "local-exec" {
    command = "az keyvault list-deleted > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "resource_aka" {
  depends_on = [null_resource.named_test_resource]
  value      = "azure://${jsondecode(data.local_file.input.content)[0].id}"
}

output "resource_aka_lower" {
  depends_on = [null_resource.named_test_resource]
  value      = "azure://${lower(jsondecode(data.local_file.input.content)[0].id)}"
}

output "region" {
  value = jsondecode(data.local_file.input.content)[0].properties.location
}

output "resource_name" {
  value = jsondecode(data.local_file.input.content)[0].name
}

output "resource_id" {
  value = jsondecode(data.local_file.input.content)[0].id
}

output "subscription_id" {
  value = var.azure_subscription
}

