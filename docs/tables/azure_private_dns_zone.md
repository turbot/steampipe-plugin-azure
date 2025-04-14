---
title: "Steampipe Table: azure_private_dns_zone - Query Azure Private DNS Zones using SQL"
description: "Allows users to query Azure Private DNS Zones, specifically providing insights into their properties, records, and associated metadata."
folder: "DNS"
---

# Table: azure_private_dns_zone - Query Azure Private DNS Zones using SQL

Azure Private DNS Zone is a service within Microsoft Azure that allows you to use your own domain name, rather than the Azure-provided names. It provides a simple, reliable, secure DNS service to manage and resolve domain names in a Virtual Network without the need for custom DNS solutions. Azure Private DNS Zones helps you to customize domain names for Azure services, manage DNS records, and improve network security.

## Table Usage Guide

The `azure_private_dns_zone` table provides insights into Azure Private DNS Zones within Azure DNS. As a network administrator, explore zone-specific details through this table, including properties, records, and associated metadata. Utilize it to uncover information about zones, such as their status, the number of records, and their associated virtual networks.

## Examples

### Basic info
Explore which private DNS zones are present in your Azure environment and determine the associated resource groups and tags for better resource management and categorization.

```sql+postgres
select
  name,
  resource_group,
  tags
from
  azure_private_dns_zone;
```

```sql+sqlite
select
  name,
  resource_group,
  tags
from
  azure_private_dns_zone;
```

### List private DNS zones with record sets
Identify private DNS zones in Azure that have more than one record set. This is useful for managing and organizing DNS resources efficiently.

```sql+postgres
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_record_sets > 1;
```

```sql+sqlite
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_record_sets > 1;
```

### List private DNS zones linked to no virtual networks
Explore which private DNS zones in Azure are not linked to any virtual networks. This can be useful in identifying potential areas of network optimization or redundancy elimination.

```sql+postgres
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_virtual_network_links_with_registration = 0;
```

```sql+sqlite
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_virtual_network_links_with_registration = 0;
```