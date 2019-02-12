% Collins Golang CLI
% Michael Schuett
% February 12, 2019

Synopsis
========

`collins query` [options] [hostname]...

Description
===========

Collins query provides a wrapper round the collins find endpoint. Most
input is just taken from passed in flags and converted into a collins
query using CQL. Multiple output formats are supported such as Table,
YAML, JSON, and Links. Table which is the default output format can be
used to pipe the output into other commands such as set, power, and
provision.

Query Examples
------------

A basic example to query a specific tag and additionally tell us the ip of the
asset would be:

    collins query -t M000001 -x ip_address

Options
=======
   
Query options {.options}
---------------

`-t` *VALUES*, `--tag` *VALUES*

:   Specify tag names to return they must be the full tag and in the case that
    multiple are specified they must be separated by commas.

`-Z`, `--remote-lookup`

:   Remote lookup will utilize collins remote lookup feature to query not only
    the current collins instances but other instances it has been linked with
    as well.

`-T` *VALUES*, `--type` *VALUES*

:   Specify the type of asset that should be returned. Common uses for this
    might be for selecting all power switches or server nodes. When multiple
    values are specified they must be separated by commas.

        -T SERVER_NODE
        --type POWER_SWITCH,SERVER_NODE

`-n` *VALUES*, `--nodeclass` *VALUES*

:   Specify the nodeclass to query. When multiple values are specified they
    must be separated by commas. By default this does a wildcard match at 
    the start and end of the string so a value of `dev` would also match
    `clusterdevnode`.

`-p` *VALUES*, `--pool` *VALUES*

:   Specify the pool to query. When multiple values are specified they must be
    separated by commas. By default this does a wildcard match at the start 
    and end of the string so a value of `prod` would also match `preprod01`.


`-r` *VALUES*, `--role` *VALUES*

:   Specify the primary role to query. When multiple values are specified they
    must be separated by commas. By default this does a wildcard match at the
    start and end of the string so a value of `MYSQL` would also match
    `MYSQL-DEV`.

`-R` *VALUES*, `--secondary-role` *VALUES*

:   Specify the secondary role to query. When multiple values are specified they
    must be separated by commas. By default this does a wildcard match at the
    start and end of the string so a value of `MASTER` would also match
    `MESOS-MASTER`.

`-i` *VALUES*, `--ip-address` *VALUES*

:   Specify the ip address to query. When multiple values are specified they
    must be separated by commas. This does not do wildcard matching.

`-S` *VALUE*, `--status` *VALUE*

:   Query only for assets that are in a given state. Status takes a argument
    in the form of status:state. An example would be `Allocated:Running`.
    State is however not required and you can just query for `Allocated` as
    well.

`-a` *VALUE*, `--attribute` *VALUE*

:   This is the bread and butter you can query for any key value pair you want
    using this flag. It can be specified multiple times and looks like
    `key:value`.

`-o` *VALUE*, `--operation` *VALUE*

:   Specify operation to join all flags with. By default AND is used. *VALUE*
    can be:

    ::: {#input-formats}
    - `AND`
    - `OR`
    :::

`-s` *VALUE*, `--size` *VALUE*

:   Set the number of results to query for at once. This is by default set to
    be 100. If you are running this from some place where network latency is a
    huge issue this may be useful to increase.

`--limit` *VALUE*

:   Limit the total number of assets that will be returned for the query. By
    default everything will be returned.

Format options {.options}
---------------

`-H`, `--show-header`

:   If you are rendering output in table format show the header above each
    column.

`-c`, *VALUES*, `--columns` *VALUES*

:   Specify a list of all of the attributes that should be shown for table
    formatting. By default we first check if the value is a custom one that we
    know about and need to compute such as total number of CPUs. If it's not a
    special column we check if it's in attributes when outputting it. The
    following special values are available.

    By default tag, hostname, nodeclass, status, pool, primary_role, and
    seconadary_role are returned.

    ::: {#column-values}
    - tag
    - status
    - state
    - classification
    - ip_address
    - cpu_cores
    - cpu_threads
    - cpu_speed_ghz
    - cpu_description
    - gpu_description
    - cpu_product
    - gpu_product
    - cpu_vendor
    - gpu_vendor
    - memory_size_bytes
    - memory_size_total
    - memory_description
    - memory_banks_total
    :::

`-x`, *VALUES*, `--extra-columns` *VALUES*

:   Instead of overwritting the default output columns with `-c` you can
    append to them using `-x`. The same documentaion applies to this.

`-f` *VALUE*, `--field-separator` *VALUE*

:   This sets the delimiter that all columns are seperated by. By default all
    columns are seperated with a tab.

Robot options {.options}
---------------

`-l`, `--link`

:   Output the assets with a link to them in the web UI.

`-j`, `--json`

:   Output the assets as a JSON array.

`-y`, `--yaml`

:   Output the assets as a YAML array.
