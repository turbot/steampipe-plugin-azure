---
title: "Steampipe Table: azure_backup_policy - Query Azure Backup Policies using SQL"
description: "Allows users to query Azure Backup Policies, providing a comprehensive view of all policies that are used within Azure Backup."
---

# Table: azure_backup_policy - Query Azure Backup Policies using SQL

Azure Backup is a cloud-based service offered by Microsoft Azure that allows organizations to back up and restore data across their Azure and non-Azure environments. It provides a way to manage, monitor, and act upon the data protection of infrastructure resources in a scalable and reliable manner. Azure Backup Policies are entities within this service that define the backup schedule, retention policy, and other settings for backup operations.

## Table Usage Guide

The `azure_backup_policy` table provides insights into Azure Backup Policies within Azure Backup. As a DevOps engineer, explore policy-specific details through this table, including the backup schedule, retention settings, and other policy configurations. Utilize it to uncover information about policies, such as their associated backup vaults, and the protected items they are associated with.

## Examples

### Basic info

Analyze the settings to understand the backup schedule and retention policy of your Azure Backup policies. This will help in assessing the data protection strategy of your infrastructure resources.

```sql+postgres
select
  id,
  name,
  resource_group,
  vault_name,
  type
from
  azure_backup_policy;
```

```sql+sqlite
select
  id,
  name,
  resource_group,
  vault_name,
  type
from
  azure_backup_policy;
```

