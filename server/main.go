package main

import (
	"easy-share/server"
	"easy-share/utils"
	"os"
	"path/filepath"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		utils.LogFatal(err)
	}
	err = os.Chdir(filepath.Dir(ex))
	if err != nil {
		utils.LogFatal(err)
	}

	params, _ := utils.SmartArgs("--data|-d=/data:,--port|-p=80:,--test,--sudo", os.Args[1:])
	dataDir := params["--data"].GetValue()
	server.Testing = params["--test"].IsSet()
	utils.UseSudo = params["--sudo"].IsSet()

	err = server.Init(dataDir)
	if err != nil {
		utils.LogFatal(err)
	}

	utils.AddSignalHandler([]os.Signal{server.SHUTDOWN}, func(sig os.Signal) {
		utils.LogF("Received signal %s\n", sig)
		switch sig {
		case server.SHUTDOWN:
			os.Exit(0)
		}
	})

	// Launch webserver
	server.WebServer(params["--port"].GetValue())
}
