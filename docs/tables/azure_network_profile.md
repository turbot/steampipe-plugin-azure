---
title: "Steampipe Table: azure_network_profile - Query Azure Network Profiles using SQL"
description: "Allows users to query Azure Network Profiles, providing access to network configuration details, subnet assignments, and container networking settings."
folder: "Network"
---

# Table: azure_network_profile - Query Azure Network Profiles using SQL

Azure Network Profile is a network configuration template for Azure resources, particularly used with Azure Container Instances. It specifies network properties for resources, such as the subnet into which they should be deployed. Network profiles are essential when deploying container groups to virtual networks, enabling secure communication between containers and other resources within the virtual network.

## Table Usage Guide

The `azure_network_profile` table provides insights into network profiles within Azure. As a network administrator or developer, you can use this table to examine network profile configurations, understand which subnets are associated with container deployments, and verify the network settings for containerized applications. This is valuable for security assessments, network infrastructure planning, and ensuring proper network isolation for containerized workloads.

## Examples

### Basic network profile information

Explore the basic attributes of your Azure Network Profiles to understand their configuration and associated container networks.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  resource_group
from
  azure_network_profile;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  resource_group
from
  azure_network_profile;
```

### List network profiles by their provisioning state

Determine the areas in which network profiles have different provisioning states to identify profiles that may require attention.

```sql+postgres
select
  name,
  id,
  provisioning_state,
  resource_group
from
  azure_network_profile
order by
  provisioning_state;
```

```sql+sqlite
select
  name,
  id,
  provisioning_state,
  resource_group
from
  azure_network_profile
order by
  provisioning_state;
```

### Get container network interface configuration details

Examine container network interface configuration to understand how containers are connected to the network, including subnet associations and IP assignment methods.

```sql+postgres
select
  name,
  c->>'name' as interface_name,
  c->'properties'->'ipConfigurations'->0->'properties'->'subnet'->>'id' as subnet_id,
  c->'properties'->'ipConfigurations'->0->>'name' as ip_config_name
from
  azure_network_profile,
  jsonb_array_elements(container_network_interface_configurations) as c;
```

```sql+sqlite
select
  name,
  json_extract(c.value, '$.name') as interface_name,
  json_extract(json_extract(json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations'), '$[0]'), '$.properties'), '$.subnet.id') as subnet_id,
  json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations'), '$[0].name') as ip_config_name
from
  azure_network_profile,
  json_each(container_network_interface_configurations) as c;
```

### Find network profiles associated with specific virtual networks

Identify which network profiles are connected to specific virtual networks to better manage your network infrastructure.

```sql+postgres
select
  name,
  id,
  resource_group,
  c->'properties'->'ipConfigurations'->0->'properties'->'subnet'->>'id' as subnet_id,
  split_part(c->'properties'->'ipConfigurations'->0->'properties'->'subnet'->>'id', '/', 9) as virtual_network_name
from
  azure_network_profile,
  jsonb_array_elements(container_network_interface_configurations) as c
where
  c->'properties'->'ipConfigurations'->0->'properties'->'subnet'->>'id' is not null;
```

```sql+sqlite
select
  p.name,
  p.id,
  p.resource_group,
  json_extract(json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations[0]'), '$.properties'), '$.subnet.id') as subnet_id,
  substr(
    json_extract(json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations[0]'), '$.properties'), '$.subnet.id'),
    instr(json_extract(json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations[0]'), '$.properties'), '$.subnet.id'), 'virtualNetworks/') + 16,
    instr(substr(json_extract(json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations[0]'), '$.properties'), '$.subnet.id'),
      instr(json_extract(json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations[0]'), '$.properties'), '$.subnet.id'), 'virtualNetworks/') + 16), '/') - 1
  ) as virtual_network_name
from
  azure_network_profile p,
  json_each(p.container_network_interface_configurations) as c
where
  json_extract(json_extract(json_extract(json_extract(c.value, '$.properties'), '$.ipConfigurations[0]'), '$.properties'), '$.subnet.id') is not null;
```
