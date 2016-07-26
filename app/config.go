package app

type Config struct {
	Host        string
	Port        int
	RootDir     string
	Index       string
	Token       string
	ExcludeDirs []string
}
