# Table: azure_mssql_virtual_machine

Azure SQL virtual machines are lift-and-shift ready for existing applications that require fast migration to the cloud with minimal changes or no changes. SQL virtual machines offer full administrative control over the SQL Server instance and underlying OS for migration to Azure.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  provisioning_state,
  sql_image_offer,
  sql_server_license_type,
  region
from
  azure_mssql_virtual_machine;
```

### List failed virtual machines

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_mssql_virtual_machine
where
  provisioning_state = 'Failed';
```
