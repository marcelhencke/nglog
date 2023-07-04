package config

var Flags struct {
	FilePath             string
	Template             string
	DebugMode            bool
	DebugMaxLines        int
	LogLineCompleteRegex string
	LogLineFastCGIRegex  string
	LogLinePhpMsgSplit   string
	NginxVarRegex        string
}
