package gsf

import (
    "errors"

    tiledb "github.com/TileDB-Inc/TileDB-Go"
)

var ErrAddFilters = errors.New("Error Adding Filter To FilterList")

// ArrayOpen is a helper func for opening a tiledb array.
func ArrayOpen(ctx *tiledb.Context, uri string, mode tiledb.QueryType) (*tiledb.Array, error) {
    array, err := tiledb.NewArray(ctx, file_uri)
    if err != nil {
        return nil, err
    }

    err = array.Open(mode)
    if err != nil {
        array.Free()
        return nil, err
    }

    return array, nil
}

// AddFilters sequentially appends compression filters to the filter pipeline list.
func AddFilters(filter_list *tiledb.FilterList, filter ...*tiledb.Filter) error {
    for _, filt := range filter {
        err := filter_list.AddFilter(filt)
        if err != nil {
            return err
        }
    }

    return nil
}

// ZstdFilter initialises the Zstandard compression filter and sets the compression
// level.
func ZstdFilter(ctx *tiledb.Context, level int32) (*tiledb.Filter, error) {
    filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_ZSTD)
    if err != nil {
        return nil, err
    }

    err = filt.SetOption(tiledb.TILEDB_COMPRESSION_LEVEL, level)
    if err != nil {
        filt.Free()
        return nil, err
    }

    return filt, nil
}

// GzipFilter initialises the deflate compression filter and sets the compression
// level.
func GzipFilter(ctx *tiledb.Context, level int32) (*tiledb.Filter, error) {
    filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_GZIP)
    if err != nil {
        return nil, err
    }

    err = filt.SetOption(tiledb.TILEDB_COMPRESSION_LEVEL, level)
    if err != nil {
        filt.Free()
        return nil, err
    }

    return filt, nil
}

// Lz4Filter initialises the LZ4 compression filter and sets the compression
// level.
func Lz4Filter(ctx *tiledb.Context, level int32) (*tiledb.Filter, error) {
    filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_LZ4)
    if err != nil {
        return nil, err
    }

    err = filt.SetOption(tiledb.TILEDB_COMPRESSION_LEVEL, level)
    if err != nil {
        filt.Free()
        return nil, err
    }

    return filt, nil
}

// RleFilter initialises the Run Length Encoding compression filter and sets the
// compression level.
func RleFilter(ctx *tiledb.Context, level int32) (*tiledb.Filter, error) {
    filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_RLE)
    if err != nil {
        return nil, err
    }

    err = filt.SetOption(tiledb.TILEDB_COMPRESSION_LEVEL, level)
    if err != nil {
        filt.Free()
        return nil, err
    }

    return filt, nil
}

// Bzip2Filter initialises the Burrows-Wheeler compression filter and sets the
// compression level.
func Bzip2Filter(ctx *tiledb.Context, level int32) (*tiledb.Filter, error) {
    filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_BZIP2)
    if err != nil {
        return nil, err
    }

    err = filt.SetOption(tiledb.TILEDB_COMPRESSION_LEVEL, level)
    if err != nil {
        filt.Free()
        return nil, err
    }

    return filt, nil
}

// AttachFilters acts as a helper for when setting the same pipeline filter list to
// a bunch of attributes.
func AttachFilters(filter_list *tiledb.FilterList, attrs ...*tiledb.Attribute) error {
    for _, attr := range attrs {
         err := attr.SetFilterList(filter_list)
         if err != nil {
             return err
         }
    }

    return nil
}
