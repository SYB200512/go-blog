package flag

import (
	"errors"
	"fmt"
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
	sqlExportFlag = &cli.BoolFlag{
		Name:  "sql-export",
		Usage: "Exports SQL data to a specified file.",
	}
	sqlImportFlag = &cli.BoolFlag{
		Name:  "sql-import",
		Usage: "Imports the SQL data from a specified file.",
	}
	esFlag = &cli.BoolFlag{
		Name:  "es",
		Usage: "Initializes the Elasticsearch index.",
	}
	esExportFlag = &cli.BoolFlag{
		Name:  "es-export",
		Usage: "Exports data from Elasticsearch to a specified file.",
	}
	esImportFlag = &cli.StringFlag{
		Name:  "es-import",
		Usage: "Imports data into Elasticsearch from a specified file.",
	}
	adminFlag = &cli.BoolFlag{
		Name:  "admin",
		Usage: "Creates an administrator using the name, email and address specified in the config.yaml file.",
	}
)

func Run(c *cli.Context) {
	// 检查是否有多个标志被设置
	if c.NumFlags() > 1 {
		err := cli.NewExitError("Only one flag is allowed", 1)
		global.Log.Error("Invalid command usage:", zap.Error(err))
		os.Exit(1)
	}
	// 检查是否有标志被设置
	switch {
	// 初始化 MySQL 数据库结构
	case c.Bool(sqlFlag.Name):
		if err := SQL(); err != nil {
			global.Log.Error("Failed to initialize the MySQL database structure:", zap.Error(err))
			return
		} else {
			global.Log.Info("Successfully initialized the MySQL database structure")
		}
		// 导出 SQL 数据
	case c.Bool(sqlExportFlag.Name):
		if err := SQLExport(); err != nil {
			global.Log.Error("Failed to export SQL data:", zap.Error(err))
		} else {
			global.Log.Info("Successfully exported SQL data")
		}
		// 导入 SQL 数据
	case c.IsSet(sqlImportFlag.Name):
		if errs := SQLImport(c.String(sqlImportFlag.Name)); len(errs) > 0 {
			var combinedErrors string
			for _, err := range errs {
				combinedErrors += err.Error() + "\n"
			}
			err := errors.New(combinedErrors)
			global.Log.Error("Failed to import SQL data:", zap.Error(err))
		} else {
			global.Log.Info("Successfully imported SQL data")
		}
		// 初始化 Elasticsearch 索引
	case c.Bool(esFlag.Name):
		if err := Elasticsearch(); err != nil {
			global.Log.Error("Failed to create ES indices:", zap.Error(err))
		} else {
			global.Log.Info("Successfully created ES indices")
		}
	case c.Bool(esExportFlag.Name):
		if err := ElasticsearchExport(); err != nil {
			global.Log.Error("Failed to export ES data:", zap.Error(err))
		} else {
			global.Log.Info("Successfully exported ES data")
		}
	case c.IsSet(esImportFlag.Name):
		if num, err := ElasticsearchImport(c.String(esImportFlag.Name)); err != nil {
			global.Log.Error("Failed to import ES data:", zap.Error(err))
		} else {
			global.Log.Info(fmt.Sprintf("Successfully imported ES data, totaling %d records", num))
		}
	case c.Bool(adminFlag.Name):
		if err := Admin(); err != nil {
			global.Log.Error("Failed to create an administrator:", zap.Error(err))
		} else {
			global.Log.Info("Successfully created an administrator")
		}
	default:
		err := cli.NewExitError("unknown command", 1)
		global.Log.Error(err.Error(), zap.Error(err))
	}
}

// NewApp 创建 CLI 应用
func NewApp() *cli.App {

	app := cli.NewApp()
	app.Name = "Go Blog"
	app.Flags = []cli.Flag{
		sqlFlag,
		sqlExportFlag,
		sqlImportFlag,
		esFlag,
		esExportFlag,
		esImportFlag,
		adminFlag,
	}
	app.Action = Run
	return app
}

// InitFlag 初始化标志
func InitFlag() {
	if len(os.Args) > 1 {
		app := NewApp()
		err := app.Run(os.Args)
		if err != nil {
			global.Log.Error("Application execution encountered an error:", zap.Error(err))
			os.Exit(1)
		}

		if os.Args[1] == "-h" || os.Args[1] == "-help" {
			fmt.Println("Display help message...")
		}
		os.Exit(0)
	}
}
