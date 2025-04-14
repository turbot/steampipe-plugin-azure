---
title: "Steampipe Table: azure_compute_resource_sku - Query Azure Compute Resource SKUs using SQL"
description: "Allows users to query Azure Compute Resource SKUs, providing details on available virtual machines, their capabilities, restrictions, and pricing tiers."
folder: "Compute"
---

# Table: azure_compute_resource_sku - Query Azure Compute Resource SKUs using SQL

Azure Compute Resource SKUs represent the purchasable SKUs of Azure virtual machines. These SKUs detail the capabilities, restrictions, and pricing tiers of the available virtual machines. This information is crucial for understanding the options and limitations when deploying Azure virtual machines.

## Table Usage Guide

The `azure_compute_resource_sku` table provides insights into Azure Compute Resource SKUs. As a cloud architect or DevOps engineer, use this table to explore the capabilities, restrictions, and pricing tiers of available Azure virtual machines. Utilize it to make informed decisions on the deployment and scaling of Azure virtual machines based on their SKU details.

## Examples

### Compute resources sku info
Explore the different tiers, sizes, and types of compute resources available in your Azure environment. This can help you understand your options and plan your resource allocation more effectively.

```sql+postgres
select
  name,
  tier,
  size,
  family,
  kind
from
  azure_compute_resource_sku;
```

```sql+sqlite
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
Analyze the settings to understand the capacity range of Azure compute resources. This can help in assessing the scalability of your resources and planning for future capacity needs.

```sql+postgres
select
  name,
  default_capacity,
  maximum_capacity,
  minimum_capacity
from
  azure_compute_resource_sku;
```

```sql+sqlite
select
  name,
  default_capacity,
  maximum_capacity,
  minimum_capacity
from
  azure_compute_resource_sku;
```

### List of all premium type disks and location
Determine the areas in which premium type disks are located to optimize resource management and allocation strategies. This can be particularly useful in identifying potential cost savings or efficiency improvements.

```sql+postgres
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

```sql+sqlite
select
  name,
  resource_type as tier,
  json_each.value as location
from
  azure_compute_resource_sku,
  json_each(locations)
where
  resource_type = 'disks'
  and tier = 'Premium';
```