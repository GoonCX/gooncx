package server

var usage = `
Usage: gooncx -bind-addr 0.0.0.0:4000
`

type Server struct {
}

func Usage() string {
  return usage
}

func ReleaseVersion() string {
  return "0.0.1"
}
