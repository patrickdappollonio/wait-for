package sources

type Logger interface {
	Printf(fmt string, args ...interface{})
}

const (
	downPrefix = "> down:"
	upPrefix   = "> up:  "
)
