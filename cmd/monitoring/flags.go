package monitoring

import "github.com/spf13/pflag"

// Flags denotes command-line flags for monitoring configuration
type Flags struct {
	Service string
	Errors  bool
	Profile bool
}

// Attach attaches Flags' variables to the given flagset
func (f *Flags) Attach(flags *pflag.FlagSet, service string) {
	flags.StringVar(&f.Service, "service-name", service, "set service name for reporting")
	flags.BoolVar(&f.Errors, "errors", false, "enable error reporting")
	flags.BoolVar(&f.Profile, "profile", false, "enable profiling")
}
