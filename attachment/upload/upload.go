package upload

type Upload interface {
	// Init initializes options
	Init(options ...Option) error
	// The Logger options
	Options() Options
	// Download file from remote or local file system
	Save([]byte) string
	// String returns the name of logger
	String() string
}