package main

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb"
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
		&cli.BoolFlag{
			Name:    "range",
			Aliases: []string{"r"},
		},
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

		cid := cctx.String("cid")
		output := cctx.String("output")
		isRange := cctx.Bool("range")

		return getFile(cid, output, isRange)
	},
}

func getFile(cid string, output string, isRange bool) error {
	address := os.Getenv("LOCATOR_API_INFO")
	opts := []config.Option{
		config.AddressOption(address),
	}

	if isRange {
		opts = append(opts, config.TraversalModeOption(config.TraversalModeRange))
	}

	client, err := titan.New(opts...)
	if err != nil {
		return err
	}

	size, reader, err := client.GetFile(context.Background(), cid)
	if err != nil {
		return err
	}
	defer reader.Close()

	if output == "" {
		io.Copy(io.Discard, reader)
		return nil
	}

	fileName := output
	if isRange {
		fileName = fmt.Sprintf("%s.car", output)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	bar := pb.New64(size).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	barR := bar.NewProxyReader(reader)

	bar.Start()
	defer bar.Finish()

	_, err = io.Copy(file, barR)
	if err != nil {
		return err
	}

	if !isRange {
		return nil
	}

	decodeCARFile(fileName, output)

	return os.Remove(fileName)
}
