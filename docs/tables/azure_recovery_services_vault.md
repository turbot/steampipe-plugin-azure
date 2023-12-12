---
title: "Steampipe Table: azure_recovery_services_vault - Query Azure Recovery Services Vaults using SQL"
description: "Allows users to query Azure Recovery Services Vaults, providing details about each vault's configuration, status, and associated resources."
---

# Table: azure_recovery_services_vault - Query Azure Recovery Services Vaults using SQL

Azure Recovery Services Vault is a backup service within Microsoft Azure that allows you to protect and recover your data and applications. It provides a centralized place to manage backups and disaster recovery, and it supports a range of Azure services, including virtual machines, SQL databases, and file shares. Azure Recovery Services Vault helps you ensure the availability and integrity of your Azure resources.

## Table Usage Guide

The `azure_recovery_services_vault` table provides insights into Recovery Services Vaults within Azure. As a system administrator, you can use this table to explore vault-specific details, including backup policies, protected items, and recovery points. This table is especially useful for ensuring that your backup and recovery strategies are properly implemented and managed.

## Examples

### Basic info
Explore the various elements of your Azure Recovery Services Vaults, such as their names, IDs, regions, and types. This can be useful in understanding the overall structure and organization of your vaults, aiding in better management and oversight.

```sql+postgres
select
  name,
  id,
  region,
  type
from
  azure_recovery_services_vault;
```

```sql+sqlite
select
  name,
  id,
  region,
  type
from
  azure_recovery_services_vault;
```

### List failed recovery service vaults
Discover the segments that have unsuccessful recovery service vaults in Azure, which can be crucial for identifying and addressing potential issues in your data recovery strategy. This query is beneficial in maintaining robust data protection and business continuity plans.

```sql+postgres
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

```sql+sqlite
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