% storage.conf(5) Container Storage Configuration File
% Dan Walsh
% May 2017

# NAME
storage.conf - Syntax of Container Storage configuration file

# DESCRIPTION
The STORAGE configuration file specifies all of the available container storage options
for tools using shared container storage, but in a TOML format that can be more easily modified
and versioned.

# FORMAT
The [TOML format][toml] is used as the encoding of the configuration file.
Every option and subtable listed here is nested under a global "storage" table.
No bare options are used. The format of TOML can be simplified to:

    [table]
    option = value

    [table.subtable1]
    option = value

    [table.subtable2]
    option = value

## STORAGE TABLE

The `storage` table supports the following options:

**graphroot**=""
  container storage graph dir (default: "/var/lib/containers/storage")
  Default directory to store all writable content created by container storage programs

**runroot**=""
  container storage run dir (default: "/var/run/containers/storage")
  Default directory to store all temporary writable content created by container storage programs

**driver**=""
  container storage driver (default is "overlay")
  Default Copy On Write (COW) container storage driver

### STORAGE OPTIONS TABLE 

The `storage.options` table supports the following options:

**additionalimagestores**=[]
  Paths to additional container image stores. Usually these are read/only and stored on remote network shares.

**size**=""
  Maximum size of a container image.  Default is 10GB.  This flag can be used to set quota
  on the size of container images.

**override_kernel_check**=""
  Tell storage drivers to ignore kernel version checks.  Some storage drivers assume that if a kernel is too
  old, the driver is not supported.  But for kernels that have had the drivers backported, this flag
  allows users to override the checks

# HISTORY
May 2017, Originally compiled by Dan Walsh <dwalsh@redhat.com>
Format copied from crio.conf man page created by Aleksa Sarai <asarai@suse.de>
