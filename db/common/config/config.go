package config

type DatabaseOpts struct {
	Type string
	// <ip>:<port>
	Host         string
	Name         string
	User         string
	Password     string
	Path         string
	MaxOpenConns int
	MaxIdleConns int
	Charset      string
}

var Database DatabaseOpts
