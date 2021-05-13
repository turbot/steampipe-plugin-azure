# Table: azure_security_center_contact

Azure security center contact details for the subscription.

## Examples

### Basic info

```sql
select
  id,
  email,
  alert_notifications,
  alerts_to_admins
from
  azure_security_center_contact;
```

### List security center contacts not configured with email notifications

```sql
select
  id,
  email,
  alert_notifications,
  alerts_to_admins
from
  azure_security_center_contact
where
  email != '';
```
