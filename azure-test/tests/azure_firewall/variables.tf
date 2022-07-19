
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
  default     = "3510ae4d-530b-497d-8f30-53b9616fc6c1"
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
  location = "West US"
}

resource "azurerm_virtual_network" "named_test_resource" {
  name                = var.resource_name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_subnet" "named_test_resource" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  virtual_network_name = azurerm_virtual_network.named_test_resource.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku_tier = "Standard"

  ip_configuration {
    name                 = var.resource_name
    subnet_id            = azurerm_subnet.named_test_resource.id
    public_ip_address_id = azurerm_public_ip.named_test_resource.id
  }

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_firewall.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_firewall.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_firewall.named_test_resource.id
}

output "public_ip_id" {
  value = azurerm_public_ip.named_test_resource.id
}

output "subnet_id" {
  value = azurerm_subnet.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
