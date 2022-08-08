variable "resource_name" {
  type        = string
  default     = "turbot-test-20201001-create-update"
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
  description = "Azure environment used for the test."
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.78.0"
    }
  }
}

provider "azurerm" {
  # Cannot be passed as a variable
  features {}
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West Europe"
}

resource "azurerm_virtual_network" "named_test_resource" {
  name                = var.resource_name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_subnet" "named_test_resource" {
  name                 = var.resource_name
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  virtual_network_name = azurerm_virtual_network.named_test_resource.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_hpc_cache" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.named_test_resource.id
  sku_name            = "Standard_2G"
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_hpc_cache.named_test_resource]
  value      = "azure://${azurerm_hpc_cache.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_hpc_cache.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_hpc_cache.named_test_resource.id
}

output "resource_id_upper" {
  value = "/subscriptions/${var.azure_subscription}/resourceGroups/${upper(var.resource_name)}/providers/Microsoft.StorageCache/caches/${var.resource_name}"
}

output "subscription_id" {
  value = var.azure_subscription
}

output "sku_name" {
  value = azurerm_hpc_cache.named_test_resource.sku_name
}
