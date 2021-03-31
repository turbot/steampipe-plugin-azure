connection "azure" {
  plugin          = "azure"                 
  
  # If no subscription id is specified for a connection, the current active 
  # subscription per the `az` cli will be used.
  #subscription_id = "00000000-0000-0000-0000-000000000000"

  # If no credentials are specified and the SDK environment variables are not set, 
  # the plugin will use the active credentials from the `az` cli. You can run 
  # `az login` to set up these credentials.  For a full list of options, see the 
  # documentation at https://hub.steampipe.io/plugins/turbot/azure

}
