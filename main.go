package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	"./cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "eos-net"
	app.Usage = "To boot or join a eos network."
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Han Fei",
			Email: "hanfei1009@126.com",
		},
	}

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "start",
			Aliases: []string{"s"},
			Usage: "start the eos network",
			Action: func(c *cli.Context) error {
				cmd.StartNodeos("test")
				return nil
			},
		},
		{
			Name: "boot",
			Aliases: []string{"b"},
			Usage: "boot the eos network",
			Action: func(c *cli.Context) error {
				cmd.Boot()
				return nil
			},
		},
		{
			Name: "join",
			Aliases: []string{"j"},
			Usage: "join the eos network",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
		{
			Name: "vote",
			Aliases: []string{"v"},
			Usage: "vote the producer",
			Action: func(c *cli.Context) error {
				cmd.Vote()
				return nil
			},
		},
		{
			Name: "resign",
			Aliases: []string{"r"},
			Usage: "resign the eosio to eosio.prods",
			Action: func(c *cli.Context) error {
				cmd.Resign()
				return nil
			},
		},
		{
			Name: "test",
			Aliases: []string{"t"},
			Usage: "test",
			Action: func(c *cli.Context) error {
				cmd.Test()
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("Use --help or -h to see what you should do!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
