package cmds

var (
	Flags struct {
		Verbose    bool     // -v, --verbose
		Extensions []string // --extensions, The file extensions cppgo should accept. Defaults to "pgo".
		GoBinary   string   // --go, Use this go executable instead of the one in the path
		CppBinary  string   // --cpp, Use this CPP instead of the one in the path
	}
)
