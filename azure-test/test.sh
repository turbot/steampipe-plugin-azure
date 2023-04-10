#!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 YELLOW="\033[0;33m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 ENDCOLOR="\e[0m"

# Define your function here
run_test () {
   echo -e "${YELLOW}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 > temp.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BOLDGREEN}Passed -> $1 ${ENDCOLOR}"
    echo $1 >> passed_tests.txt
   fi
  echo -e "$1" >> resource_list.txt && cat temp.txt | grep "resource_name" >> resource_list.txt && echo -e "\n\n" >> resource_list.txt
  cat temp.txt >> output.txt
  rm -rf temp.txt
 }

 # output.txt - store output of each test
 # failed_tests.txt - names of failed test
 # passed_tests.txt names of passed test

 # removes files from previous test
# rm -rf output.txt failed_tests.txt passed_tests.txt
 date >> output.txt
 date >> failed_tests.txt
 date >> passed_tests.txt
 date >> resource_list.txt

run_test azure_app_service_web_app
run_test azure_application_insight
run_test azure_bastion_host
run_test azure_cognitive_account
run_test azure_compute_availability_set
run_test azure_compute_disk_access
run_test azure_compute_virtual_machine
run_test azure_compute_virtual_machine_scale_set
run_test azure_cosmosdb_account
run_test azure_cosmosdb_sql_database
run_test azure_data_factory_dataset
run_test azure_express_route_circuit
run_test azure_hpc_cache
run_test azure_iothub
run_test azure_iothub_dps
run_test azure_key_vault_managed_hardware_security_module
run_test azure_kusto_cluster
run_test azure_logic_app_workflow
run_test azure_mariadb_server
run_test azure_mssql_elasticpool
run_test azure_mssql_virtual_machine
run_test azure_mysql_server
run_test azure_route_table
run_test azure_signalr_service
run_test azure_sql_database
run_test azure_storage_sync

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt