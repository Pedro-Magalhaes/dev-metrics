package metrics

// Essas variáveis serão preenchidas no momento do build via -ldflags
var (
	Version   = "dev"
	GitCommit = "none"
	BuildTime = "unknown"
)
