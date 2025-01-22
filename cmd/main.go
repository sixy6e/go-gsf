package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/alitto/pond"
	"github.com/urfave/cli/v2"

	"github.com/sixy6e/go-gsf"
)

func convert_gsf(gsf_uri, config_uri, outdir_uri string, in_memory, metadata_only, dense bool) error {
	var (
		out_uri string
		err     error
		dir     string
		file    string
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

	log.Println("Writing proc-info")
	out_uri = filepath.Join(outdir_uri, file+"-proc-info.json")
	_, err = gsf.WriteJson(out_uri, config_uri, proc_info)
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
		log.Println("Processing Attitude")
		out_uri = filepath.Join(outdir_uri, file+"-attitude.tiledb")
		att := src.AttitudeRecords(&file_info)
		err = att.ToTileDB(out_uri, config_uri)
		if err != nil {
			return err
		}

		log.Println("Processing SVP")
		out_uri = filepath.Join(outdir_uri, file+"-svp.tiledb")
		svp := src.SoundVelocityProfileRecords(&file_info)
		err = svp.ToTileDB(out_uri, config_uri)
		if err != nil {
			return err
		}

		log.Println("Reading and writing swath bathymetry ping data")
		err = src.SbpToTileDB(&file_info, config_uri, outdir_uri, dense)
		if err != nil {
			return err
		}
	}

	log.Println("Finished GSF:", gsf_uri)

	return nil
}

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
