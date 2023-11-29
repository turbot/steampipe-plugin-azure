---
title: "Steampipe Table: azure_automation_account - Query Azure Automation Accounts using SQL"
description: "Allows users to query Azure Automation Accounts."
---

# Table: azure_automation_account - Query Azure Automation Accounts using SQL

Azure Automation is a service in Microsoft Azure that allows users to automate their manual, long-running, error-prone, and frequently repeated tasks. It provides process automation, update management and configuration features, and integrates with other popular DevOps tools. Azure Automation helps users to focus on work that adds business value by reducing the time spent on routine tasks.

## Table Usage Guide

The 'azure_automation_account' table provides insights into Automation Accounts within Azure Automation. As a DevOps engineer, explore account-specific details through this table, including the account's name, ID, location, and type. Utilize it to uncover information about accounts, such as their provisioning state, creation time, last modified time, and their SKU. The schema presents a range of attributes of the Automation Account for your analysis, like the subscription ID, tenant ID, resource group, and associated tags.

## Examples

### Basic info
Explore the different automation accounts within your Azure environment, including their names and associated resource groups. This can help you manage and organize your resources more effectively.

```sql
select
  name,
  id,
  resource_group,
  type
from
  azure_automation_account;
```

### List accounts that are created in last 30 days
Gain insights into newly created accounts in the past month. This query is useful for tracking recent account activity and managing resources within Azure automation.

```sql
select
  name,
  id,
  resource_group,
  type,
  creation_time
from
  azure_automation_account
where
  creation_time >= now() - interval '30' day;
```

### List accounts that are suspended
Explore which Azure Automation accounts are currently suspended. This can be useful in identifying any potential issues or disruptions in your automation tasks and processes.

```sql
select
  name,
  id,
  resource_group,
  type,
  creation_time,
  state
from
  azure_automation_account
where
  state = 'AccountStateSuspended';
```