---
title: "Steampipe Table: azure_management_group - Query Azure Management Groups using SQL"
description: "Allows users to query Azure Management Groups, providing a hierarchical structure for unified policy and access management across multiple Azure subscriptions."
folder: "Management Group"
---

# Table: azure_management_group - Query Azure Management Groups using SQL

Azure Management Groups offer a level of scope above subscriptions. They provide a hierarchical structure for unified policy and access management across multiple Azure subscriptions. Management groups allow you to organize subscriptions into containers called "management groups" and apply your governance conditions to the management groups.

## Table Usage Guide

The `azure_management_group` table provides insights into Management Groups within Azure. As a system administrator or a DevOps engineer, explore group-specific details through this table, including group hierarchy, subscription associations, and associated metadata. Utilize it to uncover information about groups, such as their structure, the subscriptions they contain, and the policies applied to them.

**Important notes:**
- You need to have at least read access to the specific management group to query this table.

## Examples

### Basic info
Explore the management groups within your Azure environment to understand their types and the tenants they belong to. This can help in identifying who last updated these groups, aiding in accountability and tracking changes.

```sql+postgres
select
  id,
  name,
  type,
  tenant_id,
  updated_by
from
  azure_management_group;
```

```sql+sqlite
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
Explore the updated information of Azure Management Groups, including the associated children groups. This is useful for understanding the hierarchical structure and changes made within your Azure Management Groups.

```sql+postgres
select
  name,
  updated_by,
  jsonb_pretty(children) as children
from
  azure_management_group;
```

```sql+sqlite
select
  name,
  updated_by,
  children
from
  azure_management_group;
```

### List parent details for management groups
Explore which management groups in Azure have been recently modified and by whom. This can provide insights into changes in the organizational structure and help maintain accountability.

```sql+postgres
select
  name,
  updated_by,
  jsonb_pretty(parent) as parent
from
  azure_management_group;
```

```sql+sqlite
select
  name,
  updated_by,
  parent
from
  azure_management_group;
```