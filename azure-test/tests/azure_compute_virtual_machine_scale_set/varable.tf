
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
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
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

resource "azurerm_public_ip" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  allocation_method   = "Static"
  domain_name_label   = azurerm_resource_group.named_test_resource.name

  tags = {
    Name = var.resource_name
  }
}

resource "azurerm_lb" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.named_test_resource.id
  }
}

resource "azurerm_lb_backend_address_pool" "named_test_resource" {
  resource_group_name = azurerm_resource_group.named_test_resource.name
  loadbalancer_id     = azurerm_lb.named_test_resource.id
  name                = var.resource_name
}

resource "azurerm_lb_nat_pool" "named_test_resource" {
  resource_group_name            = azurerm_resource_group.named_test_resource.name
  name                           = "ssh"
  loadbalancer_id                = azurerm_lb.named_test_resource.id
  protocol                       = "Tcp"
  frontend_port_start            = 50000
  frontend_port_end              = 50119
  backend_port                   = 22
  frontend_ip_configuration_name = "PublicIPAddress"
}

resource "azurerm_virtual_machine_scale_set" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = 2
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_data_disk {
    lun           = 0
    caching       = "ReadWrite"
    create_option = "Empty"
    disk_size_gb  = 10
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username       = "myadmin"
    admin_password       = "AdminPassword@123"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  network_profile {
    name    = "terraformnetworkprofile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      primary                                = true
      subnet_id                              = azurerm_subnet.named_test_resource.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.named_test_resource.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.named_test_resource.id]
    }
  }

  tags = {
    Name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_virtual_machine_scale_set.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_virtual_machine_scale_set.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_virtual_machine_scale_set.named_test_resource.id
}
