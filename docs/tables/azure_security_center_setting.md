---
title: "Steampipe Table: azure_security_center_setting - Query Azure Security Center Settings using SQL"
description: "Allows users to query Azure Security Center Settings, specifically the configuration data, providing insights into security settings and potential discrepancies."
---

# Table: azure_security_center_setting - Query Azure Security Center Settings using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of data centers and provides advanced threat protection across hybrid workloads in the cloud. It provides a centralized way to monitor and respond to security issues across your Azure resources, including virtual machines, databases, web applications, and more. Azure Security Center helps you stay informed about the security status and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `azure_security_center_setting` table provides insights into settings within Azure Security Center. As a Security Analyst, explore setting-specific details through this table, including configurations, contact details, and auto provisioning settings. Utilize it to uncover information about settings, such as those with auto provisioning enabled, the contact details for security notifications, and the verification of security configurations.

## Examples

### Basic info
Discover the segments that are enabled in the Azure Security Center. This query is useful for quickly assessing the active areas of your security configuration.

```sql+postgres
select
  id,
  name,
  enabled
from
  azure_security_center_setting;
```

```sql+sqlite
select
  id,
  name,
  enabled
from
  azure_security_center_setting;
```

### List the enabled settings for security center
Explore which security settings are currently activated in the Azure Security Center to ensure your system is adequately protected and compliant with security protocols. This is useful for maintaining a secure environment and identifying any potential gaps in your security configuration.

```sql+postgres
select
  id,
  name,
  type
from
  azure_security_center_setting
where
  enabled;
```

```sql+sqlite
select
  id,
  name,
  type
from
  azure_security_center_setting
where
  enabled = 1;
```