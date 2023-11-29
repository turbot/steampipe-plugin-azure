---
title: "Steampipe Table: azure_recovery_services_backup_job - Query Azure Recovery Services Backup Jobs using SQL"
description: "Allows users to query Azure Recovery Services Backup Jobs."
---

# Table: azure_recovery_services_backup_job - Query Azure Recovery Services Backup Jobs using SQL

Azure Recovery Services is a service within Microsoft Azure that provides data protection and disaster recovery capabilities. It enables backup and restore functionalities for Azure Virtual Machines, SQL workloads, and on-premises Windows Servers. Azure Recovery Services helps maintain data availability and ensures business continuity during planned and unplanned outages.

## Table Usage Guide

The 'azure_recovery_services_backup_job' table provides insights into backup jobs within Azure Recovery Services. As a DevOps engineer, explore job-specific details through this table, including job status, start and end times, and associated metadata. Utilize it to uncover information about jobs, such as those with errors, the duration of jobs, and the verification of backup items. The schema presents a range of attributes of the backup job for your analysis, like the job ID, backup management type, duration, and associated tags.

## Examples

### Basic info
Analyze the settings to understand the specifics of backup jobs in a particular Azure recovery services vault. This can help in evaluating the backup strategy and ensuring data recovery measures are in line with your organization's policies.

```sql
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
Explore the specifics of different jobs, such as the type, associated activities, management methods, and operational status. This can provide insights into job performance and help identify areas for optimization.

```sql
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