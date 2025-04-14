---
title: "Steampipe Table: azure_consumption_usage - Query Azure Cost Management using SQL"
description: "Allows users to query Azure Cost Management usage details for the defined scope."
folder: "Consumption"
---

# Table: azure_consumption_usage - Query Azure Consumption Usage using SQL

Azure Consumption Usage provide the ability to explore cost and usage data via multidimensional analysis, where creating customized filters and expressions allow you to answer consumption-related questions for your Azure resources.

## Table Usage Guide

The `azure_consumption_usage` table provides the comprehensive data about the resources and services you have used in your Azure subscription, along with the associated costs. This information is crucial for managing and optimizing Azure costs, understanding billing, and monitoring resource utilization.

**Important notes:**
- By default this table returns the result for subscription scope.
- This table can provide consumption usage details for the previous one year.
- For improved performance, it is advised that you use the optional qual `filter` to limit the result set to a specific time period .
- This table supports optional quals. Queries with optional quals are optimized to use Consumption Usage filters. Optional quals are supported for the following columns:
  - `filter`: May be used to filter usageDetails by properties/resourceGroup, properties/instanceName, properties/resourceId, properties/chargeType, properties/reservationId, properties/publisherType or tags. The filter supports 'eq', 'lt', 'gt', 'le', 'ge', and 'and'. It does not currently support 'ne', 'or', or 'not'. Tag filter is a key value pair string where key and value is separated by a colon (:). PublisherType Filter accepts two values azure and marketplace and it is currently supported for Web Direct Offer Type."
  - `metric`: Allows to select different type of cost/usage records. Possible values are 'actualcost', 'amortizedcost' or 'usage'.
  - `scope`: The scope associated with usage details operations. This includes '/subscriptions/{subscriptionId}/' for subscription scope, '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}' for Billing Account scope, '/providers/Microsoft.Billing/departments/{departmentId}' for Department scope, '/providers/Microsoft.Billing/enrollmentAccounts/{enrollmentAccountId}' for EnrollmentAccount scope and '/providers/Microsoft.Management/managementGroups/{managementGroupId}' for Management Group scope. For subscription, billing account, department, enrollment account and management group, you can also add billing period to the scope using '/providers/Microsoft.Billing/billingPeriods/{billingPeriodName}'. For e.g. to specify billing period at department scope use '/providers/Microsoft.Billing/departments/{departmentId}/providers/Microsoft.Billing/billingPeriods/{billingPeriodName}'. Also, Modern Commerce Account scopes are '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}' for billingAccount scope, '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}' for billingProfile scope, 'providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}/invoiceSections/{invoiceSectionId}' for invoiceSection scope, and 'providers/Microsoft.Billing/billingAccounts/{billingAccountId}/customers/{customerId}' specific for partners.
  - `expand`: May be used to expand the 'properties/additionalInfo' or 'properties/meterDetails' within a list of usage details. By default, these fields are not included when listing usage details.

## Examples

### Basic info
This query is useful for getting a broad overview of resource consumption in Azure. It can help in cost management, resource optimization, and understanding how different Azure resources are being utilized. The data retrieved can be instrumental for in-depth analysis, especially when dealing with complex Azure environments with multiple resources and services.

```sql+postgres
select
  name,
  id,
  scope,
  kind,
  etag,
  type
from
  azure_consumption_usage;
```

```sql+sqlite
select
  name,
  id,
  scope,
  kind,
  etag,
  type
from
  azure_consumption_usage;
```

### Get legacy consumption usage in a subscription
 This is beneficial for organizations looking to get insights into their legacy resource usage in Azure, aiding in decision-making regarding migration, cost management, and resource optimization.

```sql+postgres
select
  name,
  id,
  scope,
  kind,
  etag,
  type
from
  azure_consumption_usage
where
  kind = 'legacy';
```

```sql+sqlite
select
  name,
  id,
  scope,
  kind,
  etag,
  type
from
  azure_consumption_usage
where
  kind = 'legacy';
```

### Filter actual cost stastics of legacy consumption usage
Extract detailed insights into their Azure consumption, focusing on legacy resources and actual cost metrics. It aids in financial management, resource optimization, and strategic planning in the context of Azure cloud services.

```sql+postgres
select
  name,
  id,
  metric,
  kind,
  legacy_usage_detail ->> 'BillingAccountID' as billing_account_id,
  legacy_usage_detail ->> 'BillingAccountName' as billing_account_name,
  legacy_usage_detail ->> 'BillingPeriodStartDate' as billing_period_start_date,
  legacy_usage_detail ->> 'BillingPeriodEndDate' as billing_period_end_date,
  legacy_usage_detail ->> 'Product' as product,
  legacy_usage_detail ->> 'Quantity' as quantity,
  legacy_usage_detail ->> 'Cost' as cost,
  legacy_usage_detail ->> 'BillingCurrency' as billing_currency,
  legacy_usage_detail ->> 'ChargeType' as charge_type,
  legacy_usage_detail ->> 'IsAzureCreditEligible' as is_azure_credit_eligible,
  legacy_usage_detail ->> 'ResourceID' as resource_id
from
  azure_consumption_usage
where
  kind = 'legacy'
  and metric = 'actualcost';
```

