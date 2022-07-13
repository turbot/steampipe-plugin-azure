
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
  # version         = "=1.36.0"
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

resource "azurerm_virtual_network" "named_test_resource" {
  name                = var.resource_name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_subnet" "named_test_resource" {
  name                 = var.resource_name
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  virtual_network_name = azurerm_virtual_network.named_test_resource.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  ip_configuration {
    name                          = var.resource_name
    subnet_id                     = azurerm_subnet.named_test_resource.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "named_test_resource" {
  name = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  network_interface_ids = [
    azurerm_network_interface.named_test_resource.id
  ]
  vm_size = "Standard_DS1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer = "UbuntuServer"
    sku = "16.04-LTS"
    version = "latest"
  }

  storage_os_disk {
    name = azurerm_virtual_network.named_test_resource.name
    caching = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_virtual_machine.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_virtual_machine.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_name_upper" {
  value = "${upper(var.resource_name)}"
}

output "resource_id" {
  value = azurerm_virtual_machine.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
