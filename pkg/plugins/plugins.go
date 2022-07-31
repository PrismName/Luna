package plugins

type PROTO string

var (
	HTTP PROTO = "HTTP"
	RDP  PROTO = "RDP"
)

type PluginsInfos struct {
	Name         string
	Auhtor       string
	VulnType     string
	VulnID       string
	VulnDate     string
	VulnRefrence string
	VulnDesc     string
}

func SetupPocs(v interface{}) {}

func SetupExps(v interface{}) {}
