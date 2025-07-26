package flags

import "flag"

type Flags struct {
	ResetProgress   bool
	FinishAudiobook bool
	Parsed          bool
}

const (
	FlagReset    = "reset"
	FlagComplete = "finish"
)

func New() *Flags {
	f := &Flags{}
	flag.BoolVar(&f.ResetProgress, FlagReset, false, "Reset the audiobook generation process")
	flag.BoolVar(&f.FinishAudiobook, FlagComplete, false, "Finish audiobook generation with currently processed chapters")
	flag.Parse()
	f.Parsed = true

	return f
}
