
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
  # default     = "3510ae4d-530b-497d-8f30-53b9616fc6c1"
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  version         = "=2.43.0"
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

resource "azurerm_express_route_circuit" "named_test_resource" {
  name                  = var.resource_name
  resource_group_name   = azurerm_resource_group.named_test_resource.name
  location              = azurerm_resource_group.named_test_resource.location
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  tags = {
    Name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_express_route_circuit.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_express_route_circuit.named_test_resource.id)}"
}

output "sku_tire" {
  value = "Standard"
}

output "sku_family" {
  value = "MeteredData"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_express_route_circuit.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
