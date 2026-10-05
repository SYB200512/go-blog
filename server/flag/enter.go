package flag

import (
	"os"
	"server/global"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	sqlFlag = &cli.BoolFlag{
		Name:  "sql",
		Usage: "Initializes the MySQL database structure",
	}
)

func Run(c *cli.Context) {
	switch {
	case c.Bool(sqlFlag.Name):
		if err := SQL(); err != nil {
			global.Log.Error("Failed to initialize the MySQL database structure:", zap.Error(err))
			return
		} else {
			global.Log.Info("Successfully initialized the MySQL database structure")
		}
	}
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Go Blog"
	app.Flags = []cli.Flag{
		sqlFlag,
	}
	app.Action = Run
	return app
}

func InitFlag() {
	app := NewApp()
	err := app.Run(os.Args)
	if err != nil {
		global.Log.Error("Application execution encountered an error:", zap.Error(err))
		os.Exit(1)
	}
	os.Exit(0)
}
