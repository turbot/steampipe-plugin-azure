# Table: azure_recovery_services_vault

A Recovery Services vault is a storage entity in Azure that houses data. The data is typically copies of data, or configuration information for virtual machines (VMs), workloads, servers, or workstations. You can use Recovery Services vaults to hold backup data for various Azure services such as IaaS VMs (Linux or Windows) and Azure SQL databases.

## Examples

### Basic info

```sql
select
  name,
  id,
  region,
  type
from
  azure_recovery_services_vault;
```

### List failed recovery service vaults

```sql
select
  name,
  id,
  type,
  provisioning_state,
  region
from
  azure_recovery_services_vault
where
provisioning_state = 'Failed';
```
