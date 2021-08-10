# Table: azure_stream_analytics_job

An Azure Stream Analytics job consists of an input, query, and an output. Stream Analytics ingests data from Azure Event Hubs (including Azure Event Hubs from Apache Kafka), Azure IoT Hub, or Azure Blob Storage. The query, which is based on SQL query language, can be used to easily filter, sort, aggregate, and join streaming data over a period of time.

## Examples

### Basic info

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