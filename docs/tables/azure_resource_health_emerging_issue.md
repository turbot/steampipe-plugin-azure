# Table: azure_resource_health_emerging_issue

An "emerging issue" in the context of Azure Resource Health refers to a situation where Azure identifies a potential problem or degradation of service that might impact your resources.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  resource_group
from
  azure_resource_health_emerging_issue;
```

### Get status banner details of emerging issues

```sql
select
  name,
  id,
  s ->> 'Title' as status_banner_title,
  s ->> 'Message' as status_banner_message,
  s ->> 'Cloud' as status_banner_cloud,
  s ->> 'LastModifiedTime' as status_banner_last_modified_time
from
  azure_resource_health_emerging_issue,
  jsonb_array_elements(status_banners) as s;
```

### Get status event details of emerging issues

```sql
select
  name,
  e ->> 'Title' as event_title,
  e ->> 'Description' as event_description,
  e ->> 'TrackingID' as tracking_id,
  e ->> 'StartTime' as start_time,
  e ->> 'Cloud' as cloud,
  e ->> 'Severity' as severity,
  e ->> 'Stage' as stage,
  e ->> 'Published' as published,
  e ->> 'LastModifiedTime' as last_modified_time,
  e ->> 'Impacts' as impacts
from
  azure_resource_health_emerging_issue,
  jsonb_array_elements(status_active_events) as e;
```

### Get emerging issue details of virtual machines

```sql
select
  name,
  id,
  type,
  status_banners,
from
  azure_resource_health_emerging_issue
where
  type = 'Microsoft.Compute/virtualMachines';
```