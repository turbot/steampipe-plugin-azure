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

# Make sure Network Watcher is disabled, in a particular region that are mentioned
# Creation of more than 1 network watchers for this subscription in this region, is not possible
resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_network_security_group" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_network_watcher" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_storage_account" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location

  account_tier              = "Standard"
  account_kind              = "StorageV2"
  account_replication_type  = "LRS"
  enable_https_traffic_only = true
}

resource "azurerm_log_analytics_workspace" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "named_test_resource" {
  name                 = var.resource_name
  network_watcher_name = azurerm_network_watcher.named_test_resource.name
  resource_group_name  = azurerm_resource_group.named_test_resource.name

  network_security_group_id = azurerm_network_security_group.named_test_resource.id
  storage_account_id        = azurerm_storage_account.named_test_resource.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = true
    workspace_id          = azurerm_log_analytics_workspace.named_test_resource.workspace_id
    workspace_region      = azurerm_log_analytics_workspace.named_test_resource.location
    workspace_resource_id = azurerm_log_analytics_workspace.named_test_resource.id
    interval_in_minutes   = 10
  }
}

output "resource_aka" {
  value = "azure:///subscriptions/${var.azure_subscription}/resourceGroups/${var.resource_name}/providers/Microsoft.Network/networkWatchers/${var.resource_name}/flowLogs/Microsoft.Network${var.resource_name}${var.resource_name}"
}

output "resource_aka_lower" {
  value = "azure:///subscriptions/${var.azure_subscription}/resourcegroups/${var.resource_name}/providers/microsoft.network/networkwatchers/${var.resource_name}/flowlogs/microsoft.network${var.resource_name}${var.resource_name}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_network_watcher_flow_log.named_test_resource.id
}

output "storage_account_id" {
  value = azurerm_storage_account.named_test_resource.id
}

output "network_security_group_id" {
  value = azurerm_network_security_group.named_test_resource.id
}

output "workspace_id" {
  value = azurerm_log_analytics_workspace.named_test_resource.workspace_id
}

output "subscription_id" {
  value = var.azure_subscription
}
