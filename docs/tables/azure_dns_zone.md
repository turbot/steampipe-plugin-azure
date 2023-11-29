---
title: "Steampipe Table: azure_dns_zone - Query Azure DNS Zones using SQL"
description: "Allows users to query Azure DNS Zones, providing detailed information about each DNS zone in the Azure account."
---

# Table: azure_dns_zone - Query Azure DNS Zones using SQL

Azure DNS Zones is a service within Microsoft Azure that allows you to host your DNS domain in Azure. It provides name resolution using Microsoft Azure infrastructure, and you can use it to manage and resolve domain names in a virtual network. Azure DNS Zones is globally distributed, highly available, and designed to handle millions of queries per second.

## Table Usage Guide

The 'azure_dns_zone' table delivers comprehensive insights into DNS Zones within Microsoft Azure. As a network administrator, you can leverage this table to explore detailed information about each DNS zone, including its properties, record sets, and associated metadata. The table is particularly useful for understanding the configuration of DNS zones, such as which record sets are associated with each zone, the number of record sets in each zone, and the type of each record set. The schema presents a wide range of attributes of the DNS zone for your analysis, such as the zone name, resource group name, record set count, and associated tags.

## Examples

### Basic info
Explore which resource groups in your Azure DNS Zone are tagged for specific purposes. This allows for efficient management and organization of resources within your network.

```sql
select
  name,
  resource_group,
  tags
from
  azure_dns_zone;
```

### List public DNS zones with record sets
Determine the areas in which public DNS zones have more than one record set in Azure. This can help in understanding the complexity of your DNS configuration and identify potential areas for consolidation or simplification.

```sql
select
  name,
  resource_group
from
  azure_dns_zone
where
  number_of_record_sets > 1;
```

### List public DNS zones with delegated name servers
Explore the public DNS zones that have been delegated to non-Azure name servers, which can be essential in assessing the distribution of your DNS management responsibilities. This query can help identify potential areas of risk or inefficiency in your current DNS management strategy.

```sql
select
  name,
  resource_group,
  ns
from
  azure_dns_zone, jsonb_array_elements_text(name_servers) as ns
where
  zone_type = 'Public'
  and ns not like '%.azure-dns.%.';
```