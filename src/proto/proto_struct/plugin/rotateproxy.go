package plugin

const (
	RotateProxyPluginName = "RotateProxy"
	RotateProxyImageName  = "pokemonscan/pokeball_rotateproxy"
)

type RotateProxyConfig struct {
	CommandArgs string `json:"command_args"`
	ListenPort  int    `json:"listen_port"`
	Email       string `json:"email"`
	Token       string `json:"token"`
	Check       string `json:"check"`
	CheckWords  string `json:"check_words"`
	SocksUser   string `json:"socks_user"`
	SocksPasswd string `json:"socks_passwd"`
}
