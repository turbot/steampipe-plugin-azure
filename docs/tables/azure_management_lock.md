---
title: "Steampipe Table: azure_management_lock - Query Azure Management Locks using SQL"
description: "Allows users to query Azure Management Locks, particularly their properties and associated resources, providing insights into the locks' configurations and status."
folder: "Resource Manager"
---

# Table: azure_management_lock - Query Azure Management Locks using SQL

Azure Management Lock is a feature within Microsoft Azure that helps prevent accidental deletion or modification of Azure resources. It allows administrators to apply a 'CanNotDelete' or 'ReadOnly' lock on a subscription, resource group, or resource to protect it from being inadvertently deleted or modified. These locks can be used across various Azure resources, including virtual machines, storage accounts, and more.

## Table Usage Guide

The `azure_management_lock` table provides insights into Management Locks within Microsoft Azure. As an Azure administrator or a DevOps engineer, explore lock-specific details through this table, including their level (CanNotDelete or ReadOnly), scope, and associated resources. Utilize it to uncover information about locks, such as those applied on critical resources, to ensure their accidental deletion or modification is prevented.

## Examples

### List of resources where the management locks are applied
This example demonstrates how to identify resources that have management locks applied to them within the Azure environment. This could be useful for administrators who need to manage access controls or troubleshoot issues related to locked resources.

```sql+postgres
select
  name,
  split_part(id, '/', 8) as resource_type,
  split_part(id, '/', 9) as resource_name
from
  azure_management_lock;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```


### Resources and lock levels
Uncover the details of specific Azure resources and their associated lock levels. This can help you assess what resources are locked at what level, aiding in security and access management.

```sql+postgres
select
  name,
  split_part(id, '/', 8) as resource_type,
  split_part(id, '/', 9) as resource_name,
  lock_level
from
  azure_management_lock;
```

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```