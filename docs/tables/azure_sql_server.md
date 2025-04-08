---
title: "Steampipe Table: azure_sql_server - Query Azure SQL Servers using SQL"
description: "Allows users to query Azure SQL Servers, providing detailed information about their configurations, locations, versions, and more."
folder: "SQL"
---

# Table: azure_sql_server - Query Azure SQL Servers using SQL

Azure SQL Server is a fully managed relational database service, which is a part of the broader Microsoft Azure Platform. It offers the broadest SQL Server engine compatibility and powers your cloud applications with AI-built-in, secure and manageable data platform. The service provides automatic updates, scaling, provisioning, backups, and monitoring, leaving developers free to focus on application design and optimization.

## Table Usage Guide

The `azure_sql_server` table provides insights into SQL Servers within Microsoft Azure. As a database administrator or developer, explore server-specific details through this table, including server versions, locations, configurations, and more. Utilize it to uncover information about servers, such as their current state, the number of databases, the firewall rules, and the performance tier.

## Examples

### List servers that have auditing disabled
Identify instances where auditing is disabled on Azure SQL servers. This is beneficial for enhancing security measures by pinpointing potential weaknesses in your server configurations.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  json_extract(audit.value, '$.properties.state') as audit_policy_state
from
  azure_sql_server,
  json_each(server_audit_policy) as audit
where
  json_extract(audit.value, '$.properties.state') = 'Disabled';
```

### List servers with an audit log retention period less than 90 days
Assess the elements within your system to identify servers that have an audit log retention period of less than 90 days. This is useful to ensure compliance with data retention policies and to identify potential risks associated with short retention periods.

```sql+postgres
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

```sql+sqlite
select
  name,
  s.id,
  json_extract(audit.value, '$.properties.retentionDays') as audit_policy_retention_days
from
  azure_sql_server as s,
  json_each(server_audit_policy) as audit
where
  json_extract(audit.value, '$.properties.retentionDays') < 90;
```

### List servers that have advanced data security disabled
This query helps identify servers where advanced data security is turned off. This is useful for quickly pinpointing potential security risks in your server infrastructure.

```sql+postgres
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

```sql+sqlite
select
  name,
  s.id,
  json_extract(security.value, '$.properties.state') as security_alert_policy_state
from
  azure_sql_server as s,
  json_each(server_security_alert_policy) as security
where
  json_extract(security.value, '$.properties.state') = 'Disabled';
```

### List servers that have Advanced Threat Protection types set to All
Determine the areas in which Azure SQL servers have their Advanced Threat Protection set to 'All'. This can help to assess the security measures in place and identify any potential vulnerabilities.

```sql+postgres
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

```sql+sqlite
select
  name,
  s.id,
  json_extract(security.value, '$.properties.disabledAlerts') as security_alert_policy_state
from
  azure_sql_server as s,
  json_each(server_security_alert_policy) as security,
  json_each(json_extract(security.value, '$.properties.disabledAlerts')) as disabled_alerts
where
  json_array_length(json_extract(security.value, '$.properties.disabledAlerts')) = 1
  and disabled_alerts.value = '';
```

### List servers that do not have an Active Directory admin set
Identify Azure SQL servers that are potentially vulnerable due to the absence of an Active Directory administrator. This can help in enhancing security by ensuring all servers have designated administrators.

```sql+postgres
select
  name,
  id
from
  azure_sql_server
where
  server_azure_ad_administrator is null;
```

```sql+sqlite
select
  name,
  id
from
  azure_sql_server
where
  server_azure_ad_administrator is null;
```

### List servers for which TDE protector is encrypted with the service-managed key
Determine the servers where the Transparent Data Encryption (TDE) protector is encrypted using a service-managed key. This is useful for understanding your server's encryption setup and ensuring it aligns with your organization's security policies.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  json_extract(encryption.value, '$.kind') as encryption_protector_kind
from
  azure_sql_server,
  json_each(encryption_protector) as encryption
where
  json_extract(encryption.value, '$.kind') = 'servicemanaged';
```