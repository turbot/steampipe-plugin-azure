
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
  location = "East US"
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = var.resource_name
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  location                 = azurerm_resource_group.named_test_resource.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_server" "named_test_resource" {
  name                         = var.resource_name
  resource_group_name          = azurerm_resource_group.named_test_resource.name
  location                     = azurerm_resource_group.named_test_resource.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  tags = {
    name = var.resource_name
  }
}

resource "azurerm_sql_firewall_rule" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  server_name         = azurerm_sql_server.named_test_resource.name
  start_ip_address    = "10.0.17.62"
  end_ip_address      = "10.0.17.62"
}

output "resource_aka" {
  value = "azure://${azurerm_sql_server.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_sql_server.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "firewall_rule_id" {
  value = azurerm_sql_firewall_rule.named_test_resource.id
}

output "resource_id" {
  value = azurerm_sql_server.named_test_resource.id
}

output "location" {
  value = lower(azurerm_sql_server.named_test_resource.location)
}

output "subscription_id" {
  value = var.azure_subscription
}

# To reduce the risk of accidentally exporting sensitive data that was intended
# to be only internal, Terraform requires that any root module output
# containing sensitive data be explicitly marked as sensitive, to confirm your
# intent.

# If you do intend to export this data, annotate the output value as sensitive
# by adding the following argument:
#     sensitive = true

# Though it is not being used any where so it is commented

# output "auditing_policy" {
#   value = azurerm_sql_server.named_test_resource.extended_auditing_policy
# }
