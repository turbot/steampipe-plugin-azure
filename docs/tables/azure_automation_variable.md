---
title: "Steampipe Table: azure_automation_variable - Query Azure Automation Variables using SQL"
description: "Allows users to query Azure Automation Variables, providing a comprehensive view of all variables that are used within Azure Automation."
folder: "Automation"
---

# Table: azure_automation_variable - Query Azure Automation Variables using SQL

Azure Automation is a cloud-based service offered by Microsoft Azure that allows organizations to automate and configure certain tasks across their Azure and non-Azure environments. It provides a way to manage, monitor, and act upon the infrastructure resources in a scalable and reliable manner. Azure Automation Variables are entities within this service that store values which can be used in runbooks, DSC configurations, and other features of Azure Automation.

## Table Usage Guide

The `azure_automation_variable` table provides insights into Azure Automation Variables within Azure Automation. As a DevOps engineer, explore variable-specific details through this table, including the type of value it holds, whether it is encrypted, and its associated metadata. Utilize it to uncover information about variables, such as their current values, and the automation accounts they are associated with.

## Examples

### Basic info
Analyze the settings to understand the encryption status of variables in your Azure automation accounts. This will help in assessing the security posture of your automation workflows.

```sql+postgres
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable;
```

```sql+sqlite
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable;
```

### List variables that are unencrypted
Discover the segments that are unencrypted within your Azure Automation variables. This can be beneficial in identifying potential security risks and ensuring data protection standards are met.

```sql+postgres
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  not is_encrypted;
```

```sql+sqlite
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  not is_encrypted;
```

### List variables created in last 30 days
Explore the recent changes in your Azure Automation environment by identifying variables that were created within the last 30 days. This can help you monitor and control your automation tasks, ensuring they align with your current needs and standards.

```sql+postgres
select
  id,
  name,
  account_name,
  creation_time,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  creation_time >= now() - interval '30' day;
```

```sql+sqlite
select
  id,
  name,
  account_name,
  creation_time,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  creation_time >= datetime('now', '-30 day');
```

### Get details of a variable
Discover the specifics of a particular variable within a given account and resource group. This is useful in understanding the variable's attributes, such as its type and encryption status, which can aid in managing and securing your Azure automation tasks.

```sql+postgres
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  account_name = 'turbot_account'
  and name = 'turbot'
  and resource_group = 'turbot_rg';
```

```sql+sqlite
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  account_name = 'turbot_account'
  and name = 'turbot'
  and resource_group = 'turbot_rg';
```