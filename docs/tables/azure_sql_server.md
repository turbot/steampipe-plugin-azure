# Table: azure_sql_server

An Azure SQL server is a relational database management system. As a database server, it is a software product with the primary function of storing and retrieving data as requested by other software applications—which may run either on the same computer or on another computer across a network (including the Internet).

## Examples

### List servers that have auditing disabled

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

### List servers with an audit log retention period that is less than 90 days

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
