# Table: azure_data_lake_analytics_account

Azure Data Lake Analytics is an on-demand analytics job service that simplifies big data. Instead of deploying, configuring, and tuning hardware, you write queries to transform your data and extract valuable insights. The analytics service can handle jobs of any scale instantly by setting the dial for how much power you need. You only pay for your job when it is running, making it cost-effective.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_analytics_account;
```

### List suspended data lake analytics accounts

```sql
select
  name,
  id,
  type,
  state,
  provisioning_state
from
  azure_data_lake_analytics_account
where
  state = 'Suspended';
```

### List data lake analytics accounts with firewall disabled

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_analytics_account
where
  firewall_state = 'Disabled';
```
