---
title: "Steampipe Table: azure_resource_link - Query Azure Resource Links using SQL"
description: "Allows users to query Azure Resource Links, providing insights into the interconnections between various Azure resources."
folder: "Resource Manager"
---

# Table: azure_resource_link - Query Azure Resource Links using SQL

Azure Resource Links is a feature within Microsoft Azure that allows you to create and manage links between resources. These links can be used for organizing resources and defining dependencies between them. Azure Resource Links helps you to understand the relationships and dependencies between your Azure resources.

## Table Usage Guide

The `azure_resource_link` table provides insights into the interconnections between various Azure resources. As a cloud architect or a DevOps engineer, you can explore link-specific details through this table, including the source and target of each link, and the properties of the link. Utilize it to uncover information about resource dependencies, such as those with circular dependencies or orphaned resources, and to aid in resource management and organization.

## Examples

### Basic Info
Discover the segments that connect different resources in your Azure environment. This query is particularly useful for understanding the relationships and dependencies between your resources, aiding in efficient resource management and troubleshooting.

```sql+postgres
select
  name,
  id,
  type,
  source_id,
  target_id
from
  azure_resource_link;
```

```sql+sqlite
select
  name,
  id,
  type,
  source_id,
  target_id
from
  azure_resource_link;
```

### List resource links with virtual machines
Determine the areas in which resources are linked with virtual machines in your Azure environment. This can be useful for managing and understanding dependencies between your resources.

```sql+postgres
select
  name,
  id,
  source_id,
  target_id
from
  azure_resource_link
where
  source_id LIKE '%virtualmachines%';
```

```sql+sqlite
select
  name,
  id,
  source_id,
  target_id
from
  azure_resource_link
where
  source_id LIKE '%virtualmachines%';
```