---
title: "Steampipe Table: azure_batch_account - Query Azure Batch Accounts using SQL"
description: "Allows users to query Azure Batch Accounts, specifically retrieving details such as the account name, resource group, location, and subscription ID."
folder: "Batch"
---

# Table: azure_batch_account - Query Azure Batch Accounts using SQL

Azure Batch is a cloud-based job scheduling service that parallelizes and distributes the processing of large volumes of data across many computers. It is designed to manage and run hundreds to thousands of tasks concurrently. This service reduces the time and cost associated with processing large amounts of data.

## Table Usage Guide

The `azure_batch_account` table provides insights into Batch Accounts within Azure Batch service. As a data engineer, explore account-specific details through this table, including account name, resource group, location, and subscription ID. Utilize it to uncover information about accounts, such as those with specific locations, the resource groups they belong to, and the verification of subscription IDs.

## Examples

### Basic info
Explore which Azure Batch accounts are currently provisioned, along with their dedicated core quotas and regional locations. This can be particularly useful for managing resources and optimizing cloud infrastructure.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  dedicated_core_quota,
  region
from
  azure_batch_account;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  dedicated_core_quota,
  region
from
  azure_batch_account;
```

### List failed batch accounts
Identify instances where Azure batch account provisioning has failed. This is useful for troubleshooting and understanding the areas where resource allocation has been unsuccessful.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  dedicated_core_quota,
  region
from
  azure_batch_account
where
  provisioning_state = 'Failed';
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  dedicated_core_quota,
  region
from
  azure_batch_account
where
  provisioning_state = 'Failed';
```