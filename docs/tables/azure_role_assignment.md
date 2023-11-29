---
title: "Steampipe Table: azure_role_assignment - Query Azure Identity and Access Management Role Assignments using SQL"
description: "Allows users to query Azure Role Assignments"
---

# Table: azure_role_assignment - Query Azure Identity and Access Management Role Assignments using SQL

Azure Role Assignments are a security mechanism used within Microsoft Azure to manage access to Azure resources. They define the operations that a user, group, service principal, or managed identity can perform on a particular resource. Role assignments are crucial for effective access management and security in Azure.

## Table Usage Guide

The 'azure_role_assignment' table provides insights into role assignments within Azure Identity and Access Management (IAM). As a security administrator, explore role assignment-specific details through this table, including role definitions, principals, and scope. Utilize it to uncover information about role assignments, such as those with broad permissions, the relationships between principals and roles, and the scope of each role assignment. The schema presents a range of attributes of the role assignment for your analysis, like the role ID, principal ID, scope, and role definition ID.

## Examples

### Basic info
Explore the identities and types of principals assigned to roles in your Azure environment, enabling you to better manage access and permissions. This is particularly useful in maintaining security and ensuring only authorized users have access to specific resources.

```sql
select
  name,
  id,
  principal_id,
  principal_type
from
  azure_role_assignment;
```

### List of role assignments which has permission at root level
Determine the areas in which certain role assignments have root-level permissions. This is useful for understanding the distribution of access rights within your Azure environment.

```sql
select
  name,
  id,
  scope
from
  azure_role_assignment
where
  scope = '/';
```

### List of role assignments which has subscription level permission and full access to the subscription
Explore which users have full access and subscription level permissions in Azure. This is beneficial for managing user permissions and ensuring the security of your Azure resources.

```sql
select
  ra.name as roll_assignment_name,
  rd.role_name
from
  azure_role_assignment ra
  join azure_role_definition rd on ra.role_definition_id = rd.id
  cross join jsonb_array_elements(rd.permissions) as perm
where
  ra.scope like '/subscriptions/%'
  and perm -> 'actions' = '["*"]';
```