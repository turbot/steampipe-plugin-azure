# Table: azure_sql_database

An Azure SQL Database is an intelligent, scalable, relational database service built for the cloud. Optimise performance and durability with automated, AI-powered features that are always up to date.

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


### List databases for a specific server

```sql
select
  name,
  id,
  server_name,
  location,
  edition
from
  azure_sql_database
where
  server_name = 'test-steampipe';
```


### List databases having system edition

```sql
select
  name,
  id,
  server_name,
  location,
  edition
from
  azure_sql_database
where
  edition = 'Basic';
```


### List databases which are online

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
  status = 'Online';
```


### List databases which are encrypted

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
  transparent_data_encryption ->> 'status' = 'Enabled';
```


### List databases created before last 7 days

```sql
select
  name,
  id,
  server_name,
  location,
  creation_date
from
  azure_sql_database
where
  creation_date < (current_date - 7);
```