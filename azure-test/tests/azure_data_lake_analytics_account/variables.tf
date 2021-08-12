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
  default     = "3510ae4d-530b-497d-8f30-53b9616fc6c1"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
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
  location = "East US2"
}

resource "azurerm_data_lake_store" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  encryption_state    = "Enabled"
  encryption_type     = "ServiceManaged"
}

resource "azurerm_data_lake_analytics_account" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location

  default_store_account_name = azurerm_data_lake_store.named_test_resource.name
}

output "resource_aka" {
  value = "azure://${azurerm_data_lake_analytics_account.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_data_lake_analytics_account.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_data_lake_analytics_account.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
