---
title: "Steampipe Table: azure_security_center_contact - Query Azure Security Center Contacts using SQL"
description: "Allows users to query Azure Security Center Contacts, providing insights into contact details, alert notifications, and alert email settings."
folder: "Security Center"
---

# Table: azure_security_center_contact - Query Azure Security Center Contacts using SQL

Azure Security Center Contacts is a feature within Microsoft Azure that allows you to manage and configure the security contact details in Azure Security Center. These contact details are used by Azure to send notifications regarding security alerts, recommendations, and other important security information. It is a crucial component of Azure's security management system, providing a streamlined way to receive and manage security notifications.

## Table Usage Guide

The `azure_security_center_contact` table provides insights into the contact details configured in Azure Security Center. As a security administrator, explore contact-specific details through this table, including alert notifications, and alert email settings. Utilize it to manage and monitor the communication of security alerts and recommendations from Azure to the designated contacts.

## Examples

### Basic info
Explore which security center contacts in your Azure environment have alert notifications enabled. This helps to identify who is receiving alerts and whether any necessary contacts are missing from the notifications list.

```sql+postgres
select
  id,
  email,
  notifications_by_role,
  notifications_sources
from
  azure_security_center_contact;
```

```sql+sqlite
select
  id,
  email,
  notifications_by_role,
  notifications_sources
from
  azure_security_center_contact;
```

### List security center contacts not configured with email notifications
Determine areas in which Security Center contacts have been set up without email notifications. This is useful to ensure that all necessary parties are receiving important security alerts and updates.

```sql+postgres
select
  id,
  email,
  notifications_by_role,
  notifications_sources
from
  azure_security_center_contact
where
  email != '';
```

```sql+sqlite
select
  id,
  email,
  notifications_by_role,
  notifications_sources
from
  azure_security_center_contact
where
  email != '';
```