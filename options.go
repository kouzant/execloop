package execloop

type Options struct {
	Logger Logger
}

func DefaultOptions() Options {
	return Options{
		Logger: defaultLog,
	}
}

func (o Options) WithLogger(logger Logger) Options {
	o.Logger = logger
	return o
}
