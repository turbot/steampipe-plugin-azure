# Table: azure_eventgrid_domain

An event domain is a management tool for large numbers of Event Grid topics related to the same application. You can think of it as a meta-topic that can have thousands of individual topics.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_domain;
```

### List domains not configured with private endpoint connections

```sql
select
  name,
  id,
  type,
  private_endpoint_connections
from
  azure_eventgrid_domain
where
  private_endpoint_connections is null;
```

### List domains with local authentication disabled

```sql
select
  name,
  id,
  type,
  disable_local_auth
from
  azure_eventgrid_domain
where
  disable_local_auth;
```
