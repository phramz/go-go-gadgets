package logger

// Logger log abstraction
type Logger interface {
	Debug(_ ...interface{})
	Print(_ ...interface{})
	Info(_ ...interface{})
	Warn(_ ...interface{})
	Warning(_ ...interface{})
	Error(_ ...interface{})
	Fatal(_ ...interface{})
	Panic(_ ...interface{})
}

// Console log abstraction
type Console interface {
	Debugln(_ ...interface{})
	Infoln(_ ...interface{})
	Println(_ ...interface{})
	Warnln(_ ...interface{})
	Warningln(_ ...interface{})
	Errorln(_ ...interface{})
	Fatalln(_ ...interface{})
	Panicln(_ ...interface{})
}

// Format log abstraction
type Format interface {
	Debugf(_ string, _ ...interface{})
	Infof(_ string, _ ...interface{})
	Printf(_ string, _ ...interface{})
	Warnf(_ string, _ ...interface{})
	Warningf(_ string, _ ...interface{})
	Errorf(_ string, _ ...interface{})
	Fatalf(_ string, _ ...interface{})
	Panicf(_ string, _ ...interface{})
}

// Traceable log abstraction
type Traceable interface {
	Trace(_ ...interface{})
	Traceln(_ ...interface{})
	Tracef(_ string, _ ...interface{})
}

// FormatLogger log abstraction
type FormatLogger interface {
	Logger
	Format
}

// FullLogger log abstraction
type FullLogger interface {
	Logger
	Console
	Format
	Traceable
}
