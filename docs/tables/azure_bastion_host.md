---
title: "Steampipe Table: azure_bastion_host - Query Azure Bastion Hosts using SQL"
description: "Allows users to query Azure Bastion Hosts to retrieve information about the state, configurations, and associated resources."
---

# Table: azure_bastion_host - Query Azure Bastion Hosts using SQL

Azure Bastion is a fully managed network security service that provides secure and seamless RDP and SSH access to your virtual machines directly from the Azure portal. Azure Bastion is provisioned directly in your Virtual Network (VNet) and supports all VMs in your VNet. Using Azure Bastion protects your virtual machines from exposing RDP/SSH ports to the outside world, while providing secure access to manage your VMs.

## Table Usage Guide

The 'azure_bastion_host' table provides insights into Bastion Hosts within Azure Bastion service. As an IT administrator, explore host-specific details through this table, including its state, configurations, and associated resources. Utilize it to uncover information about hosts, such as those with specific configurations, the associated subnets, and the verification of their state. The schema presents a range of attributes of the Bastion Host for your analysis, like the host name, provisioning state, type, id, and associated tags.

## Examples

### Basic info
Explore which Azure Bastion Hosts are currently provisioned and where they are located. This helps in managing resources and planning deployment strategies across different regions.

```sql
select
  name,
  dns_name,
  provisioning_state,
  region,
  resource_group
from
  azure_bastion_host;
```

### List bastion hosts that are in failed state
Discover the segments that have bastion hosts in a failed state. This can help in identifying and troubleshooting problematic hosts, ensuring the stability and security of your Azure environment.

```sql
select
  name,
  dns_name,
  provisioning_state,
  region,
  resource_group
from
  azure_bastion_host
where
  provisioning_state = 'Failed';
```

### Get subnet details associated with each host
Discover the network organization of your Azure resources by identifying the specific subnets associated with each bastion host. This allows for efficient infrastructure management and helps in identifying potential network vulnerabilities.

```sql
select
  h.name as bastion_host_name,
  s.id as subnet_id,
  s.name as subnet_name,
  address_prefix
from
  azure_bastion_host h,
  jsonb_array_elements(ip_configurations) ip,
  azure_subnet s
where
  s.id = ip -> 'properties' -> 'subnet' ->> 'id';
```

### Get IP configuration details associated with each host
Explore the IP configurations linked to each host in your Azure environment to gain insights into allocation methods and SKU details. This can help in managing and optimizing your network resources in Azure.

```sql
select
  h.name as bastion_host_name,
  i.name as ip_configuration_name,
  ip_configuration_id,
  ip_address,
  public_ip_allocation_method,
  sku_name as ip_configuration_sku
from
  azure_bastion_host h,
  jsonb_array_elements(ip_configurations) ip,
  azure_public_ip i
where
  i.id = ip -> 'properties' -> 'publicIPAddress' ->> 'id';
```