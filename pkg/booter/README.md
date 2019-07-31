# Booter package

The booter package provides a Booter interface that allows to define custom ways
of booting a machine. Booter only requires a method to get its name, and a
method to boot the machine. This interface is suitable to be used by
Systemboot's `uinit`.

Each custom booter will define their own name (e.g. "netboot") and the custom
logic to boot. For example, a network booter may try to get a network
configuration via DHCPv6, download a boot program, and run it.

The custom booter also needs to provide a way to be initialized. This is usually
done by defining a function called "New<MyBooterName>" (e.g. "NewNetBooter").
This function takes as input a sequence of bytes, representing the booter
configuration, and will return a Booter object, or an error.

The exact format of the boot configuration is determined by the custom booter,
but there is a general structure that every booter configuration has to provide.
This is discussed in the Booter configuration section below

## Booter configuration

A Booter configuration is a JSON file with a simple structure. The requirements
are:

* the top level object is a map
* the map contains at least a "type" field, with a string value that holds the
  booter's name
* the JSON should not be nested. This is recommended for simplicity but is not
  strictly required

For example, the NetBooter configuration can be like the following:

```
{
    "type": "netboot",
    "method": "<method>",
    "mac": "<mac address",
    "override_url": "<url>"
}
```

where:

* "type" is required, and its value is always "netboot" (otherwise it's not
  recognized as a NetBooter)
* "method" is required, and can be either "dhcpv6", "dhcpv4" or "slaac"
* "mac" is required, and it is the MAC address of the interface that will try
  to boot from the network. It has the "aa:bb:cc:dd:ee:ff" format
* "override_url" is optional, unles "method" is "slaac", and it is the URL from
  which the booter will try to download the network boot program


## Creating a new Booter

To create a new Booter, the following things are necessary:

* define a structure for the new booter, that implements the Booter interface
  described above. I.e. implement the `TypeName` and `Boot` methods
* define a NewMyBooterName (e.g. "NewLocalBoot") that takes a sequence of bytes
  as input, and return a `Booter` or an error if it's an invalid or unknown
  configuration. The input byte sequence must contain a valid JSON configuration
  for that booter in order to return successfully
* the new booter has to be registered as a supported booter. Just add the
  `NewMyBooterName` function (however this function is called) to the
  `supportedBooterParsers` in `bootentry.go`. This array is used by
  `GetBootEntries` to test a boot configuration against all the available
  booters

