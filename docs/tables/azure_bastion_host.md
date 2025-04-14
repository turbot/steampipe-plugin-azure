---
title: "Steampipe Table: azure_bastion_host - Query Azure Bastion Hosts using SQL"
description: "Allows users to query Azure Bastion Hosts, providing detailed information about the secure, fully managed network virtual appliance that provides seamless RDP and SSH connectivity to your virtual machines over the Secure Sockets Layer (SSL)."
folder: "Network"
---

# Table: azure_bastion_host - Query Azure Bastion Hosts using SQL

Azure Bastion is a fully managed network virtual appliance that provides seamless RDP and SSH connectivity to your virtual machines over the Secure Sockets Layer (SSL). This service is provisioned directly in your Virtual Network (VNet) and supports all VMs in your VNet using SSL without any exposure through public IP addresses. It enables secure and seamless RDP/SSH connectivity to your virtual machines directly from the Azure portal over SSL.

## Table Usage Guide

The `azure_bastion_host` table provides insights into Azure Bastion Hosts within Microsoft Azure. As a network administrator, explore details about these hosts through this table, including their locations, subnet IDs, and provisioning states. Utilize it to uncover information about hosts, such as their public IP addresses, scale units, and tags, helping ensure secure and seamless connectivity to your virtual machines.

## Examples

### Basic info
Explore which Azure Bastion hosts are being used by checking their provision status and location. This can aid in understanding the distribution of resources and their operational state across different regions and groups.

```sql+postgres
select
  name,
  dns_name,
  provisioning_state,
  region,
  resource_group
from
  azure_bastion_host;
```

```sql+sqlite
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
Determine the areas in which Azure Bastion hosts are not provisioned successfully. This query is useful in identifying and troubleshooting the failed instances, allowing for prompt resolution and minimizing downtime.

```sql+postgres
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

```sql+sqlite
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
This query is useful for identifying the specific subnet details associated with each host within your Azure environment. It can provide valuable insights for network management, helping to understand the distribution of hosts across different subnets.

```sql+postgres
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

```sql+sqlite
select
  h.name as bastion_host_name,
  s.id as subnet_id,
  s.name as subnet_name,
  address_prefix
from
  azure_bastion_host h,
  json_each(h.ip_configurations) ip,
  azure_subnet s
where
  s.id = json_extract(ip.value, '$.properties.subnet.id');
```

### Get IP configuration details associated with each host
This query is used to analyze the IP configuration details associated with each host in the Azure Bastion service. It can help in understanding the allocation method and SKU of each IP configuration, thereby providing insights into the network setup of your Azure resources.

```sql+postgres
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

```sql+sqlite
select
  h.name as bastion_host_name,
  i.name as ip_configuration_name,
  ip_configuration_id,
  ip_address,
  public_ip_allocation_method,
  sku_name as ip_configuration_sku
from
  azure_bastion_host h,
  json_each(ip_configurations) ip,
  azure_public_ip i
where
  i.id = json_extract(ip.value, '$.properties.publicIPAddress.id');
```