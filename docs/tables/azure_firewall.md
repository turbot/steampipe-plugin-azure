---
title: "Steampipe Table: azure_firewall - Query Azure Firewalls using SQL"
description: "Allows users to query Azure Firewalls, providing insights into the configuration, status, and rules of each firewall in the Azure ecosystem."
---

# Table: azure_firewall - Query Azure Firewalls using SQL

Azure Firewall is a managed, cloud-based network security service that protects your Azure Virtual Network resources. It's a fully stateful firewall as a service with built-in high availability and unrestricted cloud scalability. Azure Firewall can centrally create, enforce, and log application and network connectivity policies across subscriptions and virtual networks.

## Table Usage Guide

The `azure_firewall` table provides insights into the firewalls within Azure. As a security engineer, explore firewall-specific details through this table, including rules, configurations, and associated metadata. Utilize it to uncover information about firewalls, such as their current status, applied rules, and the verification of connectivity policies.

## Examples

### Azure firewall location and availability zone count info
Analyze the number of availability zones for each Azure firewall and their respective regions to manage and optimize your resource distribution effectively. This can help in improving your application's resilience and availability across different regions.

```sql+postgres
select
  name,
  region,
  jsonb_array_length(availability_zones) availability_zones_count
from
  azure_firewall;
```

```sql+sqlite
select
  name,
  region,
  json_array_length(availability_zones) as availability_zones_count
from
  azure_firewall;
```

### Basic IP configuration info
This query helps you analyze your Azure firewall's IP configuration. By running this, you can gain insights into details like private IP address, allocation method, associated public IP address ID, and the virtual network it is part of, which can be crucial for network management and security purposes.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```

### List the premium category firewalls
Discover the segments that are using premium category firewalls in your Azure environment. This can help you understand where higher levels of security have been implemented.

```sql+postgres
select
  name,
  sku_tier,
  sku_name
from
  azure_firewall
where
  sku_tier = 'Premium';
```

```sql+sqlite
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
Determine the areas in your Azure network where your firewalls are potentially vulnerable due to the threat intelligence mode being turned off. This can help enhance your network security by identifying and rectifying these weak points.

```sql+postgres
select
  name,
  threat_intel_mode
from
  azure_firewall
where
  threat_intel_mode = 'Off';
```

```sql+sqlite
select
  name,
  threat_intel_mode
from
  azure_firewall
where
  threat_intel_mode = 'Off';
```