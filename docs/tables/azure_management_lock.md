---
title: "Steampipe Table: azure_management_lock - Query Azure Management Locks using SQL"
description: "Allows users to query Azure Management Locks."
---

# Table: azure_management_lock - Query Azure Management Locks using SQL

Azure Management Locks are a feature in Microsoft Azure that provides the ability to apply a lock with different levels of access control to any resource in Azure. These locks help prevent accidental deletion or modification of your Azure resources. Management Locks can be applied to resource groups, subscriptions, and individual resources, providing a flexible and robust mechanism for safeguarding your critical cloud resources.

## Table Usage Guide

The 'azure_management_lock' table provides insights into Management Locks within Microsoft Azure. As a DevOps engineer, explore lock-specific details through this table, including the lock level, notes, and owners. Utilize it to uncover information about locks, such as those with 'CanNotDelete' or 'ReadOnly' access levels, the resources associated with each lock, and the lock's owners. The schema presents a range of attributes of the Management Lock for your analysis, like the lock name, id, type, and associated tags.

## Examples

### List of resources where the management locks are applied
Determine the areas in which management locks are applied within Azure resources. This query is beneficial for understanding where your resources are secured, helping to maintain and enhance your security posture.

```sql
select
  name,
  split_part(id, '/', 8) as resource_type,
  split_part(id, '/', 9) as resource_name
from
  azure_management_lock;
```


### Resources and lock levels
Explore which resources in your Azure Management are locked and the level of these locks. This can help in understanding the security measures in place and aid in managing resource accessibility.

```sql
select
  name,
  split_part(id, '/', 8) as resource_type,
  split_part(id, '/', 9) as resource_name,
  lock_level
from
  azure_management_lock;
```