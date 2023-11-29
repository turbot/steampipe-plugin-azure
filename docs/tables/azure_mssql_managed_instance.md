---
title: "Steampipe Table: azure_mssql_managed_instance - Query Azure Managed SQL Server Instances using SQL"
description: "Allows users to query Azure Managed SQL Server Instances."
---

# Table: azure_mssql_managed_instance - Query Azure Managed SQL Server Instances using SQL

An Azure Managed SQL Server Instance is a fully managed relational database service provided by Microsoft Azure. It offers the broadest SQL Server engine compatibility and automates most of the database management functions such as upgrading, patching, backups, and monitoring. It also provides built-in intelligence that learns app patterns and adapts to maximize performance, reliability, and data protection.

## Table Usage Guide

The 'azure_mssql_managed_instance' table provides insights into Managed SQL Server Instances within Microsoft Azure. As a Database Administrator, explore instance-specific details through this table, including the instance's administrative settings, network settings, and associated metadata. Utilize it to uncover information about instances, such as their current state, the number of vCores, the maximum storage size, and the license type. The schema presents a range of attributes of the Managed SQL Server Instance for your analysis, like the instance's ID, name, type, location, and SKU.

## Examples

### Basic info
Explore the status and security settings of your managed instances in Azure's SQL service. This can be useful in assessing compliance with your organization's security policies.

```sql
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
Discover the segments that have the public data endpoint enabled in your managed instances. This can help identify potential security vulnerabilities, as these instances can be accessed publicly.

```sql
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

### List security alert policies of the managed instances
Explore the security alert policies of managed instances to understand their creation time, the alerts that have been disabled, and the email addresses linked to the policies. This can help in assessing the current security measures and making necessary improvements for better data protection.

```sql
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