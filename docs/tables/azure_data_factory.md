# Table: azure_data_factory

Azure Data Factory is the platform that solves such data scenarios. It is the cloud-based ETL and data integration service that allows to create data-driven workflows for orchestrating data movement and transforming data at scale.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state,
  etag
from
  azure_data_factory;
```


### List system assigned identity type factories

```sql
select
  name,
  id,
  type,
  identity ->> 'type' as identity_type
from
  azure_data_factory
where
  identity ->> 'type' = 'SystemAssigned';
```


### List factories with public network access enabled

```sql
select
  name,
  id,
  type,
  public_network_access
from
  azure_data_factory
where
  public_network_access = 'Enabled';
```
