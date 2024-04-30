package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/alitto/pond"
	"github.com/urfave/cli/v2"

	"github.com/sixy6e/go-gsf"
)

func convert_gsf(gsf_uri string, config_uri string, in_memory, metadata_only bool) error {
	var (
		out_uri string
		err     error
	)
	log.Println("Processing GSF:", gsf_uri)
	src := gsf.OpenGSF(gsf_uri, config_uri, in_memory)
	defer src.Close()

	log.Println("Building index; Collating metadata; Computing general QA")
	file_info := src.Info()
	proc_info := src.ProcInfo(&file_info)

	log.Println("Writing metadata")
	out_uri = gsf_uri + "-metadata.json"
	_, err = gsf.WriteJson(out_uri, config_uri, file_info.Metadata)
	if err != nil {
		return err
	}

	log.Println("Writing proc-info")
	out_uri = gsf_uri + "-proc-info.json"
	_, err = gsf.WriteJson(out_uri, config_uri, proc_info)
	if err != nil {
		return err
	}

	log.Println("Writing index")
	out_uri = gsf_uri + "-index.json"
	_, err = gsf.WriteJson(out_uri, config_uri, file_info.Index)
	if err != nil {
		return err
	}

	if !metadata_only {
		log.Println("Processing Attitude")
		out_uri = gsf_uri + "-attitude.tiledb"
		att := src.AttitudeRecords(&file_info)
		err = att.ToTileDB(out_uri, config_uri)
		if err != nil {
			return err
		}

		log.Println("Processing SVP")
		out_uri = gsf_uri + "-svp.tiledb"
		svp := src.SoundVelocityProfileRecords(&file_info)
		err = svp.ToTileDB(out_uri, config_uri)
		if err != nil {
			return err
		}

		log.Println("Reading and writing swath bathymetry ping data")
		err = src.SbpToTileDB(&file_info, config_uri)
		if err != nil {
			return err
		}
	}

	log.Println("Finished GSF:", gsf_uri)

	return nil
}

func convert_gsf_list(uri string, config_uri string, in_memory, metadata_only bool) error {
	log.Println("Searching uri:", uri)
	items := gsf.FindGsf(uri, config_uri)
	// out_uris := make([]string, len(items))
	log.Println("Number of GSFs to process:", len(items))

	// TODO; if we write the file to a different structure, we need a different extension
	// for i, name := range(items) {
	//     out_uris[i] = name + "-index.json"
	// }

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
			_ = convert_gsf(item_uri, config_uri, in_memory, metadata_only)
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
					&cli.BoolFlag{
						Name:  "in-memory",
						Usage: "Read the entire contents of a GSF file into memory before processing.",
					},
					&cli.BoolFlag{
						Name:  "metadata-only",
						Usage: "Only decode and export metadata relating to the GSF file.",
					},
					// &cli.StringFlag{
					//     Name: "out-uri",
					//     Usage: "URI or pathname to write the output file to.",
					// },
				},
				Action: func(cCtx *cli.Context) error {
					// err := create_index(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.String("out-uri"))
					err := convert_gsf(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.Bool("in-memory"), cCtx.Bool("metadata-only"))
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
					&cli.BoolFlag{
						Name:  "in-memory",
						Usage: "Read the entire contents of a GSF file into memory before processing.",
					},
					&cli.BoolFlag{
						Name:  "metadata-only",
						Usage: "Only decode and export metadata relating to the GSF files.",
					},
					// &cli.StringFlag{
					//     Name: "out-uri",
					//     Usage: "URI or pathname to write the output file to.",
					// },
				},
				Action: func(cCtx *cli.Context) error {
					// err := create_index(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.String("out-uri"))
					err := convert_gsf_list(cCtx.String("uri"), cCtx.String("config-uri"), cCtx.Bool("in-memory"), cCtx.Bool("metadata-only"))
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
