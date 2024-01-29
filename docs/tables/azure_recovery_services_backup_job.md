---
title: "Steampipe Table: azure_recovery_services_backup_job - Query Azure Recovery Services Backup Jobs using SQL"
description: "Allows users to query Azure Recovery Services Backup Jobs, specifically the job status, duration, and details, providing insights into backup operations and their outcomes."
---

# Table: azure_recovery_services_backup_job - Query Azure Recovery Services Backup Jobs using SQL

Azure Recovery Services is a service within Microsoft Azure that provides data backup and disaster recovery capabilities. It allows you to protect and recover your data in the Microsoft cloud, providing a simple, secure, and cost-effective solution for protecting your data and maintaining business continuity. The service supports backup and recovery for Azure VMs, SQL Server, Azure SQL Database, on-premises Windows Servers, and more.

## Table Usage Guide

The `azure_recovery_services_backup_job` table provides insights into backup jobs within Azure Recovery Services. As a system administrator or a backup operator, explore job-specific details through this table, including job status, duration, and details. Utilize it to monitor the status of backup operations, identify any issues, and ensure the successful completion of backup jobs.

## Examples

### Basic info
Determine the areas in which specific Azure recovery services backup jobs are performed, focusing on a specific vault. This can help to assess the distribution and management of backup jobs within your Azure environment.

```sql+postgres
select
  name,
  id,
  type,
  vault_name,
  region
from
  azure_recovery_services_backup_job
where
  vault_name = 'my-vault';
```

```sql+sqlite
select
  name,
  id,
  type,
  vault_name,
  region
from
  azure_recovery_services_backup_job
where
  vault_name = 'my-vault';
```

### Get job properties of jobs
Discover the specifics of job properties in your backup jobs. This can help you understand the type of job, its management, operational status, and timing details, providing critical insights into your backup operations.

```sql+postgres
select
  name,
  id,
  properties ->> 'JobType' as job_type,
  properties ->> 'ActivityID' as activity_id,
  properties ->> 'BackupManagementType' as backup_management_type,
  properties ->> 'EndTime' as end_time,
  properties ->> 'EntityFriendlyName' as entity_friendly_name,
  properties ->> 'Operation' as Operation,
  properties ->> 'StartTime' as start_time,
  properties ->> 'Status' as Status
from
  azure_recovery_services_backup_job;
```

```sql+sqlite
select
  name,
  id,
  json_extract(properties, '$.JobType') as job_type,
  json_extract(properties, '$.ActivityID') as activity_id,
  json_extract(properties, '$.BackupManagementType') as backup_management_type,
  json_extract(properties, '$.EndTime') as end_time,
  json_extract(properties, '$.EntityFriendlyName') as entity_friendly_name,
  json_extract(properties, '$.Operation') as Operation,
  json_extract(properties, '$.StartTime') as start_time,
  json_extract(properties, '$.Status') as Status
from
  azure_recovery_services_backup_job;
```