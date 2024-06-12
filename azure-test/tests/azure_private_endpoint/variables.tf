
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
  type = string
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
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_subnet" "service" {
  name                                          = var.resource_name
  resource_group_name                           = azurerm_resource_group.named_test_resource.name
  virtual_network_name                          = azurerm_virtual_network.named_test_resource.name
  address_prefixes                              = ["10.0.1.0/24"]
  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                                          = "endpoint"
  resource_group_name                           = azurerm_resource_group.named_test_resource.name
  virtual_network_name                          = azurerm_virtual_network.named_test_resource.name
  address_prefixes                              = ["10.0.2.0/24"]
  private_link_service_network_policies_enabled = false
}

resource "azurerm_public_ip" "named_test_resource" {
  name                = var.resource_name
  sku                 = "Standard"
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "named_test_resource" {
  name                = var.resource_name
  sku                 = "Standard"
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.named_test_resource.name
    public_ip_address_id = azurerm_public_ip.named_test_resource.id
  }
}

resource "azurerm_private_link_service" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  nat_ip_configuration {
    name      = azurerm_public_ip.named_test_resource.name
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.named_test_resource.frontend_ip_configuration[0].id,
  ]
}

resource "azurerm_private_endpoint" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = var.azure_environment
    private_connection_resource_id = azurerm_private_link_service.named_test_resource.id
    is_manual_connection           = false
  }
}

output "resource_aka" {
  value = "azure://${azurerm_private_endpoint.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_private_endpoint.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_private_endpoint.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
