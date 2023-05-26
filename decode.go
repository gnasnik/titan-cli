package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-merkledag"
	unixfile "github.com/ipfs/go-unixfs/file"
	carv2 "github.com/ipld/go-car/v2"
	"github.com/ipld/go-car/v2/blockstore"
	"io"
	"log"
	"os"
)

func decodeCARFile(CARFilePath, outputPath string) {
	fmt.Printf("Decoding CAR file %s to %s ...\n", CARFilePath, outputPath)

	// Open the CAR file
	carFile, err := os.Open(CARFilePath)
	if err != nil {
		fmt.Println("Error opening CAR file:", err)
		return
	}
	defer carFile.Close()

	// Create a blockstore
	bs, err := blockstore.OpenReadOnly(CARFilePath)
	if err != nil {
		fmt.Println("Error creating blockstore from CAR file:", err)
		return
	}

	// Create a blockservice
	blockService := blockservice.New(bs, offline.Exchange(bs))

	// Create a merkledag service
	dagService := merkledag.NewDAGService(blockService)

	// Get the root CID of the CAR file
	rootsReader, err := carv2.NewReader(carFile)
	if err != nil {
		fmt.Println("Error creating roots reader:", err)
		return
	}
	rootCIDs, err := rootsReader.Roots()
	if err != nil {
		fmt.Println("Error getting root CIDs:", err)
		return
	}

	// Get the IPLD node from the root CID
	node, err := dagService.Get(context.Background(), rootCIDs[0])
	if err != nil {
		fmt.Println("Error getting IPLD node from root CID:", err)
		return
	}

	mp4File, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	merkleNode, err := unixfile.NewUnixfsFile(context.Background(), dagService, node)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(mp4File, files.ToFile(merkleNode))
}
