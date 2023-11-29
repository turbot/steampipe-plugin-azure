---
title: "Steampipe Table: azure_key_vault_managed_hardware_security_module - Query Azure Key Vault Managed Hardware Security Modules using SQL"
description: "Allows users to query Azure Key Vault Managed Hardware Security Modules."
---

# Table: azure_key_vault_managed_hardware_security_module - Query Azure Key Vault Managed Hardware Security Modules using SQL

A Managed Hardware Security Module (HSM) is a service offered by Azure Key Vault that provides cryptographic key storage in Azure. It provides secure, FIPS 140-2 Level 3 validated, cryptographic key storage and operations using Azure Key Vault. It is designed to meet the stringent requirements of highly regulated industries that process, store, and use sensitive data.

## Table Usage Guide

The 'azure_key_vault_managed_hardware_security_module' table provides insights into Managed Hardware Security Modules within Azure Key Vault. As a security or DevOps engineer, explore module-specific details through this table, including its cryptographic keys, key operations, and associated metadata. Utilize it to uncover information about modules, such as their key identifiers, enabled status, and creation time. The schema presents a range of attributes of the Managed Hardware Security Module for your analysis, like the resource ID, name, type, and location.

## Examples

### Basic info
Analyze the settings to understand the configuration of your Azure Key Vault Managed Hardware Security Module. This query can help you assess the elements within your system, such as its name, ID, type, and region, as well as whether the soft delete option is enabled.

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
Identify instances where Azure Key Vault Managed Hardware Security Modules do not have the soft delete feature enabled. This is useful for ensuring data protection and recovery in case of accidental deletion.

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