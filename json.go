package gsf

import (
	"encoding/json"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
)

// WriteJson serialises data to a JSON file. The output location can be locally
// or an object store such as s3.
func WriteJson(file_uri string, config_uri string, data any) (int, error) {

	var config *tiledb.Config
	var err error

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

	// the vfs api auto checks for a file's existence and removes it if we are wanting to write
	stream, err := vfs.Open(file_uri, tiledb.TILEDB_VFS_WRITE)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	jsn, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return 0, err
	}

	bytes_written, err := stream.Write(jsn)

	if err != nil {
		return 0, err
	}

	return bytes_written, nil
}

// JsonDumps constructs a JSON string of the supplied data.
func JsonDumps(data any) (string, error) {
	jsn, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsn), nil
}

// JsonIndentDumps constructs a json string of the supplied data using an
// indentation of four spaces.
func JsonIndentDumps(data any) (string, error) {
	jsn, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsn), nil
}
