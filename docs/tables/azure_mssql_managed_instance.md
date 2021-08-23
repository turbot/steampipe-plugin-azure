# Table: azure_mssql_managed_instance

Azure SQL Managed Instance is the intelligent, scalable cloud database service that combines the broadest SQL Server database engine compatibility with all the benefits of a fully managed and evergreen platform as a service.

## Examples

### Basic info

```sql
select
  name,
  id,
  state,
  license_type,
  minimal_tls_version
from
  azure_mssql_managed_instance;
```

### List managed instances with public endpoint enabled

```sql
select
  name,
  id,
  state,
  license_type,
  minimal_tls_version
from
  azure_mssql_managed_instance
where
  public_data_endpoint_enabled;
```
