
variable "resource_name" {
  type        = string
  default     = "turbot-test-20220425-create-update"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default    = "3510ae4d-530b-497d-8f30-53b9616fc6c1"
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
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_eventhub" "named_test_resource" {
  name                = var.resource_name
  namespace_name      = azurerm_eventhub_namespace.named_test_resource.name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  partition_count     = 2
  message_retention   = 2
}

resource "azurerm_eventhub_authorization_rule" "named_test_resource" {
  name                = var.resource_name
  namespace_name      = azurerm_eventhub_namespace.named_test_resource.name
  eventhub_name       = azurerm_eventhub.named_test_resource.name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  listen              = true
  send                = false
  manage              = false
}

resource "azurerm_security_center_automation" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  action {
    type              = "eventhub"
    resource_id       = azurerm_eventhub.named_test_resource.id
    connection_string = azurerm_eventhub_authorization_rule.named_test_resource.primary_connection_string
  }

  source {
    event_source = "Alerts"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "High"
        property_type  = "String"
      }
    }
  }

  scopes = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}"]
}

output "resource_aka" {
  value = replace("azure://${azurerm_security_center_automation.named_test_resource.id}", "resourceGroups", "resourcegroups")
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_security_center_automation.named_test_resource.id)}"
}

output "resource_id" {
  value = azurerm_security_center_automation.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "subscription_id" {
  value = var.azure_subscription
}
