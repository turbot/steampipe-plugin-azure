---
title: "Steampipe Table: azure_key_vault_managed_hardware_security_module - Query Azure Key Vaults using SQL"
description: "Allows users to query Azure Key Vaults, specifically those managed by hardware security modules, providing insights into key management, encryption, and decryption services."
folder: "Key Vault"
---

# Table: azure_key_vault_managed_hardware_security_module - Query Azure Key Vaults using SQL

Azure Key Vault is a service within Microsoft Azure that provides secure key management and cryptographic protection services. It offers solutions for securely storing and accessing secrets, keys, and certificates, while also providing logging for all key usage. A managed hardware security module (HSM) in Azure Key Vault provides cryptographic key storage in FIPS 140-2 Level 3 validated HSMs.

## Table Usage Guide

The `azure_key_vault_managed_hardware_security_module` table provides insights into Azure Key Vaults managed by hardware security modules. As a security engineer, explore vault-specific details through this table, including keys, secrets, and certificates, and their associated metadata. Utilize it to uncover information about key usage, key permissions, and the verification of cryptographic protection services.

## Examples

### Basic info
Explore the configuration of Azure's Key Vault Managed Hardware Security Module to understand its current settings and location. This is useful for auditing security measures and ensuring data is stored in the correct geographical region.

```sql+postgres
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

```sql+sqlite
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
Identify instances where the soft delete feature is disabled in Azure Key Vault Managed Hardware Security Modules. This is useful for enhancing data security by ensuring that deleted data can be recovered.

```sql+postgres
select
  name,
  id,
  enable_soft_delete
from
  azure_key_vault_managed_hardware_security_module
where
  not enable_soft_delete;
```

```sql+sqlite
select
  name,
  id,
  enable_soft_delete
from
  azure_key_vault_managed_hardware_security_module
where
  enable_soft_delete = 0;
```