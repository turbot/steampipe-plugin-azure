# Table: azure_compute_virtual_machine_metric_cpu_utilization_daily

Azure Monitor metrics provide data about the performance of your systems. The `azure_compute_virtual_machine_metric_cpu_utilization_daily` table provides metric statistics at 24 hours intervals for the most recent 1 year.

## Examples

### Basic info

```sql
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization_daily
order by
  name,
  timestamp;
```

### CPU Over 80% average

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization_daily
where
  average > 80
order by
  name,
  timestamp;
```
