# Table: azure_security_center_setting

Azure security center settings contains different configurations in security center.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  kind
from
  azure_security_center_setting;
```

### Ensure that Microsoft Cloud App Security (MCAS) integration with Security Center is selected

```sql
select
  id,
  name,
  type,
  kind
from
  azure_security_center_setting
where
  name = 'MCAS'
  and enabled;
```

### Ensure that Windows Defender ATP (WDATP) integration with Security Center is selected

```sql
select
  id,
  name,
  type,
  kind
from
  azure_security_center_setting
where
  name = 'WDATP'
  and enabled;
```
