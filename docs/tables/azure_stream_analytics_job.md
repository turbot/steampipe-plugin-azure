---
title: "Steampipe Table: azure_stream_analytics_job - Query Azure Stream Analytics Jobs using SQL"
description: "Allows users to query Azure Stream Analytics Jobs."
---

# Table: azure_stream_analytics_job - Query Azure Stream Analytics Jobs using SQL

Azure Stream Analytics is a real-time analytics and complex event-processing engine that is designed to analyze and visualize streaming data in real-time. It provides users with the ability to set up real-time analytic computations on streaming data which can originate from various sources such as devices, sensors, websites, social media feeds, applications, infrastructure systems, and more. Azure Stream Analytics is designed to process and analyze data as it's ingested in real-time, and can handle high volumes of data from multiple sources simultaneously.

## Table Usage Guide

The 'azure_stream_analytics_job' table provides insights into Stream Analytics Jobs within Azure Stream Analytics. As a Data Engineer, explore job-specific details through this table, including job topology, transformation query, output details, and associated metadata. Utilize it to uncover information about jobs, such as those with their current state, the events processed, and the input and output of the job. The schema presents a range of attributes of the Stream Analytics Job for your analysis, like the job name, resource group, location, compatibility level, data locale, and job type.

## Examples

### Basic info
Explore which Azure Stream Analytics jobs are currently active or inactive, and identify their respective locations and subscription IDs. This information can be useful for auditing purposes or for managing and optimizing streaming jobs across different regions.

```sql
select
  name,
  id,
  job_id,
  job_state,
  region,
  subscription_id
from
  azure_stream_analytics_job;
```

### List failed stream analytics jobs
Identify instances where stream analytics jobs have failed in Azure. This can be useful for troubleshooting and understanding the areas that may require additional resources or configuration adjustments.

```sql
select
  name,
  id,
  type,
  provisioning_state,
  region
from
  azure_stream_analytics_job
where
  provisioning_state = 'Failed';
```