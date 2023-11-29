---
title: "Steampipe Table: azure_recovery_services_vault - Query Azure Recovery Services Vaults using SQL"
description: "Allows users to query Azure Recovery Services Vaults"
---

# Table: azure_recovery_services_vault - Query Azure Recovery Services Vaults using SQL

Azure Recovery Services vault is a management entity that stores recovery points created over time and provides an interface to perform backup related operations. These operations include taking on-demand backups, performing restores, and creating backup policies. It offers backup support for Azure virtual machines, SQL workloads, and on-premises VMware machines.

## Table Usage Guide

The 'azure_recovery_services_vault' table provides insights into Recovery Services Vaults within Azure Recovery Services. As a DevOps engineer, explore vault-specific details through this table, such as the vault's location, resource group, subscription ID, and associated tags. Utilize it to uncover information about each vault, including its storage redundancy and soft delete feature status. The schema presents a range of attributes of the Recovery Services Vault for your analysis, like the vault name, type, SKU name, and provisioning state.

## Examples

### Basic info
Explore the different types of recovery services vaults available in various regions of your Azure environment. This can help in managing and organizing your backup and disaster recovery resources effectively.

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
Discover the segments that have unsuccessful recovery service vaults in Azure. This is useful to pinpoint specific locations where the provisioning process failed, allowing for targeted troubleshooting and resolution.

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