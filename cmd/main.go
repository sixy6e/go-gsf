package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
	"github.com/alitto/pond"
	"github.com/urfave/cli/v2"

	"github.com/sixy6e/go-gsf"
)

// convert_gsf handles the conversion process for a single GSF file.
func convert_gsf(gsf_uri, config_uri, outdir_uri string, in_memory, metadata_only, dense bool) error {
	var (
		out_uri string
		err     error
		dir     string
		file    string
		config  *tiledb.Config
	)

	dir, file = filepath.Split(gsf_uri)
	if outdir_uri == "" {
		outdir_uri = dir
	}

	log.Println("Processing GSF:", gsf_uri)
	src := gsf.OpenGSF(gsf_uri, config_uri, in_memory)
	defer src.Close()

	log.Println("Building index; Collating metadata; Computing general QA")
	file_info := src.Info()
	proc_info := src.ProcInfo(&file_info)

	log.Println("Writing metadata")
	out_uri = filepath.Join(outdir_uri, file+"-metadata.json")
	_, err = gsf.WriteJson(out_uri, config_uri, file_info.Metadata)
	if err != nil {
		return err
	}

	log.Println("Writing index")
	out_uri = filepath.Join(outdir_uri, file+"-index.json")
	_, err = gsf.WriteJson(out_uri, config_uri, file_info.Index)
	if err != nil {
		return err
	}

	if !metadata_only {
		// get a generic config if no path provided
		if config_uri == "" {
			config, err = tiledb.NewConfig()
			if err != nil {
				return err
			}
		} else {
			config, err = tiledb.LoadConfig(config_uri)
			if err != nil {
				return err
			}
		}

		defer config.Free()

		ctx, err := tiledb.NewContext(config)
		if err != nil {
			return err
		}
		defer ctx.Free()

		grp_uri := filepath.Join(outdir_uri, file+".tiledb")
		grp, err := tiledb.NewGroup(ctx, grp_uri)
		if err != nil {
			return err
		}
		defer grp.Free()

		err = grp.Create()
		if err != nil {
			return errors.Join(err, errors.New("Error creating tiledb group"))
		}

		err = grp.Open(tiledb.TILEDB_WRITE)
		if err != nil {
			return errors.Join(err, errors.New("Error opening tiledb group in write mode"))
		}

		log.Println("Writing GSF data processing information to group metadata")
		jsn, err := gsf.JsonIndentDumps(proc_info)
		if err != nil {
			return err
		}
		err = grp.PutMetadata("Data-Processing-Information", jsn)
		if err != nil {
			return err
		}

		log.Println("Processing Attitude")
		att_name := "Attitude.tiledb"
		out_uri = filepath.Join(grp_uri, att_name)
		att := src.AttitudeRecords(&file_info)
		err = att.ToTileDB(out_uri, ctx)
		if err != nil {
			return err
		}
		err = grp.AddMember(att_name, "Attitude", true)
		if err != nil {
			return errors.Join(err, errors.New("Error adding attitude to group"))
		}

		log.Println("Processing SVP")
		svp_name := "SVP.tiledb"
		out_uri = filepath.Join(grp_uri, svp_name)
		svp := src.SoundVelocityProfileRecords(&file_info)
		err = svp.ToTileDB(out_uri, ctx)
		if err != nil {
			return err
		}
		err = grp.AddMember(svp_name, "SVP", true)
		if err != nil {
			return errors.Join(err, errors.New("Error adding svp to group"))
		}

		log.Println("Reading and writing swath bathymetry ping data")
		err = src.SbpToTileDB(&file_info, ctx, grp, grp_uri, dense)
		if err != nil {
			return err
		}
	}

	log.Println("Finished GSF:", gsf_uri)

	return nil
}

// convert_gsf_list is responsible for submitting a list of GSF files to a processing pool
// that converts each GSF file. The processing pool uses 2 * n_CPUs workers to spread the
// work across.
func convert_gsf_list(uri, config_uri, outdir_uri string, in_memory, metadata_only, dense bool) error {
	log.Println("Searching uri:", uri)
	items := gsf.FindGsf(uri, config_uri)
	log.Println("Number of GSFs to process:", len(items))

	// Create a context that will be cancelled when the user presses Ctrl+C (process receives termination signal).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// fixed pool
	n := runtime.NumCPU() * 2
	pool := pond.New(n, 0, pond.MinWorkers(n), pond.Context(ctx))
	defer pool.StopAndWait()

	for _, name := range items {
		item_uri := name
		pool.Submit(func() {
			_ = convert_gsf(item_uri, config_uri, outdir_uri, in_memory, metadata_only, dense)
			// if err != nil {
			//     return err
			// }
		})
	}

	return nil // TODO; fix this design
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name: "convert",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "gsf-uri",
						Usage: "URI or pathname to a GSF file.",
					},
					&cli.StringFlag{
						Name:  "config-uri",
						Usage: "URI or pathname to a TileDB config file.",
					},
					&cli.StringFlag{
						Name:  "outdir-uri",
						Usage: "URI or pathname to an output directory.",
					},
					&cli.BoolFlag{
						Name:  "in-memory",
						Usage: "Read the entire contents of a GSF file into memory before processing.",
					},
					&cli.BoolFlag{
						Name:  "metadata-only",
						Usage: "Only decode and export metadata relating to the GSF file.",
					},
					&cli.BoolFlag{
						Name:  "dense",
						Usage: "Create a dense TileDB array schema for the beam data. Default is sparse.",
					},
				},
				Action: func(cCtx *cli.Context) error {
					err := convert_gsf(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.String("outdir-uri"), cCtx.Bool("in-memory"), cCtx.Bool("metadata-only"), cCtx.Bool("dense"))
					return err
				},
			},
			&cli.Command{
				Name: "convert-trawl",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "uri",
						Usage: "URI or pathname to a directory containing gsf files.",
					},
					&cli.StringFlag{
						Name:  "config-uri",
						Usage: "URI or pathname to a TileDB config file.",
					},
					&cli.StringFlag{
						Name:  "outdir-uri",
						Usage: "URI or pathname to an output directory.",
					},
					&cli.BoolFlag{
						Name:  "in-memory",
						Usage: "Read the entire contents of a GSF file into memory before processing.",
					},
					&cli.BoolFlag{
						Name:  "metadata-only",
						Usage: "Only decode and export metadata relating to the GSF files.",
					},
					&cli.BoolFlag{
						Name:  "dense",
						Usage: "Create a dense TileDB array schema for the beam data. Default is sparse.",
					},
				},
				Action: func(cCtx *cli.Context) error {
					err := convert_gsf_list(cCtx.String("uri"), cCtx.String("config-uri"), cCtx.String("outdir-uri"), cCtx.Bool("in-memory"), cCtx.Bool("metadata-only"), cCtx.Bool("dense"))
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
