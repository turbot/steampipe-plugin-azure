---
title: "Steampipe Table: azure_mssql_managed_instance - Query Azure SQL Managed Instances using SQL"
description: "Allows users to query Azure SQL Managed Instances, providing a comprehensive view of the settings, configurations, and health status of these resources."
---

# Table: azure_mssql_managed_instance - Query Azure SQL Managed Instances using SQL

Azure SQL Managed Instance is a fully managed SQL Server Database Engine hosted in Azure cloud. It provides most of the SQL Server's capabilities, allowing you to migrate SQL server workloads to Azure with minimal changes. With built-in intelligence that learns your unique database patterns and adaptive performance tuning based on AI, SQL Managed Instance is a best-in-class database service.

## Table Usage Guide

The `azure_mssql_managed_instance` table provides insights into Azure SQL Managed Instances. As a DBA or a cloud architect, you can explore specific details about these instances, including their settings, configurations, and health status. Use this table to monitor and manage your SQL instances effectively, ensuring optimal performance and resource utilization.

## Examples

### Basic info
Explore the status and security settings of your managed instances in Azure's SQL service. This query is useful for understanding the licensing and encryption standards used across your instances, helping you maintain compliance and security in your database management.

```sql+postgres
select
  name,
  id,
  state,
  license_type,
  minimal_tls_version
from
  azure_mssql_managed_instance;
```

```sql+sqlite
select
  name,
  id,
  state,
  license_type,
  minimal_tls_version
from
  azure_mssql_managed_instance;
```

### List managed instances with public endpoint enabled
Identify instances where Azure's managed SQL servers have their public data endpoint enabled. This helps in assessing the elements within your setup that might be exposed to potential security risks.

```sql+postgres
select
  name,
  id,
  state,
  license_type,
  minimal_tls_version
from
  azure_mssql_managed_instance
where
  public_data_endpoint_enabled;
```

```sql+sqlite
select
  name,
  id,
  state,
  license_type,
  minimal_tls_version
from
  azure_mssql_managed_instance
where
  public_data_endpoint_enabled = 1;
```

### List security alert policies of the managed instances
Explore the security alert policies of managed instances to understand their configurations, such as creation time, disabled alerts, and retention days. This can help in assessing the security measures in place and identifying areas for improvement.

```sql+postgres
select
  name,
  id,
  policy -> 'creationTime' as policy_creation_time,
  jsonb_pretty(policy -> 'disabledAlerts') as policy_disabled_alerts,
  policy -> 'emailAccountAdmins' as policy_email_account_admins,
  jsonb_pretty(policy -> 'emailAddresses') as policy_email_addresses,
  policy ->> 'id' as policy_id,
  policy ->> 'name' as policy_name,
  policy -> 'retentionDays' as policy_retention_days,
  policy ->> 'state' as policy_state,
  policy ->> 'storageAccountAccessKey' as policy_storage_account_access_key,
  policy ->> 'storageEndpoint' as policy_storage_endpoint,
  policy ->> 'type' as policy_type
from
  azure_mssql_managed_instance,
  jsonb_array_elements(security_alert_policies) as policy;
```

```sql+sqlite
select
  name,
  id,
  json_extract(policy.value, '$.creationTime') as policy_creation_time,
  policy_disabled_alerts,
  json_extract(policy.value, '$.emailAccountAdmins') as policy_email_account_admins,
  policy_email_addresses,
  json_extract(policy.value, '$.id') as policy_id,
  json_extract(policy.value, '$.name') as policy_name,
  json_extract(policy.value, '$.retentionDays') as policy_retention_days,
  json_extract(policy.value, '$.state') as policy_state,
  json_extract(policy.value, '$.storageAccountAccessKey') as policy_storage_account_access_key,
  json_extract(policy.value, '$.storageEndpoint') as policy_storage_endpoint,
  json_extract(policy.value, '$.type') as policy_type
from
  azure_mssql_managed_instance,
  json_each(security_alert_policies) as policy;
```