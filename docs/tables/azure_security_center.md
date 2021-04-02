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


### Check the status of the Automatic provisioning of monitoring agent

```sql
select
  jsonb_pretty(auto_provisioning) as auto_provisioning
from
  azure_security_center;
```


### Ensure 'Additional email addresses' is configured with a security contact email

```sql
select
  jsonb_pretty(contact) as contact
from
  azure_security_center;
```