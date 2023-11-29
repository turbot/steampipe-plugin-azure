---
title: "Steampipe Table: azure_data_lake_analytics_account - Query Azure Data Lake Analytics Accounts using SQL"
description: "Allows users to query Azure Data Lake Analytics Accounts"
---

# Table: azure_data_lake_analytics_account - Query Azure Data Lake Analytics Accounts using SQL

Azure Data Lake Analytics is an on-demand analytics job service that simplifies big data. It allows you to focus on writing, running and managing jobs, rather than operating distributed infrastructure. Instead of deploying, configuring and tuning hardware, you write queries to transform your data and extract valuable insights.

## Table Usage Guide

The 'azure_data_lake_analytics_account' table provides insights into Azure Data Lake Analytics Accounts. As a data analyst or a big data engineer, explore account-specific details through this table, including account status, creation date, last modified date, and associated metadata. Utilize it to uncover information about accounts, such as those with specific firewall states, the maximum degree of parallelism per job, and the maximum job count. The schema presents a range of attributes of the Azure Data Lake Analytics Account for your analysis, like the account ID, name, type, and associated tags.

## Examples

### Basic info
Explore which Azure Data Lake Analytics accounts are currently provisioned to gain insights into your active data processing resources. This can help you manage your resources efficiently and plan for future capacity needs.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_analytics_account;
```

### List suspended data lake analytics accounts
Identify instances where data lake analytics accounts are suspended to enable proactive management and prevent potential disruptions in data processing.

```sql
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
Identify instances where Azure Data Lake analytics accounts have their firewall disabled. This query is useful for assessing potential security vulnerabilities in your system.

```sql
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