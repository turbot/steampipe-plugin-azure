---
title: "Steampipe Table: azure_security_center_contact - Query Azure Security Center Contacts using SQL"
description: "Allows users to query Azure Security Center Contacts."
---

# Table: azure_security_center_contact - Query Azure Security Center Contacts using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers and provides advanced threat protection across your hybrid workloads in the cloud. It allows you to manage and enforce your security policies across your Azure environment, limit your exposure to threats, and detect and respond to attacks. A contact in Azure Security Center is an entity that contains the contact details for the security center.

## Table Usage Guide

The 'azure_security_center_contact' table provides insights into the contacts within Azure Security Center. As a security administrator, explore contact-specific details through this table, including email addresses, phone numbers, and alert notifications. Utilize it to uncover information about contacts, such as those who are set to receive security alerts, and the verification of alert notifications. The schema presents a range of attributes of the Security Center contact for your analysis, like the contact name, email, phone, and alert notifications.

## Examples

### Basic info
Analyze the settings to understand the alert preferences and email contact details in your Azure Security Center. This can help you ensure that alerts are being sent to the right people and that the notification settings are configured correctly.

```sql
select
  id,
  email,
  alert_notifications,
  alerts_to_admins
from
  azure_security_center_contact;
```

### List security center contacts not configured with email notifications
Discover the segments that have security center contacts without configured email notifications. This is useful to identify potential gaps in your alert system and ensure all relevant parties are receiving necessary security updates.

```sql
select
  id,
  email,
  alert_notifications,
  alerts_to_admins
from
  azure_security_center_contact
where
  email != '';
```