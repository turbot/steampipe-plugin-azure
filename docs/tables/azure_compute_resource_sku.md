---
title: "Steampipe Table: azure_compute_resource_sku - Query Azure Compute Resource SKUs using SQL"
description: "Allows users to query Azure Compute Resource SKUs"
---

# Table: azure_compute_resource_sku - Query Azure Compute Resource SKUs using SQL

Azure Compute Resource SKUs represent the purchasable units for Azure resources, providing details about the available resources for a subscription, including their tier, size, and corresponding cost. They are used to define the size and capacity of the resources that you can provision within your Azure subscription. Each SKU represents a specific combination of resource type, tier, and size.

## Table Usage Guide

The 'azure_compute_resource_sku' table provides insights into the available SKUs for Azure Compute Resources. As a DevOps engineer, explore SKU-specific details through this table, including their tier, size, and corresponding cost. Utilize it to uncover information about SKUs, such as their capacity, family, kind, and locations. The schema presents a range of attributes of the SKU for your analysis, like the resource type, tier, size, and restrictions.

## Examples

### Compute resources sku info
Determine the characteristics of your Azure compute resources, such as their tier, size, and family. This is useful for understanding the specifics of your current resources and can aid in planning future resource allocation or optimization.

```sql
select
  name,
  tier,
  size,
  family,
  kind
from
  azure_compute_resource_sku;
```


### Azure compute resources and their capacity
Identify the capacity range of Azure compute resources to efficiently manage and allocate your cloud resources.

```sql
select
  name,
  default_capacity,
  maximum_capacity,
  minimum_capacity
from
  azure_compute_resource_sku;
```


### List of all premium type disks and location
Explore which premium type disks are in use and their locations. This is useful to manage resources and understand their distribution across various locations.

```sql
select
  name,
  resource_type tier,
  l as location
from
  azure_compute_resource_sku,
  jsonb_array_elements_text(locations) as l
where
  resource_type = 'disks'
  and tier = 'Premium';
```