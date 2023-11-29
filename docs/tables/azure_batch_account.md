---
title: "Steampipe Table: azure_batch_account - Query Azure Batch Accounts using SQL"
description: "Allows users to query Azure Batch Accounts."
---

# Table: azure_batch_account - Query Azure Batch Accounts using SQL

Azure Batch is a cloud-based job scheduling service that parallelizes and distributes the processing of large volumes of data across many computers. It is designed for high-performance computing (HPC) applications, enabling developers and scientists to run large-scale parallel and high-performance computing (HPC) applications efficiently in the cloud. Azure Batch creates and manages a pool of compute nodes (virtual machines), installs the applications you want to run, and schedules jobs to run on the nodes.

## Table Usage Guide

The 'azure_batch_account' table provides insights into Batch Accounts within Azure Batch. As a DevOps engineer, explore account-specific details through this table, including the provisioning state, pool allocation mode, and associated metadata. Utilize it to uncover information about accounts, such as those with public network access, the key vault reference, and the verification of pool allocation mode. The schema presents a range of attributes of the Batch Account for your analysis, like the account name, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which Azure Batch accounts are active and their dedicated core quota limits, to manage resource allocation and prevent potential overuse. This helps in maintaining cost-effective and efficient operations within your Azure environment.

```sql
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
Explore which batch accounts in your Azure environment have failed to provision. This is useful for identifying and addressing potential issues in resource allocation or configuration.

```sql
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