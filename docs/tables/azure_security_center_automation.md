---
title: "Steampipe Table: azure_security_center_automation - Query Azure Security Center Automations using SQL"
description: "Allows users to query Azure Security Center Automations, specifically the automation details and configurations, providing insights into security automation settings and potential vulnerabilities."
---

# Table: azure_security_center_automation - Query Azure Security Center Automations using SQL

Azure Security Center Automations is a feature within Microsoft Azure that allows you to automate responses to security alerts. It provides an automated way to respond to and manage alerts for various Azure resources, including virtual machines, databases, web applications, and more. Azure Security Center Automations helps you stay informed about the security state of your Azure resources and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `azure_security_center_automation` table provides insights into the automation settings within Azure Security Center. As a security engineer, explore automation-specific details through this table, including the automation name, resource group, and associated metadata. Utilize it to uncover information about your security automations, such as their configurations, intended actions, and the resources they are associated with.

## Examples

### Basic info
This example showcases how to determine the types and categories of automations within the Azure Security Center. This information can be useful in understanding the range of automated processes in place and their respective functions.

```sql+postgres
select
  id,
  name,
  type,
  kind
from
   azure_security_center_automation;
```

```sql+sqlite
select
  id,
  name,
  type,
  kind
from
   azure_security_center_automation;
```

### List enabled continuously export microsoft defender for cloud data
Determine the areas in which Microsoft Defender for Cloud data is continuously exported and enabled. This can be useful to ensure that your security data is being properly and consistently exported for further analysis and storage.

```sql+postgres
select
  id,
  name,
  type,
  is_enabled
from
  azure_security_center_automation
where
  is_enabled;
```

```sql+sqlite
select
  id,
  name,
  type,
  is_enabled
from
  azure_security_center_automation
where
  is_enabled;
```

### List event source details for continuously export microsoft defender for cloud data
Determine the areas in which continuous data export from Microsoft Defender for Cloud is occurring. This is useful for understanding your security posture and identifying potential areas of improvement.

```sql+postgres
select
  name,
  type,
  s ->> 'eventSource' as event_source,
  r ->> 'operator' as operator,
  r ->> 'propertyType' as property_type,
  r ->> 'expectedValue' as expected_value,
  r ->> 'propertyJPath' as property_jpath
from
  azure_security_center_automation,
  jsonb_array_elements(sources) as s,
  jsonb_array_elements(s -> 'ruleSets') as rs,
  jsonb_array_elements(rs -> 'rules') as r ;
```

```sql+sqlite
select
  name,
  a.type,
  json_extract(s.value, '$.eventSource') as event_source,
  json_extract(r.value, '$.operator') as operator,
  json_extract(r.value, '$.propertyType') as property_type,
  json_extract(r.value, '$.expectedValue') as expected_value,
  json_extract(r.value, '$.propertyJPath') as property_jpath
from
  azure_security_center_automation as a,
  json_each(sources) as s,
  json_each(json_extract(s.value, '$.ruleSets')) as rs,
  json_each(json_extract(rs.value, '$.rules')) as r ;
```