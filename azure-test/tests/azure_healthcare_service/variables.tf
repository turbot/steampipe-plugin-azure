
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
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  version = "=2.41.0"
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
  location = "East US"
}

resource "azurerm_healthcare_service" "named_test_resource" {
  depends_on          = [azurerm_resource_group.named_test_resource]
  name                = var.resource_name
  resource_group_name = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  kind                = "fhir-R4"
  cosmosdb_throughput = "2000"

  tags = {
    "name" = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_healthcare_service.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_healthcare_service.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_healthcare_service.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
