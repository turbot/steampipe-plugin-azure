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
  default     = "cdffd708-7da0-4cea-abeb-0a4c334d7f64"
  description = "Azure subscription used for the test."
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

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

locals {
  path = "${path.cwd}/info.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [azurerm_resource_group.named_test_resource]
  provisioner "local-exec" {
    command = "az iot dps create --name ${var.resource_name} --resource-group ${var.resource_name} --location ${azurerm_resource_group.named_test_resource.location} > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}


output "resource_aka" {
  depends_on = [null_resource.named_test_resource]
  value      = "azure://${jsondecode(data.local_file.input.content).id}"
}

output "resource_aka_lower" {
  depends_on = [null_resource.named_test_resource]
  value      = "azure://${lower(jsondecode(data.local_file.input.content).id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = jsondecode(data.local_file.input.content).id
}

output "location" {
  value = azurerm_resource_group.named_test_resource.location
}

output "subscription_id" {
  value = var.azure_subscription
}
