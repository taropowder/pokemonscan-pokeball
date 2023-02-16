package utils

type InfoResult struct {
	// key : hostname , value : port
	Hosts map[string][]int
	// key : Url
	Urls map[string]interface{}
	// taskId
	TaskId int32
}

func (r *InfoResult) Merge(newResult InfoResult) {
	for host, ports := range newResult.Hosts {
		if _, has := r.Hosts[host]; has {
			for _, p := range ports {
				if !IntInSlice(p, ports) {
					r.Hosts[host] = append(r.Hosts[host], p)
				}
			}
		} else {
			r.Hosts[host] = make([]int, 0)
			r.Hosts[host] = append(r.Hosts[host], ports...)
		}
	}

	for url, _ := range newResult.Urls {
		if _, ok := r.Urls[url]; url != "" && !ok {
			r.Urls[url] = url
		}
	}

}

//func (r InfoResult) Proto() *pokeball.ReportInfoArgs {
//	urls := make([]string, 0)
//	for url, _ := range r.Urls {
//		urls = append(urls, url)
//	}
//	hosts := make([]*pokeball.HostInfo, 0)
//	for host, port := range r.Hosts {
//		_ports := make([]int32, 0)
//		for _, p := range port {
//			_ports = append(_ports, int32(p))
//		}
//		//hosts = append(hosts, &bulbasaur.Host{
//		//	Host: host,
//		//	Port: _ports,
//		//})
//	}
//	report := &bulbasaur.ReportArgs{
//		Urls:   urls,
//		Hosts:  hosts,
//		TaskId: r.TaskId,
//	}
//	return report
//}
