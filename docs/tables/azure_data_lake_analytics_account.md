---
title: "Steampipe Table: azure_data_lake_analytics_account - Query Azure Data Lake Analytics Accounts using SQL"
description: "Allows users to query Azure Data Lake Analytics Accounts, providing insights into the configuration, state, and other critical details of these resources."
folder: "Data Lake Analytics"
---

# Table: azure_data_lake_analytics_account - Query Azure Data Lake Analytics Accounts using SQL

Azure Data Lake Analytics is an on-demand analytics job service that simplifies big data. Instead of deploying, configuring, and tuning hardware, you write queries to transform your data and extract valuable insights. The analytics service can handle jobs of any scale instantly by setting the dial for how much power you need.

## Table Usage Guide

The `azure_data_lake_analytics_account` table provides insights into the Azure Data Lake Analytics Accounts within Azure. As a Data Analyst or a Data Engineer, explore account-specific details through this table, including the current state, the level of commitment, and the associated metadata. Utilize it to uncover information about accounts, such as their provisioning state, firewall state, and the maximum degree of parallelism per job, which can aid in optimizing data processing and analytics tasks.

## Examples

### Basic info
Analyze the settings to understand the status and type of your Azure Data Lake Analytics accounts. This can help in managing resources and identifying any accounts that may be in an unexpected state.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_analytics_account;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_analytics_account;
```

### List suspended data lake analytics accounts
Determine the areas in which your data lake analytics accounts are suspended. This allows you to manage resources efficiently by identifying and addressing any issues with these accounts.

```sql+postgres
select
  name,
  id,
  type,
  state,
  provisioning_state
from
  azure_data_lake_analytics_account
where
  state = 'Suspended';
```

```sql+sqlite
select
  name,
  id,
  type,
  state,
  provisioning_state
from
  azure_data_lake_analytics_account
where
  state = 'Suspended';
```

### List data lake analytics accounts with firewall disabled
Explore which Data Lake analytics accounts have their firewall disabled in order to identify potential security vulnerabilities. This can assist in maintaining robust security practices and preventing unauthorized access.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_analytics_account
where
  firewall_state = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_analytics_account
where
  firewall_state = 'Disabled';
```