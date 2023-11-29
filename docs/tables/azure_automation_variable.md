---
title: "Steampipe Table: azure_automation_variable - Query Azure Automation Variables using SQL"
description: "Allows users to query Azure Automation Variables"
---

# Table: azure_automation_variable - Query Azure Automation Variables using SQL

Azure Automation is a service that allows you to automate your Azure management tasks and to orchestrate actions across external systems from right within Azure. Variables in Azure Automation are used to store values that can be accessed across runbooks and modules during their execution. These variables can store different types of values, such as strings, integers, Booleans, and DateTime values.

## Table Usage Guide

The 'azure_automation_variable' table provides insights into Variables within Azure Automation. As a DevOps engineer, explore variable-specific details through this table, including names, types, values, and descriptions. Utilize it to uncover information about variables, such as those with specific values, the encrypted status of the variables, and the last time they were updated. The schema presents a range of attributes of the Automation Variable for your analysis, like the variable id, creation time, last modified time, and associated tags.

## Examples

### Basic info
Explore the basic information of Azure Automation Variables to understand the type and encryption status. This can help in managing and securing the automation environment.

```sql
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
Discover the segments that contain unencrypted variables within your Azure Automation account. This is useful for identifying potential security risks and ensuring that all sensitive information is adequately protected.

```sql
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
Discover the segments that have been newly added in the past month, which can be useful in understanding recent changes or additions to your system. This can help in assessing the elements within your system that have been recently modified or created.

```sql
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

### Get details of a variable
Explore the specific settings of a variable within a given account and resource group in Azure Automation. This is useful for assessing the elements within your automation environment, such as identifying if a variable is encrypted or not.

```sql
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