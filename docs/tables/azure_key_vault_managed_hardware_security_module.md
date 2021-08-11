# Table: azure_key_vault_managed_hardware_security_module

Azure Key Vault Managed HSM is a fully managed, highly available, single-tenant, standards-compliant cloud service that enables you to safeguard cryptographic keys for your cloud applications, using FIPS 140-2 Level 3 validated HSMs.

## Examples

### Basic info

```sql
select
  name,
  id,
  hsm_uri,
  type,
  enable_soft_delete,
  region
from
  azure_key_vault_managed_hardware_security_module;
```

### List soft delete disabled hsm managed key vaults

```sql
select
  name,
  id,
  enable_soft_delete
from
  azure_key_vault_managed_hardware_security_module
where
  not enable_soft_delete;
```