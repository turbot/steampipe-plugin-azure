---
title: "Steampipe Table: azure_ad_service_principal - Query Azure Active Directory Service Principals using SQL"
description: "Allows users to query Azure Active Directory Service Principals, specifically the details about the service principals in the Azure Active Directory."
---

# Table: azure_ad_service_principal - Query Azure Active Directory Service Principals using SQL

An Azure Active Directory Service Principal is a security identity used by user-created applications, services, and automation tools to access specific Azure resources. It allows these resources to be secured by using Azure AD role-based access control. This identity is used to authenticate to Azure AD and obtain tokens to access resources.

## Table Usage Guide

The `azure_ad_service_principal` table provides insights into Service Principals within Azure Active Directory. As a Security Engineer, utilize this table to explore details about service principals, including their app roles, display names, and associated metadata. Use it to uncover information about service principals, such as those with specific permissions, their associated application IDs, and the verification of OAuth2 permissions.

## Examples

### List of ad service principals where service principal account is disabled
Determine the areas in which Azure ad service principals are disabled. This can be useful for identifying potential security risks or troubleshooting access issues.

```sql+postgres
select
  object_id,
  object_type,
  display_name,
  account_enabled
from
  azure_ad_service_principal
where
  not account_enabled;
```

```sql+sqlite
select
  object_id,
  object_type,
  display_name,
  account_enabled
from
  azure_ad_service_principal
where
  account_enabled = 0;
```

### List of ad service principals where app role assignment is not required
Identify instances where ad service principals in Azure do not require an app role assignment. This can be useful to streamline access control and reduce unnecessary role assignments.

```sql+postgres
select
  object_id,
  display_name,
  app_role_assignment_required
from
  azure_ad_service_principal
where
  not app_role_assignment_required;
```

```sql+sqlite
select
  object_id,
  display_name,
  app_role_assignment_required
from
  azure_ad_service_principal
where
  app_role_assignment_required = 0;
```

### Application role info of service principals
Explore the roles assigned to service principals within your Azure Active Directory. This query helps in understanding the permissions and access controls for each service principal, thereby assisting in maintaining secure and efficient system operations.

```sql+postgres
select
  object_id,
  approle ->> 'allowedMemberTypes' as allowed_member_types,
  approle ->> 'description' as description,
  approle ->> 'displayName' as display_name,
  approle -> 'isEnabled' as isEnabled,
  approle ->> 'id' as id,
  approle ->> 'value' as id
from
  azure_ad_service_principal
  cross join jsonb_array_elements(app_roles) as approle;
```

```sql+sqlite
select
  object_id,
  json_extract(approle.value, '$.allowedMemberTypes') as allowed_member_types,
  json_extract(approle.value, '$.description') as description,
  json_extract(approle.value, '$.displayName') as display_name,
  json_extract(approle.value, '$.isEnabled') as isEnabled,
  json_extract(approle.value, '$.id') as id,
  json_extract(approle.value, '$.value') as id
from
  azure_ad_service_principal,
  json_each(app_roles) as approle;
```

### Oauth 2.0 permission info of ad service principal
This query is useful for gaining insights into the permissions associated with your Azure advertising service principal. It allows you to assess whether certain permissions are enabled and understand their specific descriptions and display names, helping to maintain proper access control in your Azure environment.

```sql+postgres
select
  object_id,
  perm ->> 'adminConsentDescription' as admin_consent_description,
  perm ->> 'adminConsentDisplayName' as admin_consent_display_ame,
  perm ->> 'id' as id,
  perm ->> 'isEnabled' as is_enabled,
  perm ->> 'type' as type,
  perm ->> 'value' as value
from
  azure_ad_service_principal
  cross join jsonb_array_elements(oauth2_permissions) as perm;
```

```sql+sqlite
select
  object_id,
  json_extract(perm.value, '$.adminConsentDescription') as admin_consent_description,
  json_extract(perm.value, '$.adminConsentDisplayName') as admin_consent_display_name,
  json_extract(perm.value, '$.id') as id,
  json_extract(perm.value, '$.isEnabled') as is_enabled,
  json_extract(perm.value, '$.type') as type,
  json_extract(perm.value, '$.value') as value
from
  azure_ad_service_principal,
  json_each(oauth2_permissions) as perm;
```