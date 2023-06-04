package main

import (
    "log"
    "encoding/json"
    "os"
    "fmt"

    "github.com/urfave/cli/v2"

    "gsf/decode"
    "gsf/encode"
)

// func create_index(gsf_uri string, config_uri string, out_uri string) error {
func create_index(gsf_uri string, config_uri string) error {
    file_index := decode.Index(gsf_uri, config_uri)

    jsn, err := json.MarshalIndent(file_index, "", "    ")
    if err != nil {
        // panic(err)
        return err
    }

    // TODO; if we write the file to a different structure, we need a different extension
    out_uri := uri + "-index.json"
    _, err = encode.WriteJson(out_uri, config_uri, jsn)
    if err != nil {
        // panic(err)
        return err
    }

    return nil
}

func create_index_list(uri, string, config_uri string) error {
    items := FindGsf(uri, config_uri)
    // out_uris := make([]string, len(items))

    // TODO; if we write the file to a different structure, we need a different extension
    // for i, name := range(items) {
    //     out_uris[i] = name + "-index.json"
    // }


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
                    // &cli.StringFlag{
                    //     Name: "out-uri",
                    //     Usage: "URI or pathname to write the output file to.",
                    // },
                },
                Action: func(cCtx *cli.Context) error {
                    // err := create_index(cCtx.String("gsf-uri"), cCtx.String("config-uri"), cCtx.String("out-uri"))
                    err := create_index(cCtx.String("gsf-uri"), cCtx.String("config-uri"))
                    return err
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
