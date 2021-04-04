package main

import (
	_ "github.com/CodyGuo/godaemon"
	"github.com/iplcm/access/routes"
)

func main() {
	routes.Init().SetSession().AddRoutes().Run()
}
