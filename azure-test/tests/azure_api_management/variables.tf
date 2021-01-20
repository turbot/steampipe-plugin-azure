
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

variable "azure_resource_group" {
  type        = string
  default     = "integration_test_rg"
  description = "Name of the resource group used throughout the test."
}

data "azurerm_resource_group" "data_resource_group" {
  name = var.azure_resource_group
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

resource "azurerm_api_management" "named_test_resource" {
  name                = var.resource_name
  location            = data.azurerm_resource_group.data_resource_group.location
  resource_group_name = var.azure_resource_group
  publisher_name      = "TurbotHQ"
  publisher_email     = "test@turbot.com"

  sku_name = "Developer_1"

  policy {
    xml_content = <<XML
  <policies>
    <inbound />
    <backend />
    <outbound />
    <on-error />
  </policies>
XML
  }
  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_api_management.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_api_management.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_api_management.named_test_resource.id
}

output "location" {
  value = data.azurerm_resource_group.data_resource_group.location
}

output "subscription_id" {
  value = var.azure_subscription
}

output "resource_group_name" {
  value = var.azure_resource_group
}
