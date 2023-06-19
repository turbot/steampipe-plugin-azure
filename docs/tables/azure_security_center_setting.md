# Table: azure_security_center_setting

Azure security center settings contains different configurations in security center.

## Examples

### Basic info

```sql
select
  id,
  name,
  enabled
from
  azure_security_center_setting;
```

### List the enabled settings for security center

```sql
select
  id,
  name,
  type
from
  azure_security_center_setting
where
  enabled;
```
