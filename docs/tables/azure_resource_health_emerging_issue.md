---
title: "Steampipe Table: azure_resource_health_emerging_issue - Query Azure Service Health Emerging Issues using SQL"
description: "Allows users to query Azure Service Health Emerging Issues, providing detailed information about the emerging issues."
---

# Table: azure_resource_health_emerging_issue - Query Azure Service Health Emerging Issues using SQL

Azure Resource Health's "Emerging Issues" feature in Azure Service Health is designed to inform users about new service issues that might impact them. This tool surfaces issues where Azure acknowledges a widespread outage but may not have full clarity on its extent. Previously available only on the Azure Status page, these emerging issues are now integrated into the Service Health dashboard. This integration allows users to rely solely on Service Health for a comprehensive view of both personalized service health information and broader emerging issues, enhancing the overall awareness and management of Azure service health.

## Table Usage Guide

The `azure_resource_health_emerging_issue` table provides insights into Resource Health within Microsoft Azure. As a DevOps engineer, explore group-specific details through this table, including properties, name of the resource, resource type, refresh time, status banner, and status active events.

## Examples

### Basic info
This is useful for quickly identifying the issue at hand, the type of the emerging issue, which helps in categorizing and understanding the nature of the issue, and the last refresh timestamp for the issue’s data. This is important for understanding the recency of the information, ensuring that you are working with the latest data.

```sql+postgres
select
  name,
  id,
  type,
  refresh_timestamp
from
  azure_resource_health_emerging_issue;
```

```sql+sqlite
select
  name,
  id,
  type,
  refresh_timestamp
from
  azure_resource_health_emerging_issue;
```

### Get status banner details of emerging issues
Extracting detailed information about emerging issues in Azure Resource Health.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  json_extract(s.value, '$.Title') as status_banner_title,
  json_extract(s.value, '$.Message') as status_banner_message,
  json_extract(s.value, '$.Cloud') as status_banner_cloud,
  json_extract(s.value, '$.LastModifiedTime') as status_banner_last_modified_time
from
  azure_resource_health_emerging_issue,
  json_each(status_banners) as s;
```

### Get status event details of emerging issues
This query is useful for monitoring and analysis purposes, enabling users to keep track of ongoing issues and their characteristics within the Azure environment.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(e.value, '$.Title') as event_title,
  json_extract(e.value, '$.Description') as event_description,
  json_extract(e.value, '$.TrackingID') as tracking_id,
  json_extract(e.value, '$.StartTime') as start_time,
  json_extract(e.value, '$.Cloud') as cloud,
  json_extract(e.value, '$.Severity') as severity,
  json_extract(e.value, '$.Stage') as stage,
  json_extract(e.value, '$.Published') as published,
  json_extract(e.value, '$.LastModifiedTime') as last_modified_time,
  json_extract(e.value, '$.Impacts') as impacts
from
  azure_resource_health_emerging_issue
  json_each(status_active_events) as e;
```