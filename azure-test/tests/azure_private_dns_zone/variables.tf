
variable "resource_name" {
  type        = string
  default     = "turbot-test-20230313-create-update"
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
  location = "East US"
  tags = {
    name = var.resource_name
  }
}

resource "azurerm_private_dns_zone" "named_test_resource" {
  name                = "${var.resource_name}.local"
  resource_group_name = azurerm_resource_group.named_test_resource.name

  tags = {
    name = "${var.resource_name}.local"
  }
}

output "resource_aka" {
  # Azure always puts privateDnsZones with lowercase d and z
  value = replace("azure://${azurerm_private_dns_zone.named_test_resource.id}", "Microsoft.Network/PrivateDnsZones", "Microsoft.Network/privatednszones")
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  # Azure always puts dnsZones with with lowercase d and z
  value = replace(azurerm_private_dns_zone.named_test_resource.id, "Microsoft.Network/PrivateDnsZones", "Microsoft.Network/privatednszones")
}

output "subscription_id" {
  value = var.azure_subscription
}
