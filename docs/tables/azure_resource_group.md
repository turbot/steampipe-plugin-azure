---
title: "Steampipe Table: azure_resource_group - Query Azure Resource Groups using SQL"
description: "Allows users to query Azure Resource Groups, providing detailed information about the group's properties, location, and tags."
folder: "Resource"
---

# Table: azure_resource_group - Query Azure Resource Groups using SQL

Azure Resource Groups are essential components of Azure Resource Management, serving as logical containers for resources deployed within an Azure subscription. They provide a way to monitor, control access, provision and manage billing for collections of assets, which are required to run an application, or used by a department or a team. Azure Resource Groups offer a means to manage and organize resources based on lifecycles and application architecture, along with providing access control, consistency, and efficiency.

## Table Usage Guide

The `azure_resource_group` table provides insights into Resource Groups within Microsoft Azure. As a DevOps engineer, explore group-specific details through this table, including properties, location, and associated tags. Utilize it to manage and organize resources, control access, and manage billing for collections of assets used by applications, departments, or teams.

## Examples

### List of resource groups with their locations
Explore which Azure resource groups are located in specific regions to better manage and organize your resources. This is useful for understanding the geographical distribution of your resources for efficiency and cost-effectiveness.

```sql+postgres
select
  name,
  region 
from
  azure_resource_group;
```

```sql+sqlite
select
  name,
  region 
from
  azure_resource_group;
```

### List of resource groups without owner tag key
Explore which Azure resource groups are missing an 'owner' tag. This query assists in identifying and addressing gaps in resource ownership, aiding in resource management and accountability.

```sql+postgres
select
  name,
  tags
from
  azure_resource_group
where
  not tags :: JSONB ? 'owner';
```

```sql+sqlite
select
  name,
  tags
from
  azure_resource_group
where
  json_extract(tags, '$.owner') is null;
```