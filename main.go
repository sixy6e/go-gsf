package main

import (
    "log"
    "encoding/json"
    "os"
    //"fmt"
    "runtime"
    "os/signal"
    "context"

    "github.com/urfave/cli/v2"
    "github.com/alitto/pond"

    "gsf/decode"
    "gsf/encode"
    "gsf/search"
)

// func create_index(gsf_uri string, config_uri string, out_uri string) error {
func create_index(gsf_uri string, config_uri string, in_memory bool) error {
    log.Println("Processing GSF:", gsf_uri)
    src := decode.OpenGSF(gsf_uri, config_uri, in_memory)
    defer src.Close()
    file_index := src.Index()

    jsn, err := json.MarshalIndent(file_index, "", "    ")
    if err != nil {
        // panic(err)
        return err
    }

    // TODO; if we write the file to a different structure, we need a different extension
    out_uri := gsf_uri + "-index.json"
    _, err = encode.WriteJson(out_uri, config_uri, jsn)
    if err != nil {
        // panic(err)
        return err
    }

    log.Println("Finished GSF:", gsf_uri)

    return nil
}

func create_index_list(uri string, config_uri string, in_memory bool) error {
    log.Println("Searching uri:", uri)
    items := search.FindGsf(uri, config_uri)
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
    pool := pond.New(n, 0,  pond.MinWorkers(n), pond.Context(ctx))
    defer pool.StopAndWait()

    for _, name := range(items) {
        item_uri := name
        pool.Submit(func() {
            _ = create_index(item_uri, config_uri, in_memory)
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
                Name: "index",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name: "gsf-uri",
                        Usage: "URI or pathname to a GSF file.",
                    },
                    &cli.StringFlag{
                        Name: "config-uri",
                        Usage: "URI or pathname to a TileDB config file.",
                    },
                    &cli.BoolFlag{
                        Name: "in-memory",
                        Usage: "Read the entire contents of a GSF file into memory before processing.",
                    },
                    // &cli.StringFlag{
                    //     Name: "out-uri",
                    //     Usage: "URI or pathname to write the output file to.",
                    // },
                },
                Action: func(cCtx *cli.Context) error {
                    // err := create_index(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.String("out-uri"))
                    err := create_index(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.Bool("in-memory"))
                    return err
                },
            },
            &cli.Command{
                Name: "index-trawl",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name: "uri",
                        Usage: "URI or pathname to a directory containing gsf files.",
                    },
                    &cli.StringFlag{
                        Name: "config-uri",
                        Usage: "URI or pathname to a TileDB config file.",
                    },
                    &cli.BoolFlag{
                        Name: "in-memory",
                        Usage: "Read the entire contents of a GSF file into memory before processing.",
                    },
                    // &cli.StringFlag{
                    //     Name: "out-uri",
                    //     Usage: "URI or pathname to write the output file to.",
                    // },
                },
                Action: func(cCtx *cli.Context) error {
                    // err := create_index(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.String("out-uri"))
                    err := create_index_list(cCtx.String("uri"), cCtx.String("config-uri"), cCtx.Bool("in-memory"))
                    return err
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
