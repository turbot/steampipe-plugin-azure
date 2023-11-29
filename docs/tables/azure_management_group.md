---
title: "Steampipe Table: azure_management_group - Query Azure Management Groups using SQL"
description: "Allows users to query Azure Management Groups"
---

# Table: azure_management_group - Query Azure Management Groups using SQL

Azure Management Groups provide a way to manage access, policies, and compliance across multiple Azure subscriptions. They offer the flexibility to manage the details of Azure resources, such as Azure subscriptions and policies, at a high level. Management groups are containers for managing access, policies, and compliance across multiple subscriptions.

## Table Usage Guide

The 'azure_management_group' table provides insights into Management Groups within Azure. As a DevOps engineer, explore group-specific details through this table, including group IDs, names, types, and associated metadata. Utilize it to uncover information about groups, such as the parent and children of each group, and the level of each group in the hierarchy. The schema presents a range of attributes of the Management Group for your analysis, like the group ID, name, type, and associated tags.

## Examples

### Basic info
Explore the basic details of your Azure Management Groups to understand their types and update history. This information is useful for assessing your current Azure configurations and identifying any necessary changes.

```sql
select
  id,
  name,
  type,
  tenant_id,
  updated_by
from
  azure_management_group;
```

### List children for management groups
This query is used to examine the hierarchical structure of management groups within an Azure environment. It provides insights into which groups are nested within others and who last updated them, helping to understand the organization's resource management structure.

```sql
select
  name,
  updated_by,
  jsonb_pretty(children) as children
from
  azure_management_group;
```

### List parent details for management groups
Explore the details of parent groups within the management hierarchy to understand who made the most recent updates. This can be useful for tracking changes and maintaining organizational structure in Azure.

```sql
select
  name,
  updated_by,
  jsonb_pretty(parent) as parent
from
  azure_management_group;
```