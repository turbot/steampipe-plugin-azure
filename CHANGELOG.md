## v0.20.0 [2021-10-26]

_Enhancements_

- Updated: Add context cancellation handling to all the tables ([#343](https://github.com/turbot/steampipe-plugin-azure/pull/343))
- Recompiled plugin with [steampipe-plugin-sdk v1.7.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v170--2021-10-18) ([#400](https://github.com/turbot/steampipe-plugin-azure/pull/400))
- The configuration section in the docs/index.md file now includes additional information on different methods of setting up credentials in the `azure.spc` file

_Bug fixes_

- Fixed: Authentication now works properly for non-public cloud environments

_Deprecated_

- The following tables have been deprecated since they are now maintained in the [azuread plugin](https://hub.steampipe.io/plugins/turbot/azuread/tables). These tables will be removed in the next major version. We recommend updating any scripts or workflows that use these tables to use the equivalent tables in the Azure AD plugin instead
  - [azure_ad_group](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_ad_group) (replaced by [azuread_group](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_group))
  - [azure_ad_service_principal](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_ad_service_principal) (replaced by [azuread_service_principal](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_service_principal))
  - [azure_ad_user](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_ad_user) (replaced by [azuread_user](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_user))

## v0.19.0 [2021-10-07]
  
_What's new?_

- New tables added
  - [azure_app_configuration](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_app_configuration) ([#344](https://github.com/turbot/steampipe-plugin-azure/pull/344))
  - [azure_application_gateway](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_application_gateway) ([#316](https://github.com/turbot/steampipe-plugin-azure/pull/316))
  - [azure_cognitive_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_cognitive_account) ([#301](https://github.com/turbot/steampipe-plugin-azure/pull/301))
  - [azure_compute_disk_access](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_access) ([#288](https://github.com/turbot/steampipe-plugin-azure/pull/288))
  - [azure_databox_edge_device](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_databox_edge_device) ([#377](https://github.com/turbot/steampipe-plugin-azure/pull/377))
  - [azure_eventgrid_domain](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_eventgrid_domain) ([#314](https://github.com/turbot/steampipe-plugin-azure/pull/314))
  - [azure_eventgrid_topic](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_eventgrid_topic) ([#352](https://github.com/turbot/steampipe-plugin-azure/pull/352))
  - [azure_frontdoor](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_frontdoor) ([#362](https://github.com/turbot/steampipe-plugin-azure/pull/362))
  - [azure_hdinsight_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hdinsight_cluster) ([#395](https://github.com/turbot/steampipe-plugin-azure/pull/395))
  - [azure_healthcare_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_healthcare_service) ([#345](https://github.com/turbot/steampipe-plugin-azure/pull/345))
  - [azure_hpc_cache](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hpc_cache) ([#374](https://github.com/turbot/steampipe-plugin-azure/pull/374))
  - [azure_hybrid_compute_machine](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hybrid_compute_machine) ([#378](https://github.com/turbot/steampipe-plugin-azure/pull/378))
  - [azure_hybrid_kubernetes_connected_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hybrid_kubernetes_connected_cluster) ([#376](https://github.com/turbot/steampipe-plugin-azure/pull/376))
  - [azure_iothub_dps](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_iothub_dps) ([#363](https://github.com/turbot/steampipe-plugin-azure/pull/363))
  - [azure_kusto_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_kusto_cluster) ([#369](https://github.com/turbot/steampipe-plugin-azure/pull/369))
  - [azure_machine_learning_workspace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_machine_learning_workspace) ([#315](https://github.com/turbot/steampipe-plugin-azure/pull/315))
  - [azure_mssql_virtual_machine](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mssql_virtual_machine) ([#313](https://github.com/turbot/steampipe-plugin-azure/pull/313))
  - [azure_service_fabric_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_service_fabric_cluster) ([#310](https://github.com/turbot/steampipe-plugin-azure/pull/310))
  - [azure_signalr_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_signalr_service) ([#328](https://github.com/turbot/steampipe-plugin-azure/pull/328))
  - [azure_spring_cloud_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_spring_cloud_service) ([#347](https://github.com/turbot/steampipe-plugin-azure/pull/347))
  - [azure_storage_sync](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_storage_sync) ([#326](https://github.com/turbot/steampipe-plugin-azure/pull/326))
  - [azure_synapse_workspace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_synapse_workspace) ([#346](https://github.com/turbot/steampipe-plugin-azure/pull/346))

_Enhancements_

- Added `encryption_scope` column to `azure_storage_account` table ([#392](https://github.com/turbot/steampipe-plugin-azure/pull/392))
- Added `security_profile` column to `azure_compute_virtual_machine` table ([#387](https://github.com/turbot/steampipe-plugin-azure/pull/387))
- Added `public_network_access` and `minimal_tls_version` column to `azure_sql_server` table ([#371](https://github.com/turbot/steampipe-plugin-azure/pull/371))
- Added `guest_configuration_assignments` column to `azure_compute_virtual_machine` table ([#353](https://github.com/turbot/steampipe-plugin-azure/pull/353)) ([#380](https://github.com/turbot/steampipe-plugin-azure/pull/380))
- Added `cluster_settings` column to `azure_app_service_environment` table ([#360](https://github.com/turbot/steampipe-plugin-azure/pull/360))
- Added `developer_portal_url`, `disable_gateway`, `enable_client_certificate`, `api_version_constraint`, `certificates`, `custom_properties`, `identity_user_assigned_identities`, `virtual_network_type`, `restore`, `scm_url`, `zones` and `diagnostic_settings` columns to `azure_api_management` table ([#336](https://github.com/turbot/steampipe-plugin-azure/pull/336))
- Added `security_alert_policies` column to `azure_mssql_managed_instance` table ([#333](https://github.com/turbot/steampipe-plugin-azure/pull/333))
- Added `private_endpoint_connections` column to `azure_eventhub_namespace` table ([#331](https://github.com/turbot/steampipe-plugin-azure/pull/331))
- Added `private_endpoint_connections` column to `azure_servicebus_namespace` table ([#334](https://github.com/turbot/steampipe-plugin-azure/pull/334))
- Added `server_keys` column to `azure_mysql_server` table ([#337](https://github.com/turbot/steampipe-plugin-azure/pull/337))
- Added `identity` column to `azure_compute_virtual_machine` table ([#341](https://github.com/turbot/steampipe-plugin-azure/pull/341))
- Added `server_keys` column to `azure_postgresql_server` table ([#299](https://github.com/turbot/steampipe-plugin-azure/pull/299))
- Added `private_endpoint_connections` column to `azure_key_vault` table ([#306](https://github.com/turbot/steampipe-plugin-azure/pull/306)) ([#342](https://github.com/turbot/steampipe-plugin-azure/pull/342))
- Added `private_endpoint_connections` column to `azure_sql_server` table ([#300](https://github.com/turbot/steampipe-plugin-azure/pull/300))
- Added `private_endpoint_connections` column to `azure_data_factory` table ([#298](https://github.com/turbot/steampipe-plugin-azure/pull/298))

_Bug fixes_

- Querying column `encryption` in table `azure_servicebus_namespace` now render all available properties ([#366](https://github.com/turbot/steampipe-plugin-azure/pull/366))
- Querying column `cluster_arm_id` in table `azure_eventhub_namespace` no longer return `nil` if available ([#351](https://github.com/turbot/steampipe-plugin-azure/pull/351))
- Querying column `private_endpoint_connections` in table `azure_mysql_server` now render all available properties ([#338](https://github.com/turbot/steampipe-plugin-azure/pull/338))
- Querying column `private_endpoint_connections` in table `azure_postgresql_server` now render all available properties ([#339](https://github.com/turbot/steampipe-plugin-azure/pull/339))
- Querying column `encryption_protector` in table `azure_sql_server` now render all available properties ([#361](https://github.com/turbot/steampipe-plugin-azure/pull/361))


## v0.18.0 [2021-08-25]

_What's new?_

- New tables added
  - [azure_key_vault_deleted_vault](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_deleted_vault) ([#263](https://github.com/turbot/steampipe-plugin-azure/pull/263))
  - [azure_lb_backend_address_pool](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_backend_address_pool) ([#266](https://github.com/turbot/steampipe-plugin-azure/pull/266))
  - [azure_lb_nat_rule](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_nat_rule) ([#267](https://github.com/turbot/steampipe-plugin-azure/pull/267))
  - [azure_lb_outbound_rule](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_outbound_rule) ([#264](https://github.com/turbot/steampipe-plugin-azure/pull/264))
  - [azure_mssql_elasticpool](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mssql_elasticpool) ([#276](https://github.com/turbot/steampipe-plugin-azure/pull/276))
  - [azure_mssql_managed_instance](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mssql_managed_instance) ([#277](https://github.com/turbot/steampipe-plugin-azure/pull/277))

_Enhancements_

- Updated: Add `vulnerability_assessments` and `vulnerability_assessment_scan_records` columns in `azure_sql_database` table ([#279](https://github.com/turbot/steampipe-plugin-azure/pull/279))


## v0.17.0 [2021-08-13]

_What's new?_

- New tables added
  - [azure_compute_virtual_machine_scale_set](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_scale_set) ([#249](https://github.com/turbot/steampipe-plugin-azure/pull/249))
  - [azure_data_lake_analytics_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_lake_analytics_account) ([#253](https://github.com/turbot/steampipe-plugin-azure/pull/253))
  - [azure_lb](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb) ([#231](https://github.com/turbot/steampipe-plugin-azure/pull/231))
  - [azure_lb_probe](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_probe) ([#238](https://github.com/turbot/steampipe-plugin-azure/pull/238))
  - [azure_lb_rule](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_rule) ([#235](https://github.com/turbot/steampipe-plugin-azure/pull/235))
  - [azure_resource_link](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_resource_link) ([#252](https://github.com/turbot/steampipe-plugin-azure/pull/252))
  - [azure_search_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_search_service) ([#257](https://github.com/turbot/steampipe-plugin-azure/pull/257))

_Enhancements_

- Updated: Add `retention_policy_id`, `retention_policy_name`, `retention_policy_type` and `retention_policy_property` columns in `azure_sql_database` table ([#255](https://github.com/turbot/steampipe-plugin-azure/pull/255))

_Bug fixes_

- Fixed: Integration test issues for several tables ([#259](https://github.com/turbot/steampipe-plugin-azure/pull/259))
- Fixed: Expired CLI authentication tokens will now automatically be refreshed ([#234](https://github.com/turbot/steampipe-plugin-azure/pull/234))

## v0.16.0 [2021-08-11]

_What's new?_

- New tables added
  - [azure_batch_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_batch_account) ([#242](https://github.com/turbot/steampipe-plugin-azure/pull/242))
  - [azure_data_lake_store](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_lake_store) ([#217](https://github.com/turbot/steampipe-plugin-azure/pull/217))
  - [azure_iothub](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_iothub) ([#232](https://github.com/turbot/steampipe-plugin-azure/pull/232))
  - [azure_key_vault_managed_hardware_security_module](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_managed_hardware_security_module) ([#236](https://github.com/turbot/steampipe-plugin-azure/pull/236))
  - [azure_logic_app_workflow](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_logic_app_workflow) ([#230](https://github.com/turbot/steampipe-plugin-azure/pull/230))
  - [azure_recovery_services_vault](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_recovery_services_vault) ([#243](https://github.com/turbot/steampipe-plugin-azure/pull/243))
  - [azure_stream_analytics_job](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_stream_analytics_job) ([#246](https://github.com/turbot/steampipe-plugin-azure/pull/246))

_Enhancements_

- Updated: Add `diagnostic_settings` column in `azure_network_security_group` table ([#247](https://github.com/turbot/steampipe-plugin-azure/pull/247))
- Updated: Add `ExtensionType` property in `extensions` column for `azure_compute_virtual_machine` table ([#229](https://github.com/turbot/steampipe-plugin-azure/pull/229))
- Updated: Add `enable_automatic_updates`, `provision_vm_agent_windows`, `time_zone`, `additional_unattend_content`, `patch_settings` and `win_rm` columns in `azure_compute_virtual_machine` table ([#223](https://github.com/turbot/steampipe-plugin-azure/pull/223))
- Updated: Add `diagnostic_settings` column in `azure_servicebus_namespace` table ([#225](https://github.com/turbot/steampipe-plugin-azure/pull/225))
- Updated: Add `diagnostic_settings` column in `azure_eventhub_namespace` table ([#226](https://github.com/turbot/steampipe-plugin-azure/pull/226))
- Updated: Add `network_acls` column in `azure_key_vault` table ([#220](https://github.com/turbot/steampipe-plugin-azure/pull/220))
- Updated: Add `virtual_network_rules` column in `azure_sql_server` table ([#227](https://github.com/turbot/steampipe-plugin-azure/pull/227))
- Updated: Recompiled plugin with [steampipe-plugin-sdk v1.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v141--2021-07-20) ([#207](https://github.com/turbot/steampipe-plugin-azure/pull/207))

_Bug fixes_

- Fixed: Pagination for listing resources for all tables ([#254](https://github.com/turbot/steampipe-plugin-azure/pull/254))<br/>
  _This bug impacted all the tables in plugin. Now tables will not go into infinite loop for large number of resources._
- Fixed: Improve properties of `firewall_rules` column in `azure_sql_server` table ([#237](https://github.com/turbot/steampipe-plugin-azure/pull/237))

## v0.15.0 [2021-07-31]

_What's new?_

- New tables added
  - [azure_container_registry](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_container_registry) ([#196](https://github.com/turbot/steampipe-plugin-azure/pull/196))
  - [azure_eventhub_namespace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_eventhub_namespace) ([#201](https://github.com/turbot/steampipe-plugin-azure/pull/201))
  - [azure_mariadb_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mariadb_server) ([#206](https://github.com/turbot/steampipe-plugin-azure/pull/206))
  - [azure_redis_cache](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_redis_cache) ([#203](https://github.com/turbot/steampipe-plugin-azure/pull/203))
  - [azure_security_center_jit_network_access_policy](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_jit_network_access_policy) ([#192](https://github.com/turbot/steampipe-plugin-azure/pull/192))
  - [azure_servicebus_namespace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_servicebus_namespace) ([#200](https://github.com/turbot/steampipe-plugin-azure/pull/200))

_Enhancements_

- Updated: Add column `vnet_connection` to `azure_app_service_web_app` table ([#204](https://github.com/turbot/steampipe-plugin-azure/pull/204))

## v0.14.0 [2021-07-22]

_What's new?_

- New tables added
  - [azure_data_factory_pipeline](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_factory_pipeline) ([#169](https://github.com/turbot/steampipe-plugin-azure/pull/169))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v131--2021-07-15)

_Bug fixes_

- Fixed: `azure_virtual_network_gateway` table's parent hydrate now lists resource groups instead of virtual networks to prevent duplicate rows ([#181](https://github.com/turbot/steampipe-plugin-azure/pull/181))

## v0.13.0 [2021-07-13]

_What's new?_

- New tables added
  - [azure_compute_disk_metric_read_ops](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_read_ops) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_read_ops_daily](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_read_ops_daily) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_read_ops_hourly](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_read_ops_hourly) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_write_ops](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_write_ops) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_write_ops_daily](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_write_ops_daily) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_write_ops_hourly](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_write_ops_hourly) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_virtual_machine_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_cpu_utilization) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_virtual_machine_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_cpu_utilization_daily) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_virtual_machine_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_cpu_utilization_hourly) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_data_factory](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_factory) ([#158](https://github.com/turbot/steampipe-plugin-azure/pull/158))
  - [azure_data_factory_dataset](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_factory_dataset) ([#168](https://github.com/turbot/steampipe-plugin-azure/pull/168))
  - [azure_express_route_circuit](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_express_route_circuit) ([#170](https://github.com/turbot/steampipe-plugin-azure/pull/170))
  - [azure_virtual_network_gateway](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_virtual_network_gateway) ([#157](https://github.com/turbot/steampipe-plugin-azure/pull/157))

_Enhancements_

- Updated: `azure-sdk-for-go` to `v55.4.0+incompatible` ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
- Updated: Change several metric table function names to be consistent with naming standards ([#173](https://github.com/turbot/steampipe-plugin-azure/pull/173))

_Bug fixes_

- Fixed: Integration tests for several tables and remove unused integration tests ([#175](https://github.com/turbot/steampipe-plugin-azure/pull/175))

## v0.12.0 [2021-07-01]

_Enhancements_

- Updated: Add `lifecycle_management_policy` column to `azure_storage_account` table ([#155](https://github.com/turbot/steampipe-plugin-azure/pull/155))

## v0.11.0 [2021-06-03]

_What's new?_

- New tables added
  - [azure_policy_definition](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_policy_definition) ([#141](https://github.com/turbot/steampipe-plugin-azure/pull/141))
  - [azure_tenant](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_tenant) ([#142](https://github.com/turbot/steampipe-plugin-azure/pull/142))

## v0.10.0 [2021-05-27]

_What's new?_

- Updated plugin license to Apache 2.0 per [turbot/steampipe#488](https://github.com/turbot/steampipe/issues/488)
- New tables added
  - [azure_storage_blob](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_storage_blob) ([#133](https://github.com/turbot/steampipe-plugin-azure/pull/133))

_Bug fixes_

- Fixed: Improved error messages when we fail to get credentials from the Azure CLI ([#137](https://github.com/turbot/steampipe-plugin-azure/pull/137))

## v0.9.0 [2021-05-20]

_What's new?_

- New tables added
  - [azure_policy_assignment](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_policy_assignment) ([#123](https://github.com/turbot/steampipe-plugin-azure/pull/123))
  - [azure_security_center_setting](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_setting) ([#115](https://github.com/turbot/steampipe-plugin-azure/pull/115))
  - [azure_security_center_subscription_pricing](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_subscription_pricing) ([#135](https://github.com/turbot/steampipe-plugin-azure/pull/135))
  - [azure_subscription](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_subscription) ([#132](https://github.com/turbot/steampipe-plugin-azure/pull/132))

## v0.8.0 [2021-05-13]

_What's new?_

- New tables added
  - [azure_security_center_auto_provisioning](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_auto_provisioning) ([#117](https://github.com/turbot/steampipe-plugin-azure/pull/117))
  - [azure_security_center_contact](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_contact) ([#121](https://github.com/turbot/steampipe-plugin-azure/pull/121))

_Enhancements_

- Updated: README.md and docs/index.md now contain links to our Slack community ([#129](https://github.com/turbot/steampipe-plugin-azure/pull/129))
- Updated: Bump lodash from 4.17.20 to 4.17.21 in /azure-test ([#126](https://github.com/turbot/steampipe-plugin-azure/pull/126))

## v0.7.0 [2021-05-06]

_What's new?_

- New tables added
  - [azure_kubernetes_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_kubernetes_cluster) ([#83](https://github.com/turbot/steampipe-plugin-azure/pull/83))

_Enhancements_

- Updated: Add `identity` column to `azure_app_service_web_app` table ([#90](https://github.com/turbot/steampipe-plugin-azure/pull/90))
- Updated: Add `os_disk_vhd_uri` column to `azure_compute_virtual_machine` table ([#88](https://github.com/turbot/steampipe-plugin-azure/pull/88))

_Bug fixes_

- Fixed: The get calls in the `azure_key_vault_secret` table should not fail for disabled secrets ([#111](https://github.com/turbot/steampipe-plugin-azure/pull/111))

## v0.6.0 [2021-04-29]

_What's new?_

- New tables added
  - [azure_key_vault_key](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_key) ([#73](https://github.com/turbot/steampipe-plugin-azure/pull/73))
  - [azure_network_watcher_flow_log](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_network_watcher_flow_log) ([#61](https://github.com/turbot/steampipe-plugin-azure/pull/61))

_Enhancements_

- Updated: Add `extensions` column to `azure_compute_virtual_machine` table ([#85](https://github.com/turbot/steampipe-plugin-azure/pull/85))

_Bug fixes_

- Fixed: `azure_key_vault` table queries no longer crash when getting vault diagnostic settings ([#107](https://github.com/turbot/steampipe-plugin-azure/pull/107))
- Fixed: `deleted_time` and `last_modified_time` columns now show the correct dates in `azure_storage_container` table ([#106](https://github.com/turbot/steampipe-plugin-azure/pull/106))
- Fixed: `encryption_key_vault_properties_last_rotation_time` column now shows the correct date in `azure_storage_account` table ([#101](https://github.com/turbot/steampipe-plugin-azure/pull/101))
- Fixed: `subscription_id` column now displays the correct subscription ID in `azure_diagnostic_setting` table ([#99](https://github.com/turbot/steampipe-plugin-azure/pull/99))
- Fixed: Column name in example query in `azure_key_vault_secret` table doc ([#108](https://github.com/turbot/steampipe-plugin-azure/pull/108))

## v0.5.0 [2021-04-22]

_What's new?_

- New tables added
  - [azure_diagnostic_setting](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_diagnostic_setting) ([#70](https://github.com/turbot/steampipe-plugin-azure/pull/70))
  - [azure_key_vault_secret](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_secret) ([#76](https://github.com/turbot/steampipe-plugin-azure/pull/76))
  - [azure_log_alert](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_log_alert) ([#70](https://github.com/turbot/steampipe-plugin-azure/pull/70))
  - [azure_log_profile](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_log_profile) ([#70](https://github.com/turbot/steampipe-plugin-azure/pull/70))
  - [azure_sql_database](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_sql_database) ([#93](https://github.com/turbot/steampipe-plugin-azure/pull/93))

_Enhancements_

- Updated: Add `diagnostic_settings` column to `azure_key_vault` table ([#96](https://github.com/turbot/steampipe-plugin-azure/pull/96))
- Updated: Add `queue` prefix to various queue logging columns in `azure_storage_account` table ([#94](https://github.com/turbot/steampipe-plugin-azure/pull/94))

_Bug fixes_

- Fixed: List calls should not infinitely loop in `azure_key_vault` table ([#82](https://github.com/turbot/steampipe-plugin-azure/pull/82))

## v0.4.0 [2021-04-08]

_What's new?_

- New tables added
  - [azure_mysql_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mysql_server) ([#66](https://github.com/turbot/steampipe-plugin-azure/pull/66))
  - [azure_postgresql_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_postgresql_server) ([#64](https://github.com/turbot/steampipe-plugin-azure/pull/64))
  - [azure_storage_container](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_storage_container) ([#71](https://github.com/turbot/steampipe-plugin-azure/pull/71))

_Enhancements_

- Updated: Add `auth_settings` and `configuration` columns to `azure_app_service_function_app` table ([#77](https://github.com/turbot/steampipe-plugin-azure/pull/77))
- Updated: Add `auth_settings` and `configuration` columns to `azure_app_service_web_app` table ([#77](https://github.com/turbot/steampipe-plugin-azure/pull/77))
- Updated: Add `blob_service_logging` column to `azure_storage_account` table ([#80](https://github.com/turbot/steampipe-plugin-azure/pull/80))

_Bug fixes_

- Fixed: The table `azure_sql_server` should return `null` instead of an empty object for columns with missing data ([#68](https://github.com/turbot/steampipe-plugin-azure/pull/68))

## v0.3.1 [2021-03-25]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v0.2.6](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v026-2021-03-18)

_Documentation_

- Fixed various example queries in `azure_app_service_web_app` table document ([#59](https://github.com/turbot/steampipe-plugin-azure/pull/59))

## v0.3.0 [2021-03-11]

_What's new?_

- New tables added
  - [azure_sql_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_sql_server)

_Bug fixes_

- Removed use of deprecated `ItemFromKey` function from all tables

## v0.2.2 [2021-03-02]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve issue:
  - Fix tables failing with error similar to `Error: pq: rpc error: code = Internal desc = get hydrate function getAdGroup failed with panic interface conversion: interface {} is nil, not *graphrbac.ADGroup`([#29](https://github.com/turbot/steampipe-plugin-azure/issues/29)).

## v0.2.1 [2021-02-25]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve SDK issues:
  - Fix error for missing required quals [#40](https://github.com/turbot/steampipe-plugin-sdk/issues/42).
  - Queries fail with error socket: too many open files [#190](https://github.com/turbot/steampipe/issues/190)

## v0.2.0 [2021-02-18]

_What's new?_

- Added support for [connection configuration](https://github.com/turbot/steampipe-plugin-azure/blob/main/docs/index.md#connection-configuration). You may specify azure `tenant_id`, `subscription_id`, `client_id`, `client_secret`, `client_certificate_path`, `client_certificate_password`, `username` and `password` for each connection in a configuration file. You can have multiple azure connections, each configured for a different azure subscription.

_Enhancements_

- Updates
  - Added columns power_state, private_ips and public_ips to azure_compute_virtual_machine table ([#17](https://github.com/turbot/steampipe-plugin-azure/pull/17))

_Bug fixes_

- Breaking changes

  - Renamed earlier `azure_storage_blob` table to `azure_storage_blob_service` ([#7](https://github.com/turbot/steampipe-plugin-azure/pull/7))
  - Renamed earlier `azure_storage_table` table to `azure_storage_table_service` ([#10](https://github.com/turbot/steampipe-plugin-azure/pull/10))
  - Removed columns managed_disk_storage_account_type and os_disk_size_gb from `azure_compute_virtual_machine` table ([#17](https://github.com/turbot/steampipe-plugin-azure/pull/17))
