package main

import (
	"context"
	"fmt"
	"github.com/gnasnik/titan-sdk-go"
	"github.com/gnasnik/titan-sdk-go/config"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "titan cli",
		Usage: "titan's toolset",
		Commands: []*cli.Command{
			downloadFileCmd,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var downloadFileCmd = &cli.Command{
	Name:    "download",
	Aliases: []string{"d"},
	Usage:   "get file from titan network",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "cid",
			Aliases: []string{"c"},
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.String("cid") == "" {
			return fmt.Errorf("cid is required")
		}
		return getFile(cctx.String("cid"), cctx.String("output"))
	},
}

func getFile(cid string, output string) error {
	address := os.Getenv("LOCATOR_API_INFO")
	client, err := titan.New(
		config.AddressOption(address),
	)
	if err != nil {
		return err
	}

	_, reader, err := client.GetFile(context.Background(), cid)
	if err != nil {
		return err
	}
	defer reader.Close()

	if output == "" {
		io.Copy(io.Discard, reader)
		return nil
	}

	file, err := os.Create(output)
	if err != nil {
		return err
	}

	io.Copy(file, reader)
	return nil
}
