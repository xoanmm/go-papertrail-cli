package main

import (
	"fmt"
	"github.com/xoanmm/go-papertrail-cli/pkg/papertrail"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var version = "1.0.0"
var date = time.Now().Format(time.RFC3339)

func main() {
	cmd := buildCLI(&papertrail.App{})

	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// buildCLI creates a CLI app
func buildCLI(app *papertrail.App) *cli.App {
	d, _ := time.Parse(time.RFC3339, date)
	return &cli.App{
		Name:     "go-papertrail-cli",
		Usage:    "interacts with papertrail through its api to perform both log collection actions and the creation of systems, groups and saved searches",
		Version:  version,
		Compiled: d,
		UsageText: "go-papertrail-cli [--group-name <group-name>] [--system-wildcard <wildcard>] " +
			"[--search <search-name>] [--query <query>]",
		Authors: []*cli.Author{
			{
				Name:  "Xoan Mallon",
				Email: "xoanmallon@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "group-name",
				Usage:   "group defined or to be defined in papertrail",
				Value:   "my-log-group",
				Aliases: []string{"g"},
			},

			&cli.StringFlag{
				Name:    "system-wildcard",
				Usage:   "wildcard to be applied on the systems defined in papertrail",
				Value:   "*",
				Aliases: []string{"w"},
			},

			&cli.IntFlag{
				Name:	"destination-port",
				Usage: 	"destination port for sending the logs of the indicated system/s",
				Value:	0,
				Aliases: []string{"p"},
			},

			&cli.IntFlag{
				Name:	"destination-id",
				Usage: 	"destination id for sending the logs of the indicated system/s",
				Value:	0,
				Aliases: []string{"I"},
			},

			&cli.StringFlag{
				Name:	"ip-address",
				Usage: 	"source ip address from sending the logs of the indicated system/s",
				Value:	"",
				Aliases: []string{"i"},
			},

			&cli.StringFlag{
				Name:    "system-type",
				Usage: 	 "Type of system, can be hostname or ip-address",
				Value:   "system type",
				Aliases: []string{"t"},
			},

			&cli.StringFlag{
				Name:    "search",
				Usage: 	 "name of saved search to be performed on logs or to be created on a group",
				Value:   "saved search",
				Aliases: []string{"S"},
			},

			&cli.StringFlag{
				Name:  "query",
				Usage: "query to be performed on the group of logs or applied on the search to be created",
				Value: "*",
				Aliases: []string{"q"},
			},

			&cli.StringFlag{
				Name: 	 "action",
				Usage: 	 "Action to be performed with the information provided for papertrail, possible values only c(create) and o(obtain)",
				Value:   "c",
				Aliases: []string{"a"},
			},
		},
		Action: func(c *cli.Context) error {
			// path, _ := filepath.Abs(c.String("path"))
			logGroupName := c.String("group-name")

			papertrailActions, err := app.PapertrailNecessaryActions(&papertrail.Options{
				GroupName:              logGroupName,
				SystemWildcard:			c.String("system-wildcard"),
				DestinationPort:		c.Int("destination-port"),
				DestinationId:			c.Int("destination-id"),
				IpAddress:				c.String("ip-address"),
				SystemType:				c.String("system-type"),
				Search:					c.String("search"),
				Query:					c.String("query"),
				Action:					c.String("action"),
			})
			if len(papertrailActions) > 0 {
				fmt.Println("The created items are the following")
				for _, item := range papertrailActions {
					fmt.Printf("- %s with ID %d and name '%s'\n", item.ItemType, item.ID, item.ItemName)
				}
			}
			return err
		},
	}
}