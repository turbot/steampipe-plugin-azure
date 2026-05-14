---
title: "Steampipe Table: azure_security_center_assessment - Query Azure Security Center Assessments using SQL"
description: "Allows users to query Azure Security Center Assessments, providing security assessment results for resources across a subscription."
folder: "Security Center"
---

# Table: azure_security_center_assessment - Query Azure Security Center Assessments using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers, and provides advanced threat protection across your hybrid workloads in the cloud - whether they're in Azure or not. Security assessments are the results of security checks run against your resources, indicating whether a resource is healthy, unhealthy, or not applicable for a given security control.

## Table Usage Guide

The `azure_security_center_assessment` table provides insights into the security assessment results for resources within Azure Security Center. As a security engineer, you can explore assessment status, resource details, and remediation information through this table. Utilize it to uncover unhealthy resources, track assessment status changes over time, and identify resources that require attention.

## Examples

### Basic info
Explore the various security assessments within Azure Security Center. This allows you to understand and categorize different assessment results by their unique identifiers, names, display names, and status, providing a comprehensive overview of your security posture.

```sql+postgres
select
  id,
  name,
  display_name,
  type,
  status_code
from
  azure_security_center_assessment;
```

```sql+sqlite
select
  id,
  name,
  display_name,
  type,
  status_code
from
  azure_security_center_assessment;
```

### List unhealthy assessments
Determine the areas in which resources are marked as unhealthy in Azure Security Center. This provides a way to identify instances where security measures may need to be improved or updated.

```sql+postgres
select
  name,
  display_name,
  status_code,
  status_cause,
  status_description
from
  azure_security_center_assessment
where
  status_code = 'Unhealthy';
```

```sql+sqlite
select
  name,
  display_name,
  status_code,
  status_cause,
  status_description
from
  azure_security_center_assessment
where
  status_code = 'Unhealthy';
```

### Get resource details for each assessment
Identify the specific resource that was assessed, including its type and ID. This is useful for correlating assessment results back to the underlying resource.

```sql+postgres
select
  name,
  display_name,
  resource_name,
  resource_id,
  resource_details
from
  azure_security_center_assessment;
```

```sql+sqlite
select
  name,
  display_name,
  resource_name,
  resource_id,
  resource_details
from
  azure_security_center_assessment;
```

### List assessments with their first evaluation and last status change dates
Track how long a resource has been in its current assessment status, and when it was first evaluated. This is useful for identifying long-standing issues.

```sql+postgres
select
  name,
  display_name,
  status_code,
  status_first_evaluation_date,
  status_change_date
from
  azure_security_center_assessment
order by
  status_change_date;
```

```sql+sqlite
select
  name,
  display_name,
  status_code,
  status_first_evaluation_date,
  status_change_date
from
  azure_security_center_assessment
order by
  status_change_date;
```

### Get links and metadata for an assessment by expanding the API response
By default, the `links` and `metadata` fields are not returned by the API. Use the `expand` qualifier to request them.

```sql+postgres
select
  name,
  display_name,
  links,
  metadata
from
  azure_security_center_assessment
where
  expand = 'links,metadata';
```

```sql+sqlite
select
  name,
  display_name,
  links,
  metadata
from
  azure_security_center_assessment
where
  expand = 'links,metadata';
```

### List additional data associated with assessments
Explore extra contextual data captured alongside each assessment, which can help with deeper investigation of a finding.

```sql+postgres
select
  name,
  display_name,
  additional_data
from
  azure_security_center_assessment
where
  additional_data is not null;
```

```sql+sqlite
select
  name,
  display_name,
  additional_data
from
  azure_security_center_assessment
where
  additional_data is not null;
```