```sql+sqlite
select
  name,
  id,
  metric,
  kind,
  json_extract(legacy_usage_detail, '$.BillingAccountID') as billing_account_id,
  json_extract(legacy_usage_detail, '$.BillingAccountName') as billing_account_name,
  json_extract(legacy_usage_detail, '$.BillingPeriodStartDate') as billing_period_start_date,
  json_extract(legacy_usage_detail, '$.BillingPeriodEndDate') as billing_period_end_date,
  json_extract(legacy_usage_detail, '$.Product') as product,
  json_extract(legacy_usage_detail, '$.Quantity') as quantity,
  json_extract(legacy_usage_detail, '$.Cost') as cost,
  json_extract(legacy_usage_detail, '$.BillingCurrency') as billing_currency,
  json_extract(legacy_usage_detail, '$.ChargeType') as charge_type,
  json_extract(legacy_usage_detail, '$.IsAzureCreditEligible') as is_azure_credit_eligible,
  json_extract(legacy_usage_detail, '$.ResourceID') as resource_id
from
  azure_consumption_usage
where
  kind = 'legacy'
  and metric = 'actualcost';
```

### Get top 10 legacy consumption usages in a year
Analyze the distribution of Azure container groups based on their operating system type. This can help in understanding the usage pattern of different OS types within your Azure container groups.

```sql+postgres
select
  name,
  id,
  metric,
  kind,
  legacy_usage_detail ->> 'BillingAccountID' as billing_account_id,
  legacy_usage_detail ->> 'BillingAccountName' as billing_account_name,
  legacy_usage_detail ->> 'BillingPeriodStartDate' as billing_period_start_date,
  legacy_usage_detail ->> 'BillingPeriodEndDate' as billing_period_end_date,
  legacy_usage_detail ->> 'Cost' as cost,
  legacy_usage_detail ->> 'BillingCurrency' as billing_currency
from
  azure_consumption_usage
where
  kind = 'legacy'
and
  metric = 'actualcost'
order by
  cost desc limit 10;
```

```sql+sqlite
select
  name,
  id,
  metric,
  kind,
  json_extract(legacy_usage_detail, '$.BillingAccountID') as billing_account_id,
  json_extract(legacy_usage_detail, '$.BillingAccountName') as billing_account_name,
  json_extract(legacy_usage_detail, '$.BillingPeriodStartDate') as billing_period_start_date,
  json_extract(legacy_usage_detail, '$.BillingPeriodEndDate') as billing_period_end_date,
  json_extract(legacy_usage_detail, '$.Cost') as cost,
  json_extract(legacy_usage_detail, '$.BillingCurrency') as billing_currency
from
  azure_consumption_usage
where
  kind = 'legacy'
and
  metric = 'actualcost'
order by
  cost desc limit 10;
```

### Filter consumption usage by resource group
Discover the segments that provide information about IP addresses associated with each group. This is useful in understanding the network connectivity and accessibility of these groups within the Azure container ecosystem.

```sql+postgres
select
  name,
  id,
  metric,
  kind,
  legacy_usage_detail ->> 'BillingAccountID' as billing_account_id,
  legacy_usage_detail ->> 'BillingAccountName' as billing_account_name,
  legacy_usage_detail ->> 'BillingPeriodStartDate' as billing_period_start_date,
  legacy_usage_detail ->> 'BillingPeriodEndDate' as billing_period_end_date,
  legacy_usage_detail ->> 'Cost' as cost,
  legacy_usage_detail ->> 'BillingCurrency' as billing_currency,
  legacy_usage_detail ->> 'ResourceID' as resource_id
from
  azure_consumption_usage
where
  kind = 'legacy'
  and metric = 'actualcost'
  and filter = 'properties/resourceGroup eq ''turbot_rg''';
```

```sql+sqlite
select
  name,
  id,
  metric,
  kind,
  json_extract(legacy_usage_detail, '$.BillingAccountID') as billing_account_id,
  json_extract(legacy_usage_detail, '$.BillingAccountName') as billing_account_name,
  json_extract(legacy_usage_detail, '$.BillingPeriodStartDate') as billing_period_start_date,
  json_extract(legacy_usage_detail, '$.BillingPeriodEndDate') as billing_period_end_date,
  json_extract(legacy_usage_detail, '$.Cost') as cost,
  json_extract(legacy_usage_detail, '$.BillingCurrency') as billing_currency
from
  azure_consumption_usage
where
  kind = 'legacy'
  and metric = 'actualcost'
  and filter = 'properties/resourceGroup eq ''turbot_rg''';
```