
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
  version = "=2.44.0"
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

resource "local_file" "python_file" {
  filename          = "${path.cwd}/../../test.py"
  sensitive_content = "def test (event, context):\n\tprint ('This is a test for integration testing to check creation of a lambda function')"
}

data "archive_file" "zip" {
  depends_on  = [local_file.python_file]
  type        = "zip"
  source_file = "${path.cwd}/../../test.py"
  output_path = "${path.cwd}/../../test.zip"
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
  account_kind             = "StorageV2"
  access_tier              = "Cool"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "named_test_resource" {
  name                  = var.resource_name
  storage_account_name  = azurerm_storage_account.named_test_resource.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "named_test_resource" {
  name                   = var.resource_name
  storage_account_name   = azurerm_storage_account.named_test_resource.name
  storage_container_name = azurerm_storage_container.named_test_resource.name
  type                   = "Block"
  source                 = "${path.cwd}/../../test.zip"
}

output "resource_name" {
  value = var.resource_name
}

output "subscription_id" {
  value = var.azure_subscription
}

output "resource_id" {
  value = azurerm_storage_blob.named_test_resource.id
}
