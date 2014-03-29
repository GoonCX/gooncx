package config

import (
  "flag"
  "fmt"
  "strings"
  "strconv"
  "reflect"
  "os"
  "io/ioutil"
  "github.com/BurntSushi/toml"
)

const DefaultConfigPath = "/etc/gooncx/gooncx.conf"

type Config struct {
  SystemPath       string
  BindAddr         string `toml:"bind_addr" env:"GOONCX_BIND_ADDR"`
  ShowHelp         bool
  ShowVersion      bool
  Verbose          bool `toml:"verbose" env:"ETCD_VERBOSE"`
  VeryVerbose      bool `toml:"very_verbose" env:"ETCD_VERY_VERBOSE"`
}

func New() *Config {
  c := new(Config)
  c.SystemPath = DefaultConfigPath
  c.BindAddr = "127.0.0.1:4001"
  return c
}

func (c *Config) Load(arguments []string) error {
  var path string
  f := flag.NewFlagSet("gooncx", -1)
  f.SetOutput(ioutil.Discard)
  f.StringVar(&path, "config", "", "path to config file")
  f.Parse(arguments)
  if path != "" {
    if err := c.LoadFile(path); err != nil {
      return err
    }
  }

  if err := c.LoadEnv(); err != nil {
    return err
  }
  if err := c.LoadFlags(arguments); err != nil {
    return err
  }

  return nil
}

// Load config from env
func (c *Config) LoadEnv() error {
  if err := c.loadEnv(c); err != nil {
    return err
  }
  return nil
}


// Loads config from file
func (c *Config) LoadFile(path string) error {
  _, err := toml.DecodeFile(path, &c)
  return err
}

// Loads config from args
func (c *Config) LoadFlags(arguments []string) error {
  var path string
  f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
  f.SetOutput(ioutil.Discard)
  f.BoolVar(&c.ShowHelp, "h", false, "")
  f.BoolVar(&c.ShowHelp, "help", false, "")
  f.BoolVar(&c.ShowVersion, "version", false, "")
  f.BoolVar(&c.Verbose, "v", c.Verbose, "")
  f.StringVar(&c.BindAddr, "bind-addr", c.BindAddr, "")
  f.StringVar(&path, "config", "", "")

  if err := f.Parse(arguments); err != nil {
    return err
  }
  return nil
}

func (c *Config) loadEnv(target interface{}) error {
  value := reflect.Indirect(reflect.ValueOf(target))
  typ := value.Type()
  for i := 0; i < typ.NumField(); i++ {
    field := typ.Field(i)

    // Retrieve environment variable.
    v := strings.TrimSpace(os.Getenv(field.Tag.Get("env")))
    if v == "" {
      continue
    }

    // Set the appropriate type.
    switch field.Type.Kind() {
    case reflect.Bool:
      value.Field(i).SetBool(v != "0" && v != "false")
    case reflect.Int:
      newValue, err := strconv.ParseInt(v, 10, 0)
      if err != nil {
        return fmt.Errorf("Parse error: %s: %s", field.Tag.Get("env"), err)
      }
      value.Field(i).SetInt(newValue)
    case reflect.String:
      value.Field(i).SetString(v)
    }
  }
  return nil
}
