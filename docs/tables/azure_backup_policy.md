---
title: "Steampipe Table: azure_backup_policy - Query Azure Backup Policies using SQL"
description: "Allows users to query Azure Backup Policies, providing a comprehensive view of all policies that are used within Azure Backup."
folder: "Backup"
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
  name,
  id,
  resource_group,
  vault_name,
  type
from
  azure_backup_policy;
```

```sql+sqlite
select
  name,
  id,
  resource_group,
  vault_name,
  type
from
  azure_backup_policy;
```

### Count the number of policies per vault

Identify the number of policies associated with each backup vault to understand the distribution of policies across your Azure Backup environment.

```sql+postgres
select
  vault_name,
  count(*) as policy_count
from
  azure_backup_policy
group by
  vault_name;
```

```sql+sqlite
select
  vault_name,
  count(*) as policy_count
from
  azure_backup_policy
group by
  vault_name;
```

### List policies with Azure IAAS VM Protection backup configuration

Retrieve the policies that are associated with Azure VM backup configurations to understand the backup settings for your virtual machines.

```sql+postgres
select
  name,
  id,
  resource_group,
  vault_name,
  azure_iaas_vm_protection_policy_property,
  type
from
  azure_backup_policy
where
  azure_iaas_vm_protection_policy_property is not null;
```

```sql+sqlite
select
  name,
  id,
  resource_group,
  vault_name,
  azure_iaas_vm_protection_policy_property,
  type
from
  azure_backup_policy
where
  azure_iaas_vm_protection_policy_property is not null;
```

### List Azure IAAS VM Protection policies with retention duration less than 30 days

Identify the Azure VM backup policies with a retention duration of less than 30 days to ensure that the backup data is retained for a sufficient period.

```sql+postgres
select
  name,
  id,
  resource_group,
  vault_name,
  azure_iaas_vm_protection_policy_property->'retentionPolicy'->'dailySchedule'->>'retentionDuration' as retention_duration_days
from
  azure_backup_policy
where
  azure_iaas_vm_protection_policy_property is not null
  and cast(azure_iaas_vm_protection_policy_property->'retentionPolicy'->'dailySchedule'->'retentionDuration' ->> 'count' as integer) < 30;
```

```sql+sqlite
select
  name,
  id,
  resource_group,
  vault_name,
  azure_iaas_vm_protection_policy_property->'retentionPolicy'->'dailySchedule'->'retentionDuration'->'count' as retention_duration_days
from
  azure_backup_policy
where
  azure_iaas_vm_protection_policy_property is not null
  and cast(azure_iaas_vm_protection_policy_property->'retentionPolicy'->'dailySchedule'->'retentionDuration'->'count' as integer) > 30;
```


### Detailed view of Azure VM Workload backup policies

Gain insights into the specific configurations of backup policies tailored for Azure VM workloads. This query highlights the intricacies of these policies, including backup management type, protection coverage, and the nature of the backup (such as full, incremental, or log backups), along with their scheduling and retention details.

```sql+postgres
 select
  name as policy_name,
  vault_name,
  azure_vm_workload_protection_policy_property->>'backupManagementType' as management_type,
  jsonb_array_length(azure_vm_workload_protection_policy_property->'subProtectionPolicy') as number_of_sub_policies,
  jsonb_pretty(azure_vm_workload_protection_policy_property->'subProtectionPolicy') as sub_protection_policies
from
  azure_backup_policy
where
  azure_vm_workload_protection_policy_property is not null;
```

```sql+sqlite
select
  name as policy_name,
  vault_name,
  json_extract(azure_vm_workload_protection_policy_property, '$.backupManagementType') as management_type,
  json_array_length(json_extract(azure_vm_workload_protection_policy_property, '$.subProtectionPolicy')) as number_of_sub_policies,
  json_pretty(json_extract(azure_vm_workload_protection_policy_property, '$.subProtectionPolicy')) as sub_protection_policies
from
  azure_backup_policy
where
  azure_vm_workload_protection_policy_property is not null;
```
