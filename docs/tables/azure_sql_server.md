---
title: "Steampipe Table: azure_sql_server - Query Azure SQL Servers using SQL"
description: "Allows users to query Azure SQL Servers for information such as server names, locations, versions, administrator logins, and more."
---

# Table: azure_sql_server - Query Azure SQL Servers using SQL

Azure SQL Server is a relational database service in the Microsoft Azure cloud. It provides a scalable, highly available, and managed database service that handles most of the database management functions such as upgrading, patching, backups, and monitoring. Azure SQL Server offers the broadest SQL Server engine compatibility and an automated patching and version updates feature.

## Table Usage Guide

The 'azure_sql_server' table provides insights into SQL servers within Azure SQL Server service. As a database administrator, you can explore server-specific details through this table, including server names, locations, versions, administrator logins, and more. Utilize it to uncover information about servers, such as those with specific versions, the locations of the servers, and the administrator login details. The schema presents a range of attributes of the SQL server for your analysis, like the server name, location, version, administrator login, and associated tags.

## Examples

### List servers that have auditing disabled
Determine the areas in which auditing is disabled on your servers. This can be useful to maintain security standards and ensure all activities are properly recorded for future reference.

```sql
select
  name,
  id,
  audit -> 'properties' ->> 'state' as audit_policy_state
from
  azure_sql_server,
  jsonb_array_elements(server_audit_policy) as audit
where
  audit -> 'properties' ->> 'state' = 'Disabled';
```

### List servers with an audit log retention period less than 90 days
Determine the servers that have an audit log retention period of less than 90 days. This can be useful for identifying potential security risks and ensuring compliance with internal or external data retention policies.

```sql
select
  name,
  id,
  (audit -> 'properties' ->> 'retentionDays')::integer as audit_policy_retention_days
from
  azure_sql_server,
  jsonb_array_elements(server_audit_policy) as audit
where
  (audit -> 'properties' ->> 'retentionDays')::integer < 90;
```

### List servers that have advanced data security disabled
Discover the segments that have advanced data security turned off in your Azure SQL servers. This is particularly useful for assessing potential vulnerabilities and ensuring optimal security practices.

```sql
select
  name,
  id,
  security -> 'properties' ->> 'state' as security_alert_policy_state
from
  azure_sql_server,
  jsonb_array_elements(server_security_alert_policy) as security
where
  security -> 'properties' ->> 'state' = 'Disabled';
```

### List servers that have Advanced Threat Protection types set to All
Explore which Azure servers have their Advanced Threat Protection types set to 'All'. This is useful for assessing the security configuration of servers and identifying any potential vulnerabilities.

```sql
select
  name,
  id,
  security -> 'properties' -> 'disabledAlerts' as security_alert_policy_state
from
  azure_sql_server,
  jsonb_array_elements(server_security_alert_policy) as security,
  jsonb_array_elements_text(security -> 'properties' -> 'disabledAlerts') as disabled_alerts,
  jsonb_array_length(security -> 'properties' -> 'disabledAlerts') as alert_length
where
  alert_length = 1
  and disabled_alerts = '';
```

### List servers that do not have an Active Directory admin set
Analyze your Azure SQL servers to identify those that lack an assigned Active Directory admin. This could be beneficial in pinpointing potential security vulnerabilities or compliance issues in your infrastructure.

```sql
select
  name,
  id
from
  azure_sql_server
where
  server_azure_ad_administrator is null;
```

### List servers for which TDE protector is encrypted with the service-managed key
Explore which servers have their Transparent Data Encryption (TDE) protector encrypted with a service-managed key. This is useful for assessing the encryption status and understanding the key management scheme of your servers.

```sql
select
  name,
  id,
  encryption ->> 'kind' as encryption_protector_kind
from
  azure_sql_server,
  jsonb_array_elements(encryption_protector) as encryption
where
  encryption ->> 'kind' = 'servicemanaged';
```