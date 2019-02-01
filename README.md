Collins Golang CLI
==================

This is a clone of the original collins-cli but in golang and currently with less features.
I am aiming to have be a drop in replacement for collins-cli once it is finished. This was
mainly developed because managing ruby became a pain with the many different projects we
had that used ruby and all had it configured in different way (rbenv, rvm, system ruby).
Some projects changed ruby versions when you CD'ed into them. Some had bundle environments
that disabled your system gems. Additionally I also just wanted to write something in go.

## Feature List

|Feature                               |Completed|
|--------------------------------------|:-------:|
|Query subcommand feature parity       |90%      |
|Power subcommand feature parity       |90%      |
|Datacenter subcommand feature parity  |         |
|Modify subcommand feature parity      |100%     |
|Log subcommand feature parity         |         |
|Provision subcommand feature parity   |90%      |
|IPAM subcommand feature parity        |         |

## Notes

* The query subcommand is missing some output formatting such as YAML, JSON, and Links however
the ability to query assets is the same as the ruby version.
* Both power and provision are set at 90% because the way we log to stdout could be better. For instance
we don't output anything until the request completes. This should be switched to output when the request
starts and then when we get error/success. However since I have been logging with logrus I will need to
change that out in order to support this and likely just go to logging with fmt.
