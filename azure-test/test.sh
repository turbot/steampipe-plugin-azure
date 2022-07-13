#!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 YELLOW="\033[0;33m"
 ENDCOLOR="\e[0m"
 
 
# Define your function here
run_test () {
   echo -e "${YELLOW}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 >> output.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BOLDGREEN}Passed -> $1 ${ENDCOLOR}"
    echo $1 >> passed_tests.txt
   fi
 }
 
 # output.txt - store output of each test
 # failed_tests.txt - names of failed test
 # passed_tests.txt names of passed test

 # removes files from previous test
# rm -rf output.txt failed_tests.txt passed_tests.txt
 date >> output.txt
 date >> failed_tests.txt
 date >> passed_tests.txt


# It is taking too much time(5-6 hrs)
# run_test azure_api_management
run_test azure_app_configuration
run_test azure_app_service_function_app
run_test azure_application_security_group
run_test azure_batch_account
run_test azure_compute_availability_set
run_test azure_compute_disk
run_test azure_compute_disk_access
run_test azure_compute_disk_encryption_set
run_test azure_compute_snapshot
run_test azure_compute_virtual_machine
run_test azure_container_registry
run_test azure_cosmosdb_account
run_test azure_cosmosdb_sql_database
run_test azure_data_factory
run_test azure_data_factory_dataset
run_test azure_data_factory_pipeline
run_test azure_data_lake_analytics_account
run_test azure_data_lake_store
run_test azure_eventgrid_topic
run_test azure_eventhub_namespace
run_test azure_express_route_circuit
run_test azure_hpc_cache
run_test azure_iothub
run_test azure_iothub_dps
run_test azure_key_vault_secret
run_test azure_kubernetes_cluster
run_test azure_kusto_cluster
run_test azure_lb
run_test azure_lb_backend_address_pool
run_test azure_lb_nat_rule
run_test azure_lb_outbound_rule
run_test azure_lb_probe
run_test azure_lb_rule
run_test azure_log_alert
run_test azure_log_profile
run_test azure_logic_app_workflow
run_test azure_management_lock
run_test azure_mssql_elasticpool
run_test azure_mssql_virtual_machine
run_test azure_mysql_flexible_server
run_test azure_network_interface
run_test azure_network_security_group
run_test azure_provider
run_test azure_public_ip
run_test azure_recovery_services_vault
run_test azure_redis_cache
run_test azure_resource_group
run_test azure_role_assignment
run_test azure_role_definition
run_test azure_route_table
run_test azure_search_service
run_test azure_security_center_automation
run_test azure_security_center_contact
run_test azure_service_fabric_cluster
run_test azure_servicebus_namespace
run_test azure_spring_cloud_service
run_test azure_sql_server
run_test azure_storage_account
run_test azure_storage_blob_service
run_test azure_storage_queue
run_test azure_storage_table
run_test azure_stream_analytics_job
run_test azure_subnet
run_test azure_subscription
run_test azure_synapse_workspace
run_test azure_tenant
run_test azure_virtual_network

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt