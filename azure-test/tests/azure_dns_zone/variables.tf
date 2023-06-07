
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

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
  tags = {
    name = var.resource_name
  }
}

resource "azurerm_dns_zone" "named_test_resource" {
  name                = "${var.resource_name}.xyz"
  resource_group_name = azurerm_resource_group.named_test_resource.name

  tags = {
    name = "${var.resource_name}.xyz"
  }
}

output "resource_aka" {
  # Azure always puts dnsZones with a lowercase z
  value = replace("azure://${azurerm_dns_zone.named_test_resource.id}", "Microsoft.Network/dnsZones", "Microsoft.Network/dnszones")
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  # Azure always puts dnsZones with a lowercase z
  value = replace(azurerm_dns_zone.named_test_resource.id, "Microsoft.Network/dnsZones", "Microsoft.Network/dnszones")
}

output "subscription_id" {
  value = var.azure_subscription
}
