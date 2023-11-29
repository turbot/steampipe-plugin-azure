---
title: "Steampipe Table: azure_resource_group - Query Azure Resource Groups using SQL"
description: "Allows users to query Azure Resource Groups."
---

# Table: azure_resource_group - Query Azure Resource Groups using SQL

Azure Resource Groups are basic units in Azure that provide a logical grouping for resources deployed on Azure. These groups hold related resources for an Azure solution. Resource groups are used to manage and organize Azure resources so you can monitor, control access, provision and manage billing.

## Table Usage Guide

The 'azure_resource_group' table provides insights into Resource Groups within Azure. As a DevOps engineer, explore Resource Group-specific details through this table, including locations, managed_by details, and associated metadata. Utilize it to uncover information about Resource Groups, such as those with specific provisioning states, the tags associated with each group, and the time they were last updated. The schema presents a range of attributes of the Resource Group for your analysis, like the group ID, name, type, and associated tags.

## Examples

### List of resource groups with their locations
Explore which resource groups are located in different regions to better manage and organize your resources in Azure. This can help streamline operations and ensure resources are optimally allocated across various geographical locations.

```sql
select
  name,
  region 
from
  azure_resource_group;
```

### List of resource groups without owner tag key
Identify the Azure resource groups that lack an 'owner' tag. This is useful for pinpointing potential areas of unaccountability or mismanagement within your resources.

```sql
select
  name,
  tags
from
  azure_resource_group
where
  not tags :: JSONB ? 'owner';
```