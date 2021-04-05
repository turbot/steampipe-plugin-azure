# Table: azure_security_center

Azure Security Center provides unified security management and advanced threat protection across hybrid cloud workloads. With Security Center, you can apply security policies across your workloads, limit your exposure to threats, and detect and respond to attacks.

## Examples

### Ensure that Microsoft Cloud App Security (MCAS) integration with Security Center is selected

```sql
select
  jsonb_pretty(setting) as setting
from
  azure_security_center
where
  jsonb_path_exists(
    setting,
    '$.** ? (@.type() == "string" && @ like_regex "MCAS")'
  );
```

### Ensure that Windows Defender ATP (WDATP) integration with Security Center is selected

```sql
select
  jsonb_pretty(setting) as setting
from
  azure_security_center
where
  jsonb_path_exists(
    setting,
    '$.** ? (@.type() == "string" && @ like_regex "WDATP")'
  );
```

### Ensure that Automatic provisioning of monitoring agent is set to On

```sql
select
  p -> 'properties' ->> 'autoProvision' as autoProvision
from
  azure_security_center,
  jsonb_array_elements(auto_provisioning) as p
where
  auto_provisioning is not null
  and p -> 'properties' ->> 'autoProvision' = 'On';
```

### Ensure 'Additional email addresses' is configured with a security contact email

```sql
select
  c -> 'properties' ->> 'email' as contact_email
from
  azure_security_center,
  jsonb_array_elements(contact) as c
where
  contact is not null and c -> 'properties' ->> 'email' != '';
```
