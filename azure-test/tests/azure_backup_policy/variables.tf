variable "resource_name" {
  type        = string
  default     = "steampipe-test-03042024"
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

resource "azurerm_recovery_services_vault" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku                 = "Standard"
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 700"
  }
}

resource "azurerm_backup_policy_vm" "example" {
  depends_on    = [null_resource.delay]
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  recovery_vault_name = azurerm_recovery_services_vault.named_test_resource.name

  timezone = "UTC"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday"]
    weeks    = ["Last"]
    months   = ["January"]
  }
}

output "resource_aka" {
  value = "azure://${azurerm_backup_policy_vm.example.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_backup_policy_vm.example.id)}"
}

output "id" {
  value = azurerm_backup_policy_vm.example.id
}

output "recovery_vault_name" {
  value = azurerm_recovery_services_vault.named_test_resource.name
}

output "resource_name" {
  value = azurerm_backup_policy_vm.example.name
}
