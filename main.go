package main

import (
	"github.com/jmattheis/shargo/client"
	"github.com/jmattheis/shargo/encrypt"
	"github.com/jmattheis/shargo/host"
	"image"
	"log"
	"os"

	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"github.com/urfave/cli/v2"
)

func main() {
	a := cli.App{
		Name:  "shargo",
		Usage: "share your screen",
		Commands: []*cli.Command{
			{
				Name: "host",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "addr",
						Value: ":55555",
					},
					&cli.IntFlag{
						Name:  "display",
						Value: 0,
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "",
					},
				},
				Action: func(ctx *cli.Context) error {
					images := make(chan image.Image, 1)
					go host.ImageStream(ctx.Int("display"), images)
					return host.Server(ctx.String("addr"), images, encrypt.Sha256(ctx.String("password")))
				},
			},
			{
				Name: "connect",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "addr",
						Value: ":55555",
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "",
					},
				},
				Action: func(ctx *cli.Context) error {
					images := make(chan image.Image, 1)
					go func() {
						log.Fatal(client.Client(ctx.String("addr"), encrypt.Sha256(ctx.String("password")), images))
					}()
					gui := app.New()
					window := gui.NewWindow("Sharing")
					i := &canvas.Image{FillMode: canvas.ImageFillOriginal}
					window.SetContent(i)
					window.Show()

					go func() {
						for img := range images {
							i.Image = img
							i.Refresh()
						}
					}()
					gui.Run()
					return nil
				},
			},
		},
	}
	log.Fatal(a.Run(os.Args))

}
