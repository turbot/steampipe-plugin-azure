---
title: "Steampipe Table: azure_automation_account - Query Azure Automation Accounts using SQL"
description: "Allows users to query Azure Automation Accounts, providing insights into the configuration, status, and metadata of each account."
---

# Table: azure_automation_account - Query Azure Automation Accounts using SQL

Azure Automation is a service in Microsoft Azure that allows you to automate your Azure management tasks and to orchestrate actions across external systems from right within Azure. It enables you to automate frequent, time-consuming, and error-prone cloud management tasks. Azure Automation account is a container for holding your automation resources in Azure.

## Table Usage Guide

The `azure_automation_account` table provides insights into Automation Accounts within Microsoft Azure. As a cloud administrator, you can use this table to explore account-specific details, such as the configuration, status, and associated metadata. Leverage it to monitor and manage your automation resources, ensuring they are configured correctly and operating as expected.

## Examples

### Basic info
Explore which automation accounts are currently active within your Azure environment. This can be helpful for managing resources and understanding the types of automation accounts in use.

```sql+postgres
select
  name,
  id,
  resource_group,
  type
from
  azure_automation_account;
```

```sql+sqlite
select
  name,
  id,
  resource_group,
  type
from
  azure_automation_account;
```

### List accounts that are created in last 30 days
Explore which accounts were established in the past month. This can help in tracking recent activity and understanding the growth pattern of your accounts.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  resource_group,
  type,
  creation_time
from
  azure_automation_account
where
  creation_time >= datetime('now', '-30 day');
```

### List accounts that are suspended
Determine the areas in your Azure automation accounts where accounts are suspended. This query can be useful in identifying potential issues or disruptions in your automation workflows.

```sql+postgres
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

```sql+sqlite
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