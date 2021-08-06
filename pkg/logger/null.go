package logger

var _ Logger = &nullLogger{}

// NewNullLogger returns a logger that does log anything
func NewNullLogger() FullLogger {
	return &nullLogger{}
}

type nullLogger struct{}

func (n *nullLogger) Trace(_ ...interface{}) {}

func (n *nullLogger) Debug(_ ...interface{}) {}

func (n *nullLogger) Print(_ ...interface{}) {}

func (n *nullLogger) Info(_ ...interface{}) {}

func (n *nullLogger) Warn(_ ...interface{}) {}

func (n *nullLogger) Warning(_ ...interface{}) {}

func (n *nullLogger) Error(_ ...interface{}) {}

func (n *nullLogger) Fatal(_ ...interface{}) {}

func (n *nullLogger) Panic(_ ...interface{}) {}

func (n *nullLogger) Tracef(_ string, _ ...interface{}) {}

func (n *nullLogger) Debugf(_ string, _ ...interface{}) {}

func (n *nullLogger) Infof(_ string, _ ...interface{}) {}

func (n *nullLogger) Printf(_ string, _ ...interface{}) {}

func (n *nullLogger) Warnf(_ string, _ ...interface{}) {}

func (n *nullLogger) Warningf(_ string, _ ...interface{}) {}

func (n *nullLogger) Errorf(_ string, _ ...interface{}) {}

func (n *nullLogger) Fatalf(_ string, _ ...interface{}) {}

func (n *nullLogger) Panicf(_ string, _ ...interface{}) {}

func (n *nullLogger) Traceln(_ ...interface{}) {}

func (n *nullLogger) Debugln(_ ...interface{}) {}

func (n *nullLogger) Infoln(_ ...interface{}) {}

func (n *nullLogger) Println(_ ...interface{}) {}

func (n *nullLogger) Warnln(_ ...interface{}) {}

func (n *nullLogger) Warningln(_ ...interface{}) {}

func (n *nullLogger) Errorln(_ ...interface{}) {}

func (n *nullLogger) Fatalln(_ ...interface{}) {}

func (n *nullLogger) Panicln(_ ...interface{}) {}
