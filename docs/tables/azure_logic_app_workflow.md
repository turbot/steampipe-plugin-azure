---
title: "Steampipe Table: azure_logic_app_workflow - Query Azure Logic Apps Workflows using SQL"
description: "Allows users to query Azure Logic Apps Workflows."
---

# Table: azure_logic_app_workflow - Query Azure Logic Apps Workflows using SQL

Azure Logic Apps is a cloud service that helps you schedule, automate, and orchestrate tasks, business processes, and workflows when you need to integrate apps, data, systems, and services across enterprises or organizations. It provides a way to simplify and implement scalable integrations and workflows in the cloud. Logic Apps allows you to develop and deliver powerful integration solutions with ease.

## Table Usage Guide

The 'azure_logic_app_workflow' table provides insights into the workflows within Azure Logic Apps. As an engineer, you can explore workflow-specific details through this table, including workflow status, integration account, endpoints, and associated metadata. Utilize it to uncover information about workflows, such as those with enabled or disabled state, the integration account associated with the workflow, and the endpoints used by the workflow. The schema presents a range of attributes of the workflow for your analysis, like the workflow ID, creation date, state, and associated tags.

## Examples

### Basic info
Explore which Azure Logic App Workflows are currently active by identifying their names and types. This can help in assessing the elements within your Azure environment and managing your resources effectively.

```sql
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow;
```

### List disabled workflows
Identify instances where specific workflows in Azure Logic App have been disabled. This enables users to manage and rectify any potential disruptions in their workflow processes.

```sql
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
Uncover the details of workflows that are currently suspended within your Azure Logic App, allowing you to identify and address any potential issues or disruptions in your workflow processes.

```sql
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