---
title: "Steampipe Table: azure_ad_group - Query Azure Active Directory Groups using SQL"
description: "Allows users to query Azure Active Directory Groups, providing insights into group details, including identification, membership, and associated metadata."
---

# Table: azure_ad_group - Query Azure Active Directory Groups using SQL

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service. It helps your employees sign in and access resources in external resources, such as Microsoft Office 365, the Azure portal, and thousands of other SaaS applications. Azure AD Groups are collections of users and can be used to simplify the assignment of access rights to resources in Azure AD.

## Table Usage Guide

The `azure_ad_group` table provides insights into Azure Active Directory Groups within Microsoft Azure. As a system administrator, explore group-specific details through this table, including identification, membership, and associated metadata. Utilize it to manage access to resources, understand group composition, and maintain security compliance across your organization.

## Examples

### Basic info
Explore the groups within your Azure Active Directory to determine which ones have email capabilities enabled. This is useful for auditing purposes, ensuring that only necessary groups have email functions activated.

```sql+postgres
select
  display_name,
  object_id,
  mail,
  mail_enabled,
  mail_nickname
from
  azure_ad_group;
```

```sql+sqlite
select
  display_name,
  object_id,
  mail,
  mail_enabled,
  mail_nickname
from
  azure_ad_group;
```


### List of AD groups where security is not enabled
Determine the areas in which the security feature is not enabled in your Azure Active Directory groups. This can be useful for identifying potential vulnerabilities and taking corrective action to enhance your system's security.

```sql+postgres
select
  display_name,
  object_id,
  security_enabled
from
  azure_ad_group
where
  not security_enabled;
```

```sql+sqlite
select
  display_name,
  object_id,
  security_enabled
from
  azure_ad_group
where
  security_enabled = 0;
```


### List of AD groups where mail is not enabled
Explore which Azure Active Directory groups do not have mail enabled. This is useful to identify potential communication gaps within your organization.

```sql+postgres
select
  display_name,
  mail,
  mail_enabled
from
  azure_ad_group
where
  not mail_enabled;
```

```sql+sqlite
select
  display_name,
  mail,
  mail_enabled
from
  azure_ad_group
where
  mail_enabled is not 1;
```