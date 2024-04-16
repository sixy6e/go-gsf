package gsf

import (
	"errors"
	"reflect"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
	stgpsr "github.com/yuin/stagparser"
)

var ErrAddFilters = errors.New("Error Adding Filter To FilterList")

// ArrayOpen is a helper func for opening a tiledb array.
func ArrayOpen(ctx *tiledb.Context, uri string, mode tiledb.QueryType) (*tiledb.Array, error) {
	array, err := tiledb.NewArray(ctx, uri)
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
// compression level. Note; the compression level is meaningless for RLE, and
// is quietly ignored internally by TileDB.
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

// BitWidthReductionFilter initialises the Bit width reduction and sets the
// window size.
func BitWidthReductionFilter(ctx *tiledb.Context, window int32) (*tiledb.Filter, error) {
	filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_BIT_WIDTH_REDUCTION)
	if err != nil {
		return nil, err
	}

	err = filt.SetOption(tiledb.TILEDB_BIT_WIDTH_MAX_WINDOW, window)
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

// CreateAttr creates a tiledb attribute along with the compression filter
// pipeline. The configuration is specified by the tags attached to the
// struct type.
// Tags for tiledb include: dtype, var, ftype.
// Where dtype is datatype, var is variable length, ftype is fieldtype
// (dim or attr) for dimension or attribute (dim skips the field).
// Supported datatype values are int8, uint8, int16, uint16, int32, uint32,
// int64, uint64, float32, float64, datetime_ns.
// Tags for filters include: zstd(level=16), gzip(level=6), bysh, bish,
// lz4(level=6), rle(level=-1), bzip2(level=6), bitw(window=-1).
// Where level indicates the compression level, window indicates the window size
// (-1 indicates default), zstd is zstandard, gzip is deflate,
// rle is run length encoding, bysh is byteshuffle, bish is bitshuffle and
// bitw is bit width reduction.
// Filters will be set in the order they're specified in the tag.
// Variable length fields will have the offsets compressed using a default
// strategy of positive-delta, byteshuffle, and finally zstandard with level=16.
// An example tag is `tiledb:"dtype=uint16,ftype=attr" filters:"bysh,zstandard(level=16)"`
func CreateAttr(
	field_name string,
	filter_defs []stgpsr.Definition,
	tiledb_defs map[string]stgpsr.Definition,
	schema *tiledb.ArraySchema,
	ctx *tiledb.Context,
) error {

	var (
		tdb_dtype tiledb.Datatype
		def       stgpsr.Definition
		status    bool
	)

	def, status = tiledb_defs["dtype"]
	if !status {
		return errors.Join(ErrCreateSvpTdb, errors.New("dtype tag not found"))
	}
	dtype, _ := def.Attribute("dtype")

	// define datatype
	switch dtype {
	case "int8":
		tdb_dtype = tiledb.TILEDB_INT8
	case "uint8":
		tdb_dtype = tiledb.TILEDB_UINT8
	case "int16":
		tdb_dtype = tiledb.TILEDB_INT16
	case "uint16":
		tdb_dtype = tiledb.TILEDB_UINT16
	case "int32":
		tdb_dtype = tiledb.TILEDB_INT32
	case "uint32":
		tdb_dtype = tiledb.TILEDB_UINT32
	case "int64":
		tdb_dtype = tiledb.TILEDB_INT64
	case "uint64":
		tdb_dtype = tiledb.TILEDB_UINT64
	case "float32":
		tdb_dtype = tiledb.TILEDB_FLOAT32
	case "float64":
		tdb_dtype = tiledb.TILEDB_FLOAT64
	case "datetime_ns": // can add other datetime types when required
		tdb_dtype = tiledb.TILEDB_DATETIME_NS
	}

	attr_filts, err := tiledb.NewFilterList(ctx)
	if err != nil {
		return errors.Join(ErrCreateSvpTdb, err)
	}
	defer attr_filts.Free()

	// filter pipeline
	for _, filter := range filter_defs {
		switch filter.Name() {
		case "zstd":
			level, status := filter.Attribute("level")
			if !status {
				return errors.Join(ErrCreateSvpTdb, errors.New("zstd level not defined"))
			}
			filt, err := ZstdFilter(ctx, int32(level.(int64)))
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		case "gzip":
			level, status := filter.Attribute("level")
			if !status {
				return errors.Join(ErrCreateSvpTdb, errors.New("gzip level not defined"))
			}
			filt, err := ZstdFilter(ctx, int32(level.(int64)))
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		case "lz4":
			level, status := filter.Attribute("level")
			if !status {
				return errors.Join(ErrCreateSvpTdb, errors.New("lz4 level not defined"))
			}
			filt, err := Lz4Filter(ctx, int32(level.(int64)))
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		case "rle":
			level, status := filter.Attribute("level")
			if !status {
				return errors.Join(ErrCreateSvpTdb, errors.New("rle level not defined"))
			}
			filt, err := RleFilter(ctx, int32(level.(int64)))
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		case "bzip2":
			level, status := filter.Attribute("level")
			if !status {
				return errors.Join(ErrCreateSvpTdb, errors.New("bzip2 level not defined"))
			}
			filt, err := Bzip2Filter(ctx, int32(level.(int64)))
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		case "bitw":
			win, status := filter.Attribute("window")
			if !status {
				return errors.Join(ErrCreateSvpTdb, errors.New("bitwidth window not defined"))
			}
			filt, err := BitWidthReductionFilter(ctx, int32(win.(int64)))
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		case "bish":
			filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_BITSHUFFLE)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		case "bysh":
			filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_BYTESHUFFLE)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
			defer filt.Free()
			err = attr_filts.AddFilter(filt)
			if err != nil {
				return errors.Join(ErrCreateSvpTdb, err)
			}
		}
	}

	// create attr
	attr, err := tiledb.NewAttribute(ctx, field_name, tdb_dtype)
	if err != nil {
		return errors.Join(ErrCreateSvpTdb, err)
	}
	defer attr.Free()

	// variable length attrs
	_, status = tiledb_defs["var"]
	if status {
		attr.SetCellValNum(tiledb.TILEDB_VAR_NUM)
		if err != nil {
			return errors.Join(ErrCreateSvpTdb, err)
		}
	}

	// attach filter pipeline to attr
	err = AttachFilters(attr_filts, attr)
	if err != nil {
		return errors.Join(ErrCreateSvpTdb, err)
	}

	// attach attr to schema
	err = schema.AddAttributes(attr)
	if err != nil {
		return errors.Join(ErrCreateSvpTdb, err)
	}

	// variable length attrs filters
	// making an assumption that the var attr needs to be set on the schema
	// before we add the offsets filter pipeline to the schema
	if status {
		offset_filts, err := tiledb.NewFilterList(ctx)
		if err != nil {
			return errors.Join(ErrCreateSvpTdb, err)
		}

		dd_filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_POSITIVE_DELTA)
		if err != nil {
			return errors.Join(ErrCreateSvpTdb, err)
		}

		bysh_filt, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_BYTESHUFFLE)
		if err != nil {
			return errors.Join(ErrCreateSvpTdb, err)
		}

		zstd_filt, err := ZstdFilter(ctx, int32(16))
		if err != nil {
			return errors.Join(ErrCreateSvpTdb, err)
		}

		err = AddFilters(offset_filts, dd_filt, bysh_filt, zstd_filt)
		if err != nil {
			return errors.Join(ErrCreateSvpTdb, err)
		}

		err = schema.SetOffsetsFilterList(offset_filts)
		if err != nil {
			return errors.Join(ErrCreateSvpTdb, err)
		}
	}

	return nil
}

// sliceDimsType is a helper for determining the numver of dimensions
// and the underlying type a slice contains.
// This func is called elsewhere that is undertaking reflection on
// a struct whose fields are slices.
// Care needs to be taken in that the original caller must initialise the
// int that the dims pointer points to, is zero.
// The primary motivation was not to be explicitly calling each structs
// field, for example EM4 which contains 53 fields, and would be a lot of code.
// Multiply that for over a dozen different sensors, and that's a lot of code.
// However, it would be more explicit, and easier to follow. I have found that
// reflection is hard to follow, and I could have easily introduced more errors
// through blind assumptions, than being explicit and calling each field by name
// for serialisation.
func sliceDimsType(typ reflect.Type, dims *int) reflect.Type {
	if typ.Kind() == reflect.Slice {
		*dims += 1
		return sliceDimsType(typ.Elem(), dims)
	}

	// either not a slice, or we've buried deep enough to the underliying
	// slice type; eg uint8, float32, time.Time etc
	return typ
}
