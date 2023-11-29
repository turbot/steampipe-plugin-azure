# Table: azure_reservation_recommendation

Azure Reservations help you save money by committing to one-year or three-year plans for multiple products. Committing allows you to get a discount on the resources you use. Reservations can significantly reduce your resource costs by up to 72% from pay-as-you-go prices. Reservations provide a billing discount and don't affect the runtime state of your resources. After you purchase a reservation, the discount automatically applies to matching resources.

**Note:** We can filter out the recommendations by using the columns `look_back_period`, `resource_type` or `scope` values in the query parameter. By default the table returns the data of resource type `VirtualMachines` with the scope `Single` for the last seven days of usage to look back for recommendation.

## Examples

### Basic info

```sql
select
  name,
  id,
  region,
  scope,
  etag,
  type
from
  azure_reservation_recommendation;
```

### Get reservation recommendation details for the last 30 days

```sql
select
  name,
  tags,
  sku,
  look_back_period
from
  azure_reservation_recommendation
where
  look_back_period = 'Last30Days';
```

### List reservation recommendation of the resource type MySQL

```sql
select
  name,
  tags,
  sku,
  look_back_period,
  resource_type
from
  azure_reservation_recommendation
where
  resource_type = 'MySQL';
```

### Get legacy resrvation recommendation properties

```sql
select
  name,
  id,
  legacy_recommendation_properties ->> 'LookBackPeriod' as look_back_period,
  legacy_recommendation_properties ->> 'InstanceFlexibilityRatio' as instance_flexibility_ratio,
  legacy_recommendation_properties ->> 'InstanceFlexibilityGroup' as instance_flexibility_group,
  legacy_recommendation_properties ->> 'NormalizedSize' as normalized_size,
  legacy_recommendation_properties ->> 'RecommendedQuantityNormalized' as recommended_quantity_normalized,
  legacy_recommendation_properties -> 'MeterID' as meter_id,
  legacy_recommendation_properties ->> 'ResourceType' as resource_type,
  legacy_recommendation_properties ->> 'Term' as term,
  legacy_recommendation_properties -> 'CostWithNoReservedInstances' as cost_with_no_reserved_instances,
  legacy_recommendation_properties -> 'RecommendedQuantity' as recommended_quantity,
  legacy_recommendation_properties -> 'TotalCostWithReservedInstances' as total_cost_with_reserved_instances,
  legacy_recommendation_properties -> 'NetSavings' as net_savings,
  legacy_recommendation_properties ->> 'FirstUsageDate' as first_usage_date,
  legacy_recommendation_properties ->> 'Scope' as scope,
  legacy_recommendation_properties -> 'SkuProperties' as sku_properties
from
  azure_reservation_recommendation
where
  kind = 'legacy';
```

### Get modern resrvation recommendation properties

```sql
select
  name,
  id,
  modern_recommendation_properties ->> 'Location' as location,
  modern_recommendation_properties ->> 'LookBackPeriod' as look_back_period,
  modern_recommendation_properties ->> 'InstanceFlexibilityRatio' as instance_flexibility_ratio,
  modern_recommendation_properties ->> 'InstanceFlexibilityGroup' as instance_flexibility_group,
  modern_recommendation_properties ->> 'NormalizedSize' as normalized_size,
  modern_recommendation_properties ->> 'RecommendedQuantityNormalized' as recommended_quantity_normalized,
  modern_recommendation_properties -> 'MeterID' as meter_id,
  modern_recommendation_properties ->> 'ResourceType' as resource_type,
  modern_recommendation_properties ->> 'Term' as term,
  modern_recommendation_properties -> 'CostWithNoReservedInstances' as cost_with_no_reserved_instances,
  modern_recommendation_properties -> 'RecommendedQuantity' as recommended_quantity,
  modern_recommendation_properties -> 'TotalCostWithReservedInstances' as total_cost_with_reserved_instances,
  modern_recommendation_properties -> 'NetSavings' as net_savings,
  modern_recommendation_properties ->> 'FirstUsageDate' as first_usage_date,
  modern_recommendation_properties ->> 'Scope' as scope,
  modern_recommendation_properties -> 'SkuProperties' as sku_properties,
  modern_recommendation_properties ->> 'SubscriptionID' as subscription_id,
  modern_recommendation_properties ->> 'ResourceType' as resource_type,
  modern_recommendation_properties ->> 'SkuName' as sku_name
from
  azure_reservation_recommendation
where
  kind = 'modern';
```