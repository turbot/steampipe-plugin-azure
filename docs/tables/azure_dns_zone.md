---
title: "Steampipe Table: azure_dns_zone - Query Azure DNS Zones using SQL"
description: "Allows users to query Azure DNS Zones, providing detailed information about each DNS zone within the Azure environment."
---

# Table: azure_dns_zone - Query Azure DNS Zones using SQL

Azure DNS is a hosting service for DNS domains, providing name resolution using Microsoft Azure infrastructure. By hosting domains in Azure, it provides you with the same reliability and performance provided to Microsoft’s global network. Azure DNS also supports private DNS domains.

## Table Usage Guide

The `azure_dns_zone` table provides insights into DNS zones within Microsoft Azure. As a network administrator, explore DNS zone-specific details through this table, including record sets, number of record sets, and associated metadata. Utilize it to uncover information about DNS zones, such as those with certain properties, the relationships between different zones, and the verification of DNS settings.

## Examples

### Basic info
This query allows you to analyze the configuration of your Azure DNS zones. It helps you identify instances where specific tags are used, providing insights into the organization and management of your resources.

```sql+postgres
select
  name,
  resource_group,
  tags
from
  azure_dns_zone;
```

```sql+sqlite
select
  name,
  resource_group,
  tags
from
  azure_dns_zone;
```

### List public DNS zones with record sets
Explore which public DNS zones in your Azure environment contain more than one record set. This can help in managing and organizing your DNS records effectively.

```sql+postgres
select
  name,
  resource_group
from
  azure_dns_zone
where
  number_of_record_sets > 1;
```

```sql+sqlite
select
  name,
  resource_group
from
  azure_dns_zone
where
  number_of_record_sets > 1;
```

### List public DNS zones with delegated name servers
Determine the areas in which public DNS zones are utilizing delegated name servers, which can be beneficial for those seeking to manage or troubleshoot their DNS configurations.

```sql+postgres
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

```sql+sqlite
select
  name,
  resource_group,
  ns.value as ns
from
  azure_dns_zone,
  json_each(name_servers) as ns
where
  zone_type = 'Public'
  and ns.value not like '%.azure-dns.%.';
```