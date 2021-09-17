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
    command = "az disk-access create -g ${var.resource_name} -l ${azurerm_resource_group.named_test_resource.location} -n ${var.resource_name} > ${local.path}"
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
  value = jsondecode(data.local_file.input.content).name
}

output "resource_id" {
  value = jsondecode(data.local_file.input.content).id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "resource_name_upper" {
  value = upper(jsondecode(data.local_file.input.content).name)
}
