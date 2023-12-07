---
title: "Steampipe Table: azure_ad_user - Query Azure Active Directory Users using SQL"
description: "Allows users to query Azure Active Directory Users, providing details of user profiles including user details, email addresses, and department information."
---

# Table: azure_ad_user - Query Azure Active Directory Users using SQL

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service. It helps your employees sign in and access resources in external resources, such as Microsoft Office 365, the Azure portal, and thousands of other SaaS applications. Azure AD also includes a full suite of identity management capabilities including multi-factor authentication, device registration, role-based access control, user provisioning, and more.

## Table Usage Guide

The `azure_ad_user` table provides insights into user profiles within Azure Active Directory. As a system administrator, explore user-specific details through this table, including user details, email addresses, and department information. Utilize it to uncover information about users, such as their roles, access controls, and associated metadata.

## Examples

### Basic active directory user info
Determine the areas in which active directory users are currently active within the Azure environment. This query is beneficial in managing user access and maintaining security standards.

```sql+postgres
select
  display_name,
  user_principal_name,
  given_name,
  mail,
  account_enabled,
  object_id
from
  azure_ad_user;
```

```sql+sqlite
select
  display_name,
  user_principal_name,
  given_name,
  mail,
  account_enabled,
  object_id
from
  azure_ad_user;
```


### List of guest users in the active directory
Identify instances where guest users are present in the active directory to maintain security and access control. This query is useful in managing permissions and keeping track of external users in your system.

```sql+postgres
select
  display_name,
  user_principal_name,
  mail,
  user_type,
  usage_location
from
  azure_ad_user
where
  user_type = 'Guest';
```

```sql+sqlite
select
  display_name,
  user_principal_name,
  mail,
  user_type,
  usage_location
from
  azure_ad_user
where
  user_type = 'Guest';
```


### Password profile info of each user
This example helps in understanding the password policies applied to each user within the Azure Active Directory. It aids in determining whether users are required to change their passwords at their next login or if the password change policy is enforced, thereby assisting in maintaining security standards.

```sql+postgres
select
  display_name,
  user_principal_name,
  additional_properties -> 'passwordProfile' -> 'enforceChangePasswordPolicy' as enforce_change_password_policy,
  additional_properties -> 'passwordProfile' -> 'forceChangePasswordNextLogin' as change_password_next_login
from
  azure_ad_user;
```

```sql+sqlite
select
  display_name,
  user_principal_name,
  json_extract(additional_properties, '$.passwordProfile.enforceChangePasswordPolicy') as enforce_change_password_policy,
  json_extract(additional_properties, '$.passwordProfile.forceChangePasswordNextLogin') as change_password_next_login
from
  azure_ad_user;
```