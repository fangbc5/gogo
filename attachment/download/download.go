package download

type DownLoad interface {
	// Init initializes options
	Init(options ...Option) error
	// The Logger options
	Options() Options
	// Download file from remote or local file system
	Load(fileId string) []byte
	// String returns the name of logger
	String() string
}