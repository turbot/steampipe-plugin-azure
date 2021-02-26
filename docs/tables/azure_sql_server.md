# Table: azure_sql_server

An Azure SQL server is a relational database management system. As a database server, it is a software product with the primary function of storing and retrieving data as requested by other software applications—which may run either on the same computer or on another computer across a network (including the Internet).

## Examples

### List of server for which Auditing is not enabled

```sql
select
  name,
  id,
  audit -> 'properties' ->> 'state' as audit_policy_state
from
  azure_sql_server,
  jsonb(server_audit_policy) as audit
where
  audit -> 'properties' ->> 'state' = 'Disabled';
```

### List of server for which Auditing retention is greater than 90 days

```sql
select
  name,
  id,
  audit -> 'properties' -> 'retentionDays' as audit_policy_retention_days
from
  azure_sql_server,
  jsonb(server_audit_policy) as audit
where
  audit -> 'properties' ->> 'retentionDays' >= '90';
```

### List of server for which Advanced Data Security is not enabled

```sql
select
  name,
  id,
  security -> 'properties' ->> 'state' as security_alert_policy_state
from
  azure_sql_server,
  jsonb(server_security_alert_policy) as security
where
  security -> 'properties' ->> 'state' = 'Disabled';
```

### List of server for which Threat Detection types is set to All

```sql
select
  name,
  id,
  security -> 'properties' -> 'disabledAlerts' as security_alert_policy_state
from
  azure_sql_server,
  jsonb(server_security_alert_policy) as security,
  jsonb_array_elements_text(security -> 'properties' -> 'disabledAlerts') as disabled_alerts,
  jsonb_array_length(security -> 'properties' -> 'disabledAlerts') as alert_length
where
  alert_length = 1
  and disabled_alerts = '';
```

### List of server for which Active Directory Admin is not configured

```sql
select
  name,
  id
from
  azure_sql_server
where
  server_azure_ad_administrator is null;
```

### List of server for which TDE protector is encrypted with Service-managed key

```sql
select
  name,
  id,
  encryption ->> 'kind' as encryption_protector_kind
from
  azure_sql_server,
  jsonb(encryption_protector) as encryption
where
  encryption ->> 'kind' = 'servicemanaged';
```
