# Table: azure_security_center_automation

Microsoft Defender for Cloud generates detailed security alerts and recommendations. You can view them in the portal or through programmatic tools. You might also need to export some or all of this information for tracking with other monitoring tools in your environment.

## Examples

### Basic info 

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

### List event source details for continoourly export microsoft defender for cloud data

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