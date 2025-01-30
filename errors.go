package gsf

import (
	"errors"
)

var ErrCreateAttitudeTdb = errors.New("Error Creating Attitude TileDB Array")
var ErrWriteAttitudeTdb = errors.New("Error Writing Attitude TileDB Array")
var ErrCreateBdTdb = errors.New("Error Creating Beam Data TileDB Array")
var ErrWriteBdTdb = errors.New("Error Writing Beam Data TileDB Array")
var ErrCreateMdTdb = errors.New("Error Creating Metadata TileDB Array")
var ErrWriteMdTdb = errors.New("Error Writing Metadata TileDB Array")
var ErrCreateAttributeTdb = errors.New("Error Creating Attribute for TileDB Array")
var ErrCreateMdDenseTdb = errors.New("Error Creating Dense Metadata TileDB Array")
var ErrCreateBeamSparseTdb = errors.New("Error Creating Beam Sparse TileDB Array")
var ErrCreateSchemaTdb = errors.New("Error Creating TileDB Schema")
var ErrCreateDimTdb = errors.New("Error Creating TileDB Dimension")
var ErrSensor = errors.New("Sensor Not Supported")
var ErrWriteSensorMd = errors.New("Error Writing Sensor Metadata")
var ErrSensorImgMetadata = errors.New("Error Reading Sensor Imagery Metadata")
var ErrSensorMetadata = errors.New("Error Reading Sensor Metadata")
var ErrCreateSvpTdb = errors.New("Error Creating SVP TileDB Array")
var ErrWriteSvpTdb = errors.New("Error Writing SVP TileDB Array")
var ErrAddFilters = errors.New("Error Adding Filter To FilterList")
var ErrDims = errors.New("Error Dims Is > 2")
var ErrDtype = errors.New("Error Slice Datatype Is Unexpected") // we should not have any slices > 2D
var ErrSetBuff = errors.New("Error Setting TileDB Buffer")
var ErrFiltList = errors.New("Error Creating TileDB Filter List")
var ErrNewAttr = errors.New("Error Creating TileDB Attribute")
var ErrNewFilt = errors.New("Error Creating TileDB Filter")
var ErrSetFiltList = errors.New("Error Setting TileDB Filter List")
var ErrAddAttr = errors.New("Error Adding TileDB Attribute")
var ErrZstdFilt = errors.New("Error Creating TileDB ZStandard Filter")
