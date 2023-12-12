---
title: "Steampipe Table: azure_stream_analytics_job - Query Azure Stream Analytics Jobs using SQL"
description: "Allows users to query Azure Stream Analytics Jobs, providing insights into their configurations, statuses, and other associated details."
---

# Table: azure_stream_analytics_job - Query Azure Stream Analytics Jobs using SQL

Azure Stream Analytics is a real-time analytics service that allows you to analyze and visualize streaming data from various sources such as devices, sensors, websites, social media feeds, and applications. It enables you to set up real-time analytic computations on streaming data which can be used for anomaly detection, live dashboarding, and alerts among other scenarios. This service is designed to process and analyze high volumes of fast streaming data from multiple streams simultaneously.

## Table Usage Guide

The `azure_stream_analytics_job` table provides insights into Stream Analytics Jobs within Azure. As a data analyst or data scientist, you can explore job-specific details through this table, including job configurations, input and output details, and transformation queries. Utilize it to monitor the status and health of your Stream Analytics Jobs, understand their configurations, and ensure they are processing data as expected.

## Examples

### Basic info
Explore which stream analytics jobs are currently running in your Azure environment. This allows you to gain insights on job states and distribution across different regions and subscriptions, helping you manage resource allocation and monitor job performance.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which stream analytics jobs have failed, enabling you to focus on troubleshooting and rectifying those specific regions. This query is particularly useful for maintaining the efficiency of your Azure Stream Analytics.

```sql+postgres
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

```sql+sqlite
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