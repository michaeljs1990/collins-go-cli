# collins query - Search for assets in collins

## Usage

   collins query [command options] [arguments...]

## Options

  Global setting
   --timeout value           Timeout in seconds (0 == forever) (default: 0)
   -C value, --config value  Use specific Collins config yaml for client
   
  Query options
   -t value, --tag value             Assets with tag[s] value[,...]
   -Z, --remote-lookup               Query remote datacenters for asset
   -T value, --type value            Only show asset with type value
   -n value, --nodeclass value       Assets in nodeclass value[,...]
   -p value, --pool value            Assets in pool value[,...]
   -s value, --size value            Number of assets to return per page (default: 100)
   --limit value                     Limit total results of assets (default: 0)
   -r value, --role value            Assets in primary role
   -R value, --secondary-role value  Assets in secondary role
   -i value, --ip-address value      Assets with IP address[es]
   -S value, --status value          Asset status (and optional state after :)
   -a value, --attribute value       Arbitrary attributes and values to match in query. : between key and value
   -o value, --operation value       Sets if your query will be joined with AND or OR (default: "AND")
   
  Robot formatting
   -l, --link  Output link to assets found in web UI
   -j, --json  Output results in JSON
   -y, --yaml  Output results in YAML
   
  Table formatting
   -H, --show-header                  Show header fields in output
   -c value, --columns value          Attributes to output as columns, comma separated (default: "tag,hostname,nodeclass,status,pool,primary_role,secondary_role")
   -x value, --extra-columns value    Show these columns in addition to the default columns, comma separated
   -f value, --field-separator value  Separator between columns in output (default: "\t")

## Examples


