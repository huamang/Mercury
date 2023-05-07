package common

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Parse() {
	app := &cli.App{
		Name:    "Mercury",
		Usage:   "A scanner",
		Version: "v1.0",
		Action: func(context *cli.Context) error {
			cli.ShowAppHelp(context)
			os.Exit(0)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "output file",
				Destination: &Output,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "debug mode",
				Value:       false,
				Destination: &Debug,
			},
		},

		Commands: []*cli.Command{
			{
				Name:  "scan",
				Usage: "scan mode",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "target",
						Aliases:     []string{"t"},
						Usage:       "target host, eg: CIDR or \"-\" or \",\" to split",
						Destination: &Host,
					},
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Usage:       "target file, eg: /tmp/target.txt",
						Destination: &HostFile,
					},
					&cli.IntFlag{
						Name:        "thread",
						Aliases:     []string{"T"},
						Usage:       "number of concurrent threads",
						Value:       100,
						Destination: &Thread,
					},
					&cli.StringFlag{
						Name:        "port",
						Aliases:     []string{"p"},
						Usage:       "appoint ports, eg: 80,8080 or 80-88",
						Destination: &Port,
					},
					&cli.BoolFlag{
						Name:        "icmp",
						Usage:       "icmp mode scan,default ping mode",
						Value:       false,
						Destination: &Icmp,
					},
					&cli.BoolFlag{
						Name:        "check",
						Usage:       "just check alive",
						Value:       false,
						Destination: &CheckAlive,
					},
				},
				Action: func(context *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
