package main

import (
	"encoding/xml"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getXML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read body: %v", err)
	}

	return string(data), nil
}

type Vds struct {
	id         int     `xml:"elem>id"`
	ostempl    string  `xml:"elem>ostempl"`
	ip         string  `xml:"elem>ip"`
	domain     string  `xml:"elem>domain"`
	intname    string  `xml:"elem>intname"`
	item_cost  float64 `xml:"elem>item_cost"`
	pricelist  string  `xml:"elem>pricelist"`
	status     int     `xml:"elem>status"`
	createdate string  `xml:"elem>createdate"`
}

func main() {

	app := cli.NewApp()
	app.Name = "fcli"
	app.Version = "0.1"
	app.Usage = "Cli tool for manage Firstvds BillManager!"

	var authinfo string
	var url string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "url, u",
			Value:       "https://api.firstvds.ru/billmgr",
			Usage:       "Url for request",
			Destination: &url,
		},
		cli.StringFlag{
			Name:        "authinfo, a",
			Usage:       "user:password for the billmgr access",
			FilePath:    "~/.firstauthinfo",
			Destination: &authinfo,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list of current items",
			Action: func(c *cli.Context) error {
				if len(authinfo) < 4 {
					fmt.Println("WARN: no auth data available. You should use -a param or ~/.firstauthinfo file ")
					os.Exit(0)
				}
				url := url + "?authinfo=" + authinfo + "&func=vds&out=xml"
				outXml, err := getXML(url)
				if err != nil {
					print(err.Error())
				}
				decoder := xml.NewDecoder(outXml)
				err = decoder.Decode(&Product)
				print(outXml)
				return nil
			},
		},
		{
			Name:    "tariff",
			Aliases: []string{"t"},
			Usage:   "show list of available tariffs",
			Action: func(c *cli.Context) error {
				fmt.Println("list of tariffs: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "vds",
			Aliases: []string{"v"},
			Usage:   "action for VDS products",
			Subcommands: []cli.Command{
				{
					Name:    "order",
					Aliases: []string{"o"},
					Usage:   "order a new VDS",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:    "remove",
					Aliases: []string{"r"},
					Usage:   "remove an existing item",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
