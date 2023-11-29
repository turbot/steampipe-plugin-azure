---
title: "Steampipe Table: azure_security_center_automation - Query Azure Security Center Automations using SQL"
description: "Allows users to query Azure Security Center Automations, providing detailed information on their security automation configurations."
---

# Table: azure_security_center_automation - Query Azure Security Center Automations using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers and provides advanced threat protection across your hybrid workloads in the cloud. The Security Center Automations are part of this system, designed to provide automatic responses to specific security incidents. They offer an efficient way to remediate threats and misconfigurations, enabling a proactive approach to security management.

## Table Usage Guide

The 'azure_security_center_automation' table provides insights into the automations within Azure Security Center. As a security or DevOps engineer, explore automation-specific details through this table, including the associated resources, actions, and conditions. Utilize it to uncover information about automations, such as those related to specific security alerts, the actions taken in response, and the resources affected. The schema presents a range of attributes of the automation for your analysis, like the automation name, resource group, subscription ID, and associated tags.

## Examples

### Basic info 
Explore the types and kinds of security automations set up in your Azure Security Center. This is useful for understanding the variety and scope of automated security measures currently in place.

```sql
select
  id,
  name,
  type,
  kind
from
   azure_security_center_automation;
```

### List enabled continuously export microsoft defender for cloud data
Analyze the configuration of your Microsoft Defender for cloud data to identify which aspects are continuously exporting. This helps in keeping track of the data and ensuring that all necessary information is being exported as required.

```sql
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
Determine the details of event sources for continuous data export in Microsoft Defender for Cloud. This is useful for understanding the configuration and operators of your security automation rules, as well as identifying expected values and property types.

```sql
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