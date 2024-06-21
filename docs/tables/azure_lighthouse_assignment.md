---
title: "Steampipe Table: azure_lighthouse_assignment - Query Azure Lighthouse Assignments using SQL"
description: "Allows users to query Azure Lighthouse assignments, providing insights into the management and governance of resources across multiple tenants."
---

# Table: azure_lighthouse_assignment - Query Azure Lighthouse Assignments using SQL

Azure Lighthouse is a service within Microsoft Azure that enables cross-tenant management, allowing service providers to manage resources across multiple tenants while maintaining control and visibility. Azure Lighthouse assignments are specific configurations that apply these management capabilities to designated resources.

## Table Usage Guide

The `azure_lighthouse_assignment` table provides insights into Azure Lighthouse assignments. As a Network Administrator or Service Provider, you can explore details about each assignment, including its configuration, provisioning state, and associated registration definition. Use this table to ensure your cross-tenant management and governance assignments are correctly configured and to quickly identify any potential issues.

## Examples

### Basic info
Explore the status and details of your Azure Lighthouse assignments to understand their current state and type. This is beneficial for auditing and managing your cross-tenant resources effectively.

```sql+postgres
select
  name,
  id,
  provisioning_state,
  registration_assignment_id,
  registration_definition_id,
  type
from
  azure_lighthouse_assignment;
```

```sql+sqlite
select
  name,
  id,
  provisioning_state,
  registration_assignment_id,
  registration_definition_id,
  type
from
  azure_lighthouse_assignment;
```

### List assignments by resource group
Identify the Azure Lighthouse assignments based on their resource group. This can help in organizing and managing assignments within specific resource groups.

```sql+postgres
select
  name,
  id,
  resource_group,
  scope,
  type
from
  azure_lighthouse_assignment
where
  resource_group = 'your_resource_group_name';
```

```sql+sqlite
select
  name,
  id,
  resource_group,
  scope,
  type
from
  azure_lighthouse_assignment
where
  resource_group = 'your_resource_group_name';
```

### List assignments with specific provisioning state
Explore the Azure Lighthouse assignments that have a specific provisioning state. This helps in monitoring the status and ensuring that assignments are correctly provisioned.

```sql+postgres
select
  name,
  id,
  provisioning_state,
  scope,
  type
from
  azure_lighthouse_assignment
where
  provisioning_state = 'Succeeded';
```

```sql+sqlite
select
  name,
  id,
  provisioning_state,
  scope,
  type
from
  azure_lighthouse_assignment
where
  provisioning_state = 'Succeeded';
```

### List assignments by scope
Get an overview of Azure Lighthouse assignments based on their scope. This can assist in understanding the management scope and ensuring it aligns with your governance requirements.

```sql+postgres
select
  name,
  id,
  scope,
  registration_definition_id,
  type
from
  azure_lighthouse_assignment
where
  scope = '/subscriptions/your_subscription_id';
```

```sql+sqlite
select
  name,
  id,
  scope,
  registration_definition_id,
  type
from
  azure_lighthouse_assignment
where
  scope = '/subscriptions/your_subscription_id';
```

### Determine the scope for assignments
This query is highly useful for normalizing, analyzing, and reporting on Azure resource scopes in an environment managed by Azure Lighthouse.

```sql+postgres
select
  case
    when id like '/subscriptions/%/resourceGroups/%/providers/%/%/%' then
      substring(id from '/subscriptions/[^/]+/resourceGroups/[^/]+/providers/[^/]+/[^/]+/[^/]+')
    when id like '/subscriptions/%/resourceGroups/%' then
      substring(id from '/subscriptions/[^/]+/resourceGroups/[^/]+')
    when id like '/subscriptions/%' then
      substring(id from '/subscriptions/[^/]+')
    when id like '/providers/Microsoft.Management/managementGroups/%' then
      substring(id from '/providers/Microsoft.Management/managementGroups/[^/]+')
    else
      null
  end as scope_id,
  registration_definition_id,
  id
from
  azure_lighthouse_assignment;
```

```sql+sqlite
select
  case
    when id like '/subscriptions/%/resourceGroups/%/providers/%/%/%' then
      substr(id, 1, instr(id, '/', 3, 5) + length('/providers') - 1)
    when id like '/subscriptions/%/resourceGroups/%' then
      substr(id, 1, instr(id, '/', 3, 4) + length('/resourceGroups') - 1)
    when id like '/subscriptions/%' then
      substr(id, 1, instr(id, '/', 3, 2) - 1)
    when id like '/providers/Microsoft.Management/managementGroups/%' then
      substr(id, 1, instr(id, '/', 3, 4) + length('/managementGroups') - 1)
    else
      null
  end as scope_id,
  registration_definition_id,
  id
from
  azure_lighthouse_assignment;
```
