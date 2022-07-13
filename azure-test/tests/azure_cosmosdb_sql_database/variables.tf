
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
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  # version         = "=2.41.0"
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

resource "azurerm_cosmosdb_account" "named_test_resource" {
  name = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  offer_type = "Standard"
  kind = "GlobalDocumentDB"
  enable_automatic_failover = true
  
  consistency_policy {
    consistency_level = "BoundedStaleness"
    max_interval_in_seconds = 600
    max_staleness_prefix = 200000
  }

  geo_location {
    prefix = "turbot-cosmos-db"
    location = azurerm_resource_group.named_test_resource.location
    failover_priority = 0
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 360"
  }
  triggers = {
    "before" = "${azurerm_cosmosdb_account.named_test_resource.id}"
  }
}

resource "azurerm_cosmosdb_sql_database" "named_test_resource" {
  depends_on = [null_resource.delay]
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  account_name        = azurerm_cosmosdb_account.named_test_resource.name
  throughput          = 400
}

output "resource_aka" {
  value = "azure://${azurerm_cosmosdb_sql_database.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_cosmosdb_sql_database.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_cosmosdb_sql_database.named_test_resource.id
}
