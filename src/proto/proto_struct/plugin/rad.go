package plugin

const (
	RadPluginName        = "Rad"
	RadImageName         = "pokemonscan/pokeball_rad"
	RadDefaultConfigFile = `exec_path: ""                     # 启动chrome的路径
disable_headless: false           # 禁用无头模式
force_sandbox: false              # 强制开启sandbox；为 false 时默认开启沙箱，但在容器中会关闭沙箱。为true时强制启用沙箱，可能导致在docker中无法使用。
enable_image: false               # 启用图片显示
parent_path_detect: true          # 是否启用父目录探测功能
proxy: ""                         # 代理配置
user_agent: ""                    # 请求user-agent配置
domain_headers:                   # 请求头配置:[]{domain,map[headerKey]HeaderValue}
  - domain: '*'                   # 为哪些域名设置header，glob语法
    headers: %s                    
max_depth: 10                     # 最大页面深度限制
navigate_timeout_second: 10       # 访问超时时间，单位秒
load_timeout_second: 10           # 加载超时时间，单位秒
retry: 0                          # 页面访问失败后的重试次数
page_analyze_timeout_second: 300  # 页面分析超时时间，单位秒
max_interactive: 1000             # 单个页面最大交互次数
max_interactive_depth: 10         # 页面交互深度限制
max_page_concurrent: 10           # 最大页面并发（不大于10）
max_page_visit: 1000              # 总共允许访问的页面数量
max_page_visit_per_site: 1000     # 每个站点最多访问的页面数量
element_filter_strength: 0        # 过滤同站点相似元素强度，1-7取值，强度逐步增大，为0时不进行跨页面元素过滤
new_task_filter_config: # 检查某个链接是否应该被加入爬取队列
  %s
  hostname_disallowed: [ ]         # 不允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  port_allowed: [ ]                # 允许访问的端口, 支持的格式如: 80、80-85
  port_disallowed: [ ]             # 不允许访问的端口, 支持的格式如: 80、80-85
  path_allowed: [ ]                # 允许访问的路径，支持的格式如: test、*test*
  path_disallowed: [ ]             # 不允许访问的路径, 支持的格式如: test、*test*
  query_key_allowed: [ ]           # 允许访问的 Query Key，支持的格式如: test、*test*
  query_key_disallowed: [ ]        # 不允许访问的 Query Key, 支持的格式如: test、*test*
  fragment_allowed: [ ]            # 允许访问的 Fragment, 支持的格式如: test、*test*
  fragment_disallowed: [ ]         # 不允许访问的 Fragment, 支持的格式如: test、*test*
  post_key_allowed: [ ]            # 允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
  post_key_disallowed: [ ]         # 不允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
request_send_filter_config: # 检查某个请求是否应该被发送
  hostname_allowed: [ ]            # 允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  hostname_disallowed: [ ]         # 不允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  port_allowed: [ ]                # 允许访问的端口, 支持的格式如: 80、80-85
  port_disallowed: [ ]             # 不允许访问的端口, 支持的格式如: 80、80-85
  path_allowed: [ ]                # 允许访问的路径，支持的格式如: test、*test*
  path_disallowed: [ ]             # 不允许访问的路径, 支持的格式如: test、*test*
  query_key_allowed: [ ]           # 允许访问的 Query Key，支持的格式如: test、*test*
  query_key_disallowed: [ ]        # 不允许访问的 Query Key, 支持的格式如: test、*test*
  fragment_allowed: [ ]            # 允许访问的 Fragment, 支持的格式如: test、*test*
  fragment_disallowed: [ ]         # 不允许访问的 Fragment, 支持的格式如: test、*test*
  post_key_allowed: [ ]            # 允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
  post_key_disallowed: [ ]         # 不允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
request_output_filter_config: # 检查某个请求是否应该被输出
  hostname_allowed: [ ]            # 允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  hostname_disallowed: [ ]         # 不允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  port_allowed: [ ]                # 允许访问的端口, 支持的格式如: 80、80-85
  port_disallowed: [ ]             # 不允许访问的端口, 支持的格式如: 80、80-85
  path_allowed: [ ]                # 允许访问的路径，支持的格式如: test、*test*
  path_disallowed: [ ]             # 不允许访问的路径, 支持的格式如: test、*test*
  query_key_allowed: [ ]           # 允许访问的 Query Key，支持的格式如: test、*test*
  query_key_disallowed: [ ]        # 不允许访问的 Query Key, 支持的格式如: test、*test*
  fragment_allowed: [ ]            # 允许访问的 Fragment, 支持的格式如: test、*test*
  fragment_disallowed: [ ]         # 不允许访问的 Fragment, 支持的格式如: test、*test*
  post_key_allowed: [ ]            # 允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
  post_key_disallowed: [ ]         # 不允许访问的 Post Body 中的参数, 支持的格式如: test、*test*`
)

type RadConfig struct {
	Target           string `json:"target"`
	AllowDomains     string `json:"allow_domains"`
	DownstreamPlugin string `json:"downstream_plugin"`
	RadConfigFile    string `json:"rad_config_file"`
	Cookie           string `json:"cookie"`
}
