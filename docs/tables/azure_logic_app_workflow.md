---
title: "Steampipe Table: azure_logic_app_workflow - Query Azure Logic App Workflows using SQL"
description: "Allows users to query Azure Logic App Workflows, specifically details regarding the configuration, status, and properties of each workflow, aiding in the management and monitoring of automated business processes."
folder: "Monitor"
---

# Table: azure_logic_app_workflow - Query Azure Logic App Workflows using SQL

Azure Logic Apps is a cloud service that helps you schedule, automate, and orchestrate tasks, business processes, and workflows when you need to integrate apps, data, systems, and services across enterprises or organizations. Logic Apps simplifies how you design and build scalable solutions for app integration, data integration, system integration, enterprise application integration (EAI), and business-to-business (B2B) communication, whether in the cloud, on premises, or both. You can build workflows that automatically trigger and run whenever a specific event occurs or when new data meets the conditions that you defined.

## Table Usage Guide

The `azure_logic_app_workflow` table provides insights into Logic App Workflows within Azure. As a DevOps engineer, explore workflow-specific details through this table, including configuration, status, and associated properties. Utilize it to monitor and manage automated business processes, ensuring efficient data and system integration.

## Examples

### Basic info
Explore which Azure Logic App workflows are active or inactive and their associated types. This can be beneficial for identifying workflows that may need attention or tracking the variety of workflow types in use.

```sql+postgres
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow;
```

```sql+sqlite
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow;
```

### List disabled workflows
This query allows you to identify all disabled workflows within your Azure Logic App, helping to manage and optimize your workflow processes. This could be particularly useful for troubleshooting, auditing, or improving efficiency by identifying unused or unnecessary workflows.

```sql+postgres
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow
where
  state = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow
where
  state = 'Disabled';
```

### List suspended workflows
Determine the areas in which Azure Logic App workflows are currently suspended. This is useful for troubleshooting and identifying workflows that may need attention or modification.

```sql+postgres
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow
where
  state = 'Suspended';
```

```sql+sqlite
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow
where
  state = 'Suspended';
```