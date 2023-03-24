
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

resource "azurerm_cosmosdb_account" "named_test_resource" {
  name                      = var.resource_name
  location                  = azurerm_resource_group.named_test_resource.location
  resource_group_name       = azurerm_resource_group.named_test_resource.name
  offer_type                = "Standard"
  kind                      = "MongoDB"
  enable_free_tier          = true
  enable_automatic_failover = false

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  capabilities {
    name = "MongoDBv3.4"
  }

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 1000
  }


  geo_location {
    location          = azurerm_resource_group.named_test_resource.location
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

resource "azurerm_cosmosdb_mongo_database" "named_test_resource" {
  depends_on          = [null_resource.delay]
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  account_name        = azurerm_cosmosdb_account.named_test_resource.name
}

resource "azurerm_cosmosdb_mongo_collection" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  account_name        = azurerm_cosmosdb_account.named_test_resource.name
  database_name       = azurerm_cosmosdb_mongo_database.named_test_resource.name

  default_ttl_seconds = "777"
  shard_key           = "uniqueKey"
  throughput          = 400

  index {
    keys   = ["_id"]
    unique = true
  }
}


output "resource_aka" {
  value = "azure://${azurerm_cosmosdb_mongo_collection.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_cosmosdb_mongo_collection.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_cosmosdb_mongo_collection.named_test_resource.id
}
