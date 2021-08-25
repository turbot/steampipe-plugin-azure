# Table: azure_mssql_elasticpool

Azure SQL Database elastic pools are a simple, cost-effective solution for managing and scaling multiple databases that have varying and unpredictable usage demands.

## Examples

### Basic info

```sql
select
  name,
  id,
  state,
  dtu,
  zone_redundant
from
  azure_mssql_elasticpool;
```

### List zone redundant elastic pools

```sql
select
  name,
  id,
  state,
  dtu,
  zone_redundant
from
  azure_mssql_elasticpool
where
  zone_redundant;
```
