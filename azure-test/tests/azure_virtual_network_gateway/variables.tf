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

resource "azurerm_virtual_network" "named_test_resource" {
  name                = var.resource_name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_subnet" "named_test_resource" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  virtual_network_name = azurerm_virtual_network.named_test_resource.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "named_test_resource" {
  depends_on          = [azurerm_subnet.named_test_resource]
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  allocation_method   = "Dynamic"
}

locals {
  path = "${path.cwd}/info.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [azurerm_public_ip.named_test_resource]
  provisioner "local-exec" {
    command = "az network vnet-gateway create --gateway-type Vpn --location ${azurerm_resource_group.named_test_resource.location} --name ${var.resource_name} --no-wait --public-ip-addresses ${var.resource_name} --resource-group ${var.resource_name} --vnet ${var.resource_name}"
  }
  provisioner "local-exec" {
    command = "az network vnet-gateway show -g ${var.resource_name} -n ${var.resource_name} > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "resource_aka" {
  depends_on = [null_resource.named_test_resource]
  value      = "azure://${jsondecode(data.local_file.input.content).id}"
}

output "resource_aka_lower" {
  depends_on = [null_resource.named_test_resource]
  value      = "azure://${lower(jsondecode(data.local_file.input.content).id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = jsondecode(data.local_file.input.content).id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "region" {
  value = azurerm_resource_group.named_test_resource.location
}
