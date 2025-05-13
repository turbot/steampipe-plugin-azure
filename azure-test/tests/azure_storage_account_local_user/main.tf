provider "azurerm" {
  features {}
}

provider "random" {}

data "azurerm_client_config" "current" {}

resource "random_integer" "suffix" {
  min = 10000
  max = 99999
}

locals {
  resource_name = "${var.resource_name}-${random_integer.suffix.result}"
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = local.resource_name
  location = var.azure_location
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = replace(local.resource_name, "-", "")
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  location                 = azurerm_resource_group.named_test_resource.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled          = true
  sftp_enabled            = true
}

resource "azurerm_storage_account_local_user" "named_test_resource" {
  name                 = local.resource_name
  storage_account_id   = azurerm_storage_account.named_test_resource.id
  ssh_key_enabled     = true
  home_directory      = "/test"
  permission_scope {
    permissions {
      read   = true
      write  = true
      delete = true
      list   = true
    }
    service       = "blob"
    resource_name = "test"
  }
  ssh_authorized_key {
    description = "test key"
    key        = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsEBLXCsDtXtuskHcS6oChLsHj0JYp7ml1gVVBKAsGGAlT80BGqJyXfRwruFqaJvE/KHpe9LWXFaT05kHWGR2Dc3bZHCooOVdiblkYm4rvfuPG8iktE6YvVV3nTUfgb9oj/HIRmj/8xoVbD6zjRWAw7dr9EggUDiy86HSJ+6V8tMziRnSWJ4W2+QSDdNTCxSTKaaP3RV/nRky8cWYGfx4yiU+EH8/odPeeCc+Ts6iVcc5h5gfnqC1M8RZI3Zbnxz+v5vxY1LlR381wSMMJP8dmwnm6oQuor84fNN5AYDSLM4AxJRUqpgO6fik1saE6c9a5/V554aDt2Z7TkrwoYeCb"
  }
}

output "resource_aka" {
  value = "azure:///subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.named_test_resource.name}/providers/Microsoft.Storage/storageAccounts/${azurerm_storage_account.named_test_resource.name}/localUsers/${azurerm_storage_account_local_user.named_test_resource.name}"
}

output "resource_aka_list" {
  value = [
    "azure:///subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.named_test_resource.name}/providers/Microsoft.Storage/storageAccounts/${azurerm_storage_account.named_test_resource.name}/localUsers/${azurerm_storage_account_local_user.named_test_resource.name}"
  ]
}

output "resource_id" {
  value = azurerm_storage_account_local_user.named_test_resource.id
}

output "subscription_id" {
  value = data.azurerm_client_config.current.subscription_id
} 