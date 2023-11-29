---
title: "Steampipe Table: azure_security_center_setting - Query Azure Security Center Settings using SQL"
description: "Allows users to query Azure Security Center Settings"
---

# Table: azure_security_center_setting - Query Azure Security Center Settings using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers, and provides advanced threat protection across your hybrid workloads in the cloud - whether they're in Azure or not. It gives you the ability to protect your hybrid cloud workloads and get unified security management across your entire environment. Azure Security Center helps you prevent, detect, and respond to threats with increased visibility and control over the security of all your Azure resources.

## Table Usage Guide

The 'azure_security_center_setting' table provides insights into settings within Azure Security Center. As a Security Engineer, explore setting-specific details through this table, including the type of setting, whether it is enabled or not, and the kind of resource it is associated with. Utilize it to uncover information about settings, such as those that are disabled, those that are enabled, and the resources they are associated with. The schema presents a range of attributes of the Security Center setting for your analysis, like the setting name, type, kind, provisioning state, and associated metadata.

## Examples

### Basic info
Explore the status of your Azure Security Center settings to determine which ones are active. This can help streamline your security management by focusing on the settings currently in use.

```sql
select
  id,
  name,
  enabled
from
  azure_security_center_setting;
```

### List the enabled settings for security center
Explore which settings are enabled in the Azure Security Center to determine the areas of your system that are currently secured. This can help in identifying any potential vulnerabilities or gaps in security.

```sql
select
  id,
  name,
  type
from
  azure_security_center_setting
where
  enabled;
```