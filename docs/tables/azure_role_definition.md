---
title: "Steampipe Table: azure_role_definition - Query Azure Active Directory Role Definitions using SQL"
description: "Allows users to query Azure Active Directory Role Definitions."
---

# Table: azure_role_definition - Query Azure Active Directory Role Definitions using SQL

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service. Role Definitions are a collection of permissions. It’s a template that defines the operations that can be performed, such as read, write, and delete.

## Table Usage Guide

The 'azure_role_definition' table provides insights into role definitions within Azure Active Directory. As a DevOps engineer, you can explore role-specific details through this table, including permissions and associated metadata. Utilize it to uncover information about role definitions, such as those with specific permissions. The schema presents a range of attributes of the role definition for your analysis, like the role name, id, description, and type.

## Examples

### List the custom roles
Explore the custom roles within your Azure environment to understand their configurations and purposes. This can help in managing access and permissions more effectively.

```sql
select
  name,
  description,
  role_name,
  role_type,
  title
from
  azure_role_definition
where
  role_type = 'CustomRole';
```

### List of roles whose assignable scope is set to root('/') scope
Discover the roles within the Azure environment that have the highest level of access, specifically those set to the root ('/') scope. This can be useful for auditing purposes, allowing you to ensure only the appropriate roles have such broad permissions.

```sql
select
  name,
  role_name,
  scope
from
  azure_role_definition,
  jsonb_array_elements_text(assignable_scopes) as scope
where
  scope = '/';
```

### Permissions of all custom roles
Explore the permissions associated with all custom roles in an Azure environment. This can be useful to understand and manage access control, ensuring the right roles have the appropriate permissions.

```sql
select
  name,
  role_name,
  role_type,
  permission -> 'actions' as action,
  permission -> 'dataActions' as data_action,
  permission -> 'notActions' as no_action,
  permission -> 'notDataActions' as not_data_actions
from
  azure_role_definition
  cross join jsonb_array_elements(permissions) as permission
where
  role_type = 'CustomRole';
```