# Table: azure_security_center_contact

Azure security center contact configurations for the subscription.

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

### Ensure security contact email configured for the subscription

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
