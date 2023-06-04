package search

import (
    "path/filepath"
    tiledb "github.com/TileDB-Inc/TileDB-Go"
)

// An internal general purpose trawling function. Potentially could be globally
// exported at a later date.
// The basename is only matched with the pattern, eg
// ("*.gsf", "0060_20150624_185509_Investigator_em710.gsf")
func trawl(vfs *tiledb.VFS, pattern string, uri string, items []string) []string {
    dirs, files, err := vfs.List(uri)
    if err != nil {
        panic(err)
    }

    // check files for the matching pattern
    for _, file := range(files) {
        match, err := filepath.Match(pattern, filepath.Base(file))
        if err != nil {
            panic(err)
        }

        if match {
            items = append(items, file)
        }
    }

    // recurse over every directory
    for _, dir := range(dirs) {
        items = trawl(vfs, pattern, dir, items)
    }

    return items
}

// A specific function to recursively search for *.gsf files under a given URI.
// The function uses the TileDB Go bindings to seamlessly search either local
// filesystems or obeject stores such as AWS-S3. A TileDB config is required
// for searching object stores with permission constraints.
func FindGsf(uri string, config_uri string) []string {
    var (
        config *tiledb.Config
        err error
        items []string
        pattern string
    )

    // get a generic config if no path provided
    if config_uri == "" {
        config, err = tiledb.NewConfig()
        if err != nil {
            panic(err)
        }
    } else {
        config, err = tiledb.LoadConfig(config_uri)
        if err != nil {
            panic(err)
        }
    }

    defer config.Free()

    ctx, err := tiledb.NewContext(config)
    if err != nil {
        panic(err)
    }
    defer ctx.Free()

    vfs, err := tiledb.NewVFS(ctx, config)
    if err != nil {
        panic(err)
    }
    defer vfs.Free()

    items = make([]string, 0)
    pattern = "*.gsf"

    items = trawl(vfs, pattern, uri, items)

    return items
}
