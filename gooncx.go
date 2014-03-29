package main

import(
  "os"
  "fmt"
  "github.com/gooncx/gooncx/config"
  "github.com/gooncx/gooncx/server"
  "github.com/op/go-logging"
)

func main() {
  var log = logging.MustGetLogger("gooncx.com")
  var format = logging.MustStringFormatter("%{level} %{message}")
  logging.SetFormatter(format)
  logging.SetLevel(logging.INFO, "gooncx.com")
  var config = config.New()
  if err := config.Load(os.Args[1:]); err != nil {
    fmt.Println(server.Usage() + "\n")
    fmt.Println(err.Error() + "\n")
    os.Exit(1)
  } else if config.ShowVersion {
    fmt.Printf("gooncx version: %s\n", server.ReleaseVersion())
    os.Exit(0)
  } else if config.ShowHelp {
    fmt.Println(server.Usage() + "\n")
    os.Exit(0)
  } else if config.VeryVerbose {
    logging.SetLevel(logging.DEBUG, "gooncx.com")
  }
  log.Info("Bind address will be : %s\n", config.BindAddr)
}
