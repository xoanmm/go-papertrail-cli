package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xoanmm/go-papertrail-cli/pkg/papertrail"
	"log"
	"os"
	"strings"
	"time"
)

// dateLayout define the layout to use with format
// mm/dd/yyyy HH:MM:SSx
const dateLayout = "01/02/2006 15:04:05"

var version = "1.2.0"
var date = time.Now().Format(time.RFC3339)
var now = time.Now().UTC()
var nowDate = now.Format(dateLayout)
var nowDateLessEightHours = now.Add(-8 * time.Hour).Format(dateLayout)

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
		Usage:    "interacts with papertrail through its api to perform both log collection actions and the creation/deletion of systems, groups and saved searches",
		Version:  version,
		Compiled: d,
		UsageText: "go-papertrail-cli [--group-name <group-name>] [--system-wildcard <wildcard>] " +
			"[--search <search-name>] [--query <query>] [--action <action>] " +
			"[--delete-all-searches <delete-all-searches>] [--delete-only-searches <delete-only-searches>] " +
			"[--delete-all-systems <delete-all-systems>]  [--delete-only-systems <delete-only-systems>]" +
			"[--start-date <start-date>] [--end-date <end-date>] [--path <path>]",
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
				Name:    "destination-port",
				Usage:   "destination port for sending the logs of the indicated system/s",
				Value:   0,
				Aliases: []string{"p"},
			},

			&cli.IntFlag{
				Name:    "destination-id",
				Usage:   "destination id for sending the logs of the indicated system/s",
				Value:   0,
				Aliases: []string{"I"},
			},

			&cli.StringFlag{
				Name:    "ip-address",
				Usage:   "source ip address from sending the logs of the indicated system/s",
				Value:   "",
				Aliases: []string{"i"},
			},

			&cli.StringFlag{
				Name:    "system-type",
				Usage:   "Type of system, can be hostname or ip-address",
				Value:   "hostname",
				Aliases: []string{"t"},
			},

			&cli.StringFlag{
				Name:    "search",
				Usage:   "name of saved search to be performed on logs or to be created on a group",
				Value:   "default search",
				Aliases: []string{"S"},
			},

			&cli.StringFlag{
				Name:    "query",
				Usage:   "query to be performed on the group of logs or applied on the search to be created",
				Value:   "*",
				Aliases: []string{"q"},
			},

			&cli.StringFlag{
				Name:    "action",
				Usage:   "Action to be performed with the information provided for papertrail, possible values only c(create), o(obtain) or d(delete)",
				Value:   "c",
				Aliases: []string{"a"},
			},

			&cli.BoolFlag{
				Name:    "delete-all-searches",
				Usage:   "Indicates if all searches in a group or a specific search are going to be deleted",
				Value:   false,
				Aliases: []string{"d"},
			},

			&cli.BoolFlag{
				Name:  "delete-only-searches",
				Usage: "Indicates if only searches specified are going to be deleted",
				Value: false,
			},

			&cli.BoolFlag{
				Name:    "delete-all-systems",
				Usage:   "Indicates if all systems specified are going to be deleted",
				Value:   true,
				Aliases: []string{"D"},
			},

			&cli.BoolFlag{
				Name:  "delete-only-systems",
				Usage: "Indicates if only systems specified are going to be deleted",
				Value: false,
			},

			&cli.StringFlag{
				Name:        "start-date",
				Usage:       "filter only from a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time)",
				DefaultText: "$ACTUAL_DATE - 8hours",
				Value:       nowDateLessEightHours,
				Aliases:     []string{"s"},
			},

			&cli.StringFlag{
				Name:        "end-date",
				Usage:       "filter only until a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time)",
				DefaultText: "$ACTUAL_DATE",
				Value:       nowDate,
				Aliases:     []string{"e"},
			},

			&cli.StringFlag{
				Name:    "path",
				Usage:   "path where to store the logs",
				Value:   "/tmp",
				Aliases: []string{"P"},
			},
		},
		Action: func(c *cli.Context) error {
			// path, _ := filepath.Abs(c.String("path"))
			logGroupName := c.String("group-name")
			actionName := c.String("action")

			papertrailActions, action, err := app.PapertrailActions(&papertrail.Options{
				GroupName:          logGroupName,
				SystemWildcard:     c.String("system-wildcard"),
				DestinationPort:    c.Int("destination-port"),
				DestinationId:      c.Int("destination-id"),
				IpAddress:          c.String("ip-address"),
				SystemType:         c.String("system-type"),
				Search:             c.String("search"),
				Query:              c.String("query"),
				Action:             actionName,
				DeleteAllSystems:   c.Bool("delete-all-systems"),
				DeleteOnlySystems:  c.Bool("delete-only-systems"),
				DeleteAllSearches:  c.Bool("delete-all-searches"),
				DeleteOnlySearches: c.Bool("delete-only-searches"),
				StartDate:          c.String("start-date"),
				EndDate:            c.String("end-date"),
				Path:               c.String("path"),
			})
			printFinalResultIfNotErrorsDetected(err, action, papertrailActions)
			return err
		},
	}
}

func printFinalResultIfNotErrorsDetected(err error, actionName *string, papertrailActions []papertrail.Item) {
	if err == nil && actionName != nil {
		if !papertrail.ActionIsObtain(*actionName) {
			if len(papertrailActions) > 0 {
				log.Printf("%s actions have been carried out on the following elements\n", strings.Title(*actionName))
				for _, item := range papertrailActions {
					log.Printf("- %s with ID %d and name '%s'\n", item.ItemType, item.ID, item.ItemName)
				}
			}
		} else {
			log.Printf("%s saved in file %s", papertrailActions[0].ItemType, papertrailActions[0].ItemName)
		}
	}
}
