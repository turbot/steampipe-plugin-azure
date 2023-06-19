# Table: azure_sql_database

An Azure SQL Database is an intelligent, scalable, relational database service built for the cloud.

## Examples

### Basic info

```sql
select
  name,
  id,
  server_name,
  location,
  edition
from
  azure_sql_database;
```

### List databases that are not online

```sql
select
  name,
  id,
  server_name,
  location,
  edition,
  status
from
  azure_sql_database
where
  status != 'Online';
```

### List databases that are not encrypted

```sql
select
  name,
  id,
  server_name,
  location,
  edition,
  transparent_data_encryption ->> 'status' as encryption_status
from
  azure_sql_database
where
  transparent_data_encryption ->> 'status' != 'Enabled';
```
