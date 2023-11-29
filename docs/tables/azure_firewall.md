---
title: "Steampipe Table: azure_firewall - Query Azure Network Firewalls using SQL"
description: "Allows users to query Azure Network Firewalls for detailed information about their configuration, status, rules, and more."
---

# Table: azure_firewall - Query Azure Network Firewalls using SQL

Azure Firewall is a managed, cloud-based network security service that protects your Azure Virtual Network resources. It's a fully stateful firewall as a service with built-in high availability and unrestricted cloud scalability. You can centrally create, enforce, and log application and network connectivity policies across subscriptions and virtual networks.

## Table Usage Guide

The 'azure_firewall' table provides insights into Network Firewalls within Azure Networking. As a network administrator, explore firewall-specific details through this table, including network rules, application rules, and associated metadata. Utilize it to uncover information about firewalls, such as rules with broad coverage, the relationships between different rules, and the verification of application rules. The schema presents a range of attributes of the Network Firewall for your analysis, like the firewall name, resource group, subscription ID, and associated tags.

## Examples

### Azure firewall location and availability zone count info
Explore the distribution of Azure firewalls across different regions and gain insights into their availability zone count to optimize network security and resource allocation.

```sql
select
  name,
  region,
  jsonb_array_length(availability_zones) availability_zones_count
from
  azure_firewall;
```

### Basic IP configuration info
Determine the configuration of IP addresses in your Azure firewall. This query allows you to identify private and public IP addresses, their allocation methods, and the virtual network they're associated with, helping you maintain an overview of your network's structure and security.

```sql
select
  name,
  ip #> '{properties, privateIPAddress}' private_ip_address,
  ip #> '{properties, privateIPAllocationMethod}' private_ip_allocation_method,
  split_part(
    ip -> 'properties' -> 'publicIPAddress' ->> 'id',
    '/',
    9
  ) public_ip_address_id,
  split_part(ip -> 'properties' ->> 'subnet', '/', 9) virtual_network
from
  azure_firewall
  cross join jsonb_array_elements(ip_configurations) as ip;
```

### List the premium category firewalls
Explore which firewalls fall under the premium category in Azure. This is beneficial for assessing your current security infrastructure and planning future upgrades or budget allocation.

```sql
select
  name,
  sku_tier,
  sku_name
from
  azure_firewall
where
  sku_tier = 'Premium';
```

### List of firewalls where threat intel mode is off
Discover the segments that have their firewall's threat intelligence mode turned off. This could be useful for identifying potential security gaps in your Azure services.

```sql
select
  name,
  threat_intel_mode
from
  azure_firewall
where
  threat_intel_mode = 'Off';
```