---
title: "Steampipe Table: azure_data_protection_backup_job - Query Azure Data Protection Backup Jobs using SQL"
description: "Allows users to query Azure Data Protection Backup Jobs, providing insights into the backup processes and job statuses."
folder: "Data Protection"
---

# Table: azure_data_protection_backup_job - Query Azure Data Protection Backup Jobs using SQL

Azure Data Protection Backup Jobs enable organizations to manage and monitor the backup processes for various data sources. These jobs provide detailed information about the backup instances, job statuses, errors, and performance metrics, which are essential for ensuring data integrity and recovery.

## Table Usage Guide

The `azure_data_protection_backup_job` table provides insights into backup jobs within Microsoft Azure. As an IT Administrator or Data Protection Officer, you can explore details about each backup job, including its status, duration, errors, and associated resources. Use this table to ensure your backup processes are functioning correctly and to quickly identify and resolve any issues.

## Examples

### Basic info
Explore the status and details of your Azure Data Protection Backup Jobs to understand their current state and type. This is beneficial for auditing and managing your backup processes effectively.

```sql+postgres
select
  name,
  id,
  operation,
  operation_category,
  status,
  start_time,
  end_time
from
  azure_data_protection_backup_job;
```

```sql+sqlite
select
  name,
  id,
  operation,
  operation_category,
  status,
  start_time,
  end_time
from
  azure_data_protection_backup_job;
```

### List backup jobs with errors
Identify backup jobs that have encountered errors. This helps in troubleshooting and resolving issues to ensure data protection continuity.

```sql+postgres
select
  name,
  id,
  status,
  jsonb_array_elements_text(error_details) as error_detail
from
  azure_data_protection_backup_job
where
  status = 'Failed';
```

```sql+sqlite
select
  name,
  id,
  status,
  json_each_text(error_details) as error_detail
from
  azure_data_protection_backup_job
where
  status = 'Failed';
```

### List user-triggered backup jobs
Explore the backup jobs that were triggered manually by users. This helps in understanding the frequency and context of adhoc backups.

```sql+postgres
select
  name,
  id,
  is_user_triggered,
  start_time,
  end_time
from
  azure_data_protection_backup_job
where
  is_user_triggered = true;
```

```sql+sqlite
select
  name,
  id,
  is_user_triggered,
  start_time,
  end_time
from
  azure_data_protection_backup_job
where
  is_user_triggered = true;
```

### List backup jobs by policy
Get an overview of backup jobs associated with specific policies. This can assist in ensuring compliance with data protection policies and understanding the effectiveness of different backup strategies.

```sql+postgres
select
  name,
  id,
  policy_name,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  policy_name = 'your_policy_name';
```

```sql+sqlite
select
  name,
  id,
  policy_name,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  policy_name = 'your_policy_name';
```

### List backup jobs by data source
Explore backup jobs based on their data sources. This helps in understanding the backup coverage and identifying any gaps in data protection.

```sql+postgres
select
  name,
  id,
  data_source_name,
  data_source_type,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  data_source_name = 'your_data_source_name';
```

```sql+sqlite
select
  name,
  id,
  data_source_name,
  data_source_type,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  data_source_name = 'your_data_source_name';
```

### List backup jobs by region
Get a regional overview of your backup jobs to ensure that all geographical locations are adequately covered and to identify any regional issues.

```sql+postgres
select
  name,
  id,
  region,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  region = 'your_region';
```

```sql+sqlite
select
  name,
  id,
  region,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  region = 'your_region';
```

### List backup jobs with long duration
Identify backup jobs that have taken longer than expected to complete. This helps in identifying performance bottlenecks and optimizing backup processes.

```sql+postgres
select
  name,
  id,
  duration,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  duration > 'PT1H';  -- Jobs longer than 1 hour
```

```sql+sqlite
select
  name,
  id,
  duration,
  start_time,
  end_time,
  status
from
  azure_data_protection_backup_job
where
  duration > 'PT1H';  -- Jobs longer than 1 hour
```

### List backup jobs by status
Explore backup jobs based on their status to monitor ongoing and completed jobs, and to identify any jobs that may require attention.

```sql+postgres
select
  name,
  id,
  status,
  start_time,
  end_time
from
  azure_data_protection_backup_job
where
  status = 'InProgress';
```

```sql+sqlite
select
  name,
  id,
  status,
  start_time,
  end_time
from
  azure_data_protection_backup_job
where
  status = 'InProgress';
```