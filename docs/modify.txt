% Collins Golang CLI
% Michael Schuett
% February 12, 2019

Synopsis
========

`collins modify` [options] ...

Description
===========

Collins modify allows you to change attributes on a given asset. You can pipe
assets to the modify command from `collins query` to easily set a number of assets.

Modify Examples
------------

A basic example to modify a specific tag and set a tag called environment to D123
would look like this.

    collins modify -t M000001 -a environment:D123

Options
=======

Modify options {.options}
---------------

`-t` *VALUES*, `--tags` *VALUES*

:   Specify tag names to return they must be the full tag and in the case that
    multiple are specified they must be separated by commas.

`-a` *VALUE*, `--set-attribute` *VALUE*

:   Set a given key to a value in the format of key:value.

`-d` *VALUE*, `--delete-attribute` *VALUE*

:   Delete attribute takes the name of a key on a given asset. It will return success
    even in the case that the key does not exist on the asset.

`-l` *VALUE*, `--log` *VALUE*

:   *VALUE* contains the message that will be logged for the given asset. You can
    set `--level` to change the level at which the message is logged.


`-L` *VALUE*, `--level` *VALUE*

:   level set the log level when writing logs to collins. By default all logs
    written with `--log` will show up as a note. The following values are valid
    log levels.

    ::: {#input-formats}
    - `ERROR`
    - `DEBUG`
    - `EMERGENCY`
    - `ALERT`
    - `CRITICAL`
    - `WARNING`
    - `NOTICE`
    - `INFORMATIONAL`
    - `NOTE`
    :::

`-S` *VALUE*, `--set-state` *VALUE*

:   Set the status and optionall state that should be set for a given or multiple
    tags. The format is status:state and an example of it would look like 
    Running:allocated. When changing the state of an asset a `--reason` must be
    given.

`-r` *VALUE*, `--reason` *VALUE*

:   Set the reason for changing the state of an asset. The reason may be any not
    empty string although others will appreciate if you set it to something meaningful.
