# config

config package implements an interface to parse toml config files for
various modules. It supports configuration from files, flags and default
values specified in source code.

The package user only has to specify the config struct in source code.

Lot of the heavy lifting is done by the library. The library handles declaring
the flags based on the json tags of the struct. The test example shows a
simple wrapper that can be called by user to check for flags that 
override the defaults in the source code and the values specified in a
config file. This wrapper can be used in similar form or modified order to
support various priority order for different ways of providing the config.
In the provided test the order of priority is flags override file input and
file input overrides defaults in source code.

If not anything the library is a good source of reference for golang
flags package and reflection usage :)

TODO: Integrate environment variables into the scheme for completeness.

