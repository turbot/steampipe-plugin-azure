
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

resource "azurerm_data_factory" "named_test_resource" {
  name                = var.resource_name
  location            = "East US"
  resource_group_name = azurerm_resource_group.named_test_resource.name
  tags = {
    name = var.resource_name
  }
}

resource "azurerm_data_factory_linked_service_mysql" "named_test_resource" {
  name                = var.resource_name
  data_factory_id = azurerm_data_factory.named_test_resource.id
  connection_string   = "Server=test;Port=3306;Database=test;User=test;SSLMode=1;UseSystemTrustStore=0;Password=test"
}

resource "azurerm_data_factory_dataset_mysql" "named_test_resource" {
  name                = var.resource_name
  data_factory_id = azurerm_data_factory.named_test_resource.id
  linked_service_name = azurerm_data_factory_linked_service_mysql.named_test_resource.name
}

output "resource_aka" {
  value = "azure://${azurerm_data_factory_dataset_mysql.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_data_factory_dataset_mysql.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_data_factory_dataset_mysql.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
