# go-gsf
Prototype GSF decoder in Go.

A mini side project to learn Go at the same time as developing an alternative GSF decoder than the original C version.
As well as a more efficient GSF decoder than the Python ones that are out there.

The main focal point of deconstructing the GSF is to encode various data structures into something more mainstream, such as TileDB.

This enables the data IO to be shifted to a dedicated and generic API, enabling more tools to tap into reading and writing the data, thus enabling more users and generic users that come from outside the marin community.

Most of the internal guts of the GSF data (stored as records) can be mapped into saner data constructs, such as:

* JSON (or YAML) for metadata
  * Processing information
  * GSF file index information
  * Quality information
  * Global schema information
  * General descriptive information
* Dense (or Sparse) TileDB arrays for:
  * Beam data
  * Backscatter data
  * Sound Velocity Profile data
  * Attitude data


# Command line utilities

## Metadata

There is a command line utility for generating metadata based on contents of a GSF file. The metadata utility works two modes:

* Individual file
* Trawler; trawl/crawl a directory tree for GSF files to process


### Individual file

```Shell
$ ./gsf metadata --help
NAME:
   gsf metadata

USAGE:
   gsf metadata [command options] [arguments...]

OPTIONS:
   --gsf-uri value     URI or pathname to a GSF file.
   --config-uri value  URI or pathname to a TileDB config file.
   --in-memory         Read the entire contents of a GSF file into memory before processing. (default: false)
   --help, -h          show help
```

### Trawler

```Shell
$ ./gsf metadata-trawl --help
NAME:
   gsf metadata-trawl

USAGE:
   gsf metadata-trawl [command options] [arguments...]

OPTIONS:
   --uri value         URI or pathname to a directory containing gsf files.
   --config-uri value  URI or pathname to a TileDB config file.
   --in-memory         Read the entire contents of a GSF file into memory before processing. (default: false)
   --help, -h          show help
```
