# Table: azure_backup_job

An Azure Backup job is a task that you can define and run to perform data protection operations on your Azure resources. These jobs are typically used to back up and restore data from various Azure services, such as virtual machines, databases, and files.

**Note**: `vault_name` is required in query parameter to get the job details of the job.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  vault_name,
  region
from
  azure_backup_job
where
  vault_name = 'my-vault';
```

### Get job properties of jobs

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
  azure_backup_job
where
  vault_name = 'my-vault';
```
