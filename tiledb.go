package gsf

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
	"github.com/samber/lo"
	stgpsr "github.com/yuin/stagparser"
)

var ErrAddFilters = errors.New("Error Adding Filter To FilterList")
var ErrDims = errors.New("Error Dims is > 2")                   // we should not have any slices > 2D
var ErrDtype = errors.New("Error slice datatype is unexpected") // we should not have any slices > 2D
var ErrSetBuff = errors.New("Error setting tiledb buffer")      // we should not have any slices > 2D

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
	case "string":
		tdb_dtype = tiledb.TILEDB_STRING_UTF8
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

// sliceOffsets is a helper func to calculate the 1D array offsets for fields
// that are of variable length.
func sliceOffsets[T any](s [][]T, byte_size uint64) (slc_offset []uint64) {
	nrows := len(s)
	slc_offset = make([]uint64, nrows)
	offset := uint64(0)

	for i := 0; i < nrows; i++ {
		length := uint64(len(s[i]))
		slc_offset[i] = offset
		offset += length * byte_size
	}

	return slc_offset
}

func setStructFieldBuffers(query *tiledb.Query, t any) error {
	var (
		err error
	)

	bytesize1 := uint64(1)
	bytesize2 := uint64(2)
	bytesize4 := uint64(4)
	bytesize8 := uint64(8)

	values := reflect.ValueOf(t).Elem()
	types := reflect.TypeOf(t).Elem()
	for i := 0; i < values.NumField(); i++ {
		fld := values.Field(i)
		typ := fld.Type()

		if types.Field(i).IsExported() {
			name := types.Field(i).Name
			dims := 0
			stype := sliceDimsType(typ, &dims)

			switch dims {
			case 1:
				switch stype.Name() {
				case "int8":
					slc := fld.Interface().([]int8)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "uint8":
					slc := fld.Interface().([]uint8)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "int16":
					slc := fld.Interface().([]int16)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "uint16":
					slc := fld.Interface().([]uint16)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "int32":
					slc := fld.Interface().([]int32)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "uint32":
					slc := fld.Interface().([]uint32)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "int64":
					slc := fld.Interface().([]int64)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "uint64":
					slc := fld.Interface().([]uint64)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "float32":
					slc := fld.Interface().([]float32)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "float64":
					slc := fld.Interface().([]float64)
					_, err = query.SetDataBuffer(name, slc)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				case "Time":
					slc := fld.Interface().([]time.Time)

					// time arrays need an additional conversion for serialisation
					nrows := len(slc)
					timestamps := make([]int64, nrows)
					for t := 0; t < nrows; t++ {
						timestamps[t] = slc[t].UnixNano()
					}

					_, err = query.SetDataBuffer(name, timestamps)
					if err != nil {
						return errors.Join(ErrSetBuff, err, errors.New(name))
					}
				default:
					// some datatype we haven't accounted for
					return errors.Join(ErrDtype, errors.New(stype.Name()))
				}
			case 2:
				// these will be the variable length arrays
				// this approach won't work for say the BrbIntensity.TimeSeries
				// which is stored as a single 1D slice, and the count stored elsewhere
				// on the struct (unless we change it)
				// For var length arrays, the procedure is to create a flattened version
				// of the 2D slice, calculate byte offsets, and set the buffers for
				// both the flattened and byte offset slices
				switch stype.Name() {
				case "int8":
					slc := fld.Interface().([][]int8)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize1)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "uint8":
					slc := fld.Interface().([][]uint8)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize1)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "int16":
					slc := fld.Interface().([][]int16)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize2)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "uint16":
					slc := fld.Interface().([][]uint16)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize2)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "int32":
					slc := fld.Interface().([][]int32)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize4)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "uint32":
					slc := fld.Interface().([][]uint32)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize4)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "int64":
					slc := fld.Interface().([][]int64)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize8)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "uint64":
					slc := fld.Interface().([][]uint64)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize8)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "float32":
					slc := fld.Interface().([][]float32)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize4)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "float64":
					slc := fld.Interface().([][]float64)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize8)

					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, flt)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				case "Time":
					slc := fld.Interface().([][]time.Time)
					flt := lo.Flatten(slc)
					slc_offset := sliceOffsets(slc, bytesize8)

					// time arrays need an additional conversion for serialisation
					nrows := len(flt)
					timestamps := make([]int64, nrows)
					for t := 0; t < nrows; t++ {
						timestamps[t] = flt[t].UnixNano()
					}
					_, err = query.SetOffsetsBuffer(name, slc_offset)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}

					_, err = query.SetDataBuffer(name, timestamps)
					if err != nil {
						return errors.Join(err, errors.New(name))
					}
				default:
					// some datatype we haven't accounted for
					return errors.Join(ErrDtype, errors.New(stype.Name()))
				}
			default:
				return errors.Join(ErrDims, errors.New(strconv.Itoa(dims)))
			}
		}
	}
	return nil
}

// WriteArrayMetadata is a helper for attaching/writing metadata to a TileDB array.
// The metadata is converted to JSON before writing to TileDB.
func WriteArrayMetadata(ctx *tiledb.Context, array_uri, key string, md any) error {
	array, err := ArrayOpen(ctx, array_uri, tiledb.TILEDB_WRITE)
	if err != nil {
		return errors.Join(err, errors.New("Error opening (w) TileDB array: "+array_uri))
	}
	defer array.Free()
	defer array.Close()

	jsn, err := JsonDumps(md)
	if err != nil {
		return errors.Join(err, errors.New("Error serialising metadata to JSON"))
	}

	err = array.PutMetadata(key, jsn)
	if err != nil {
		return errors.Join(err, errors.New("Error writing metadata to array: "+array_uri))
	}

	return nil
}
