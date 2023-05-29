package plugin

const (
	CommonPluginName = "Common"
)

type CommonConfig struct {
	Image  string `json:"image"`
	Config struct {
		ConfigPath    string `json:"config_path"`
		ConfigContent string `json:"config_content"`
	} `json:"config"`
	ResultPath string `json:"result_path"`
	Command    string `json:"command"`
	File       struct {
		FilePath    string `json:"file_path"`
		FileContent string `json:"file_content"`
	} `json:"file"`
}
