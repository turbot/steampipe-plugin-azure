# Table: azure_application_insight

Application Insights is an extension of Azure Monitor and provides Application Performance Monitoring (also known as “APM”) features.

## Examples

### Basic info

```sql
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight;
```

### List application insights having retention period less than 30 days

```sql
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  retention_in_days < 30;
```

### List insights that can be queried publicly

```sql
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  public_network_access_for_query = 'Enabled';
```

### List insights that allow ingestion publicly

```sql
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  public_network_access_for_ingestion = 'Enabled';
```