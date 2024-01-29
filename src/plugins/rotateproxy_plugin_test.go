package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"testing"
)

func TestRotateProxy(t *testing.T) {
	r := RotateProxyPlugin{}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前工作目录失败:", err)
		return
	}

	fmt.Println("当前工作目录:", cwd)

	filePath := "data/config/rotate_proxy.json"
	configStr, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	//configStr := "{\n  \"whitelist_type\": [\n    \"jpg\",\n    \"gif\",\n    \"png\",\n    \"css\"\n  ],\n  \"request_intercept_rules\": {\n    \"uri\": {\n      \"login_page\": \"login\"\n    },\n    \"headers\": {\n      \"Shiro\": \"deleteMe\"\n    },\n    \"parameters\": {\n      \"password_filed\": \"passwd\"\n    }\n  },\n  \"response_intercept_rules\": {\n    \"uri\": {\n    },\n    \"headers\": {\n      \"Shiro\": \"deleteMe\"\n    },\n    \"body\": {\n      \"password_filed\": \"passwd\"\n    }\n  }\n}"
	//configStr := "{\n  \"resp_intercept_rules\": [\n    {\n      \"url\": [\"js\"],\n      \"data\": [\"{path:\"]\n    }\n  ]\n}"
	//configStr := "{\n  \"resp_intercept_rules\": [\n    {\n      \"url\": [\"js\"],\n      \"data\": [\"this.message,name\"]\n    }\n  ],\n  \"port\": 8980\n}"
	//configStr := "{\n  \"resp_intercept_rules\": [\n    {\n      \"url\": [\"js\"],\n      \"data\": [\"exactPath\"]\n    }\n  ],\n  \"port\": 8090,\n  \"ca_file\": \"LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBMmYwTzdKOWx4bnFSVW1GNVdSS1RNbGVCM3JLSGZra2V5TlhNNXVmYk0vak0zOHNGCjAzeUozL1VWZ0tOeDZHdWVuLy94b0FkY1VpV0NUQlFKUzVPdXpCWlhjTTdnZExKVFFGb3RMdEwxWVUycVFSc2sKb3NxR3NHY09pVnY4ZS9TNnN6VkZyMjFpdEtmY0IwUVkySzJkZzlZYVZDT29BSUZucHhtTTlRdGRENGZtNUlXMQoyVDJTSjk5T29UV3hSMEU3QUxERmFjRmpUTnhrTVBGQnpCOFpvOG5ERVdaSVdKeUlkK0Z2VjAwc1dvTy9oNUhjCjM1QTFvMEJJeEpOTTNxNkU0bUllcjVDRVN5aU43bURjaHlYK0RDU1RURnMzNGxzNGxPd1g5Mi93Rlc1ZExzdXgKRmtiY1BDR2piZzJDeVk5Y3hCOFdVUS9uVW82bkFpU0NvaUJlTFFJREFRQUJBb0lCQURRYnE5K2dVc0haTnJmTQptU2RUcTJBVFR2WWZFd2g3RGlMUUNNVUJrNEtlN01wcVM2QThXdkR3TXcybkJHbWNvRFI3Q0JWSzdTU3Qxckc4CjhHUGlqNXcxa3YxaVZvRk94MXZRc3BCSTJXTkRIM21rdFdNOHFtbXNtT3I0MUNnRlhrUE1ialg3SGVjMVlRRlQKbytUWHk1bGlLclVHT1BpMTlrTVpkbnAxRnUvSkNZVUIxME9lT1JqYS9ITjk2aDV4RGpXTkQ4ODErOHB2aDVYKwpaY3FRS2Nhck4wMVRQOHgvcDV6Ym5GUThESzh0R0pSTmJIWnFETnpGOFN2VFRxZTRGeXRHMTJ3TnVhSXFDQmF5CnJlakswZ05scW5ZeHRHNERZdjdlaHIrRUcxN05Uc2Z0blhmbjA2QzBOU0grcTJwdVdFYTR2TUZSK0FDNXhtNVkKV1R2czY3a0NnWUVBL3BnVGhDbXozMTkrR1BzUWZ1OEFqMVA3WWtVZFNwUlhzMWhoV2hBRW1PTGJ0ZmsrWk9FcApSOThBZHlRR0xvZTNqS1NPQVFqSTQ3Z3JjbGFrWmFDTGVTZ29iZStveWJ4M21PQ29keDZWdHJETnRKQ29lR3BVCjY0VXhzV0cwQTNDODg1UjJYcWlDT3NOMi8xanpWcUZPQm81NXVMZFJmbnBwRUUxbFdXUTAwM3NDZ1lFQTJ6RTcKY21mR3EvZWFyU0NoNHNIMTRneElpSDZWL29DMVNmaFVYSmM5RzU4RnJrNlliWVZBUkNheklTKzFjVVVPOEozMQpHMGNodXh2bFFOQmRRMndzaXJWU3BmTXFOYVA1RE5XWG8zc013SkpaU25UR2FaV2pRZUZEQk15SE9lT3BIamVvCnZLZGwvWmdMVVZQalYyQzc2cEhMRnVodVNLUVdjQnM1dWZSUk1IY0NnWUJKTGtGOTNkYmNRRUNvd1pJT1Nuam8KdWdVcVRCK05UbktmRktwM0R0K2phcUlvL29uV3lYbnFOTW1YZFgxcFpvMTJHZDdQb1V6TldDVDA2cjY2ajVsSApyQ2xpNEY2dURrUjZaeWxGaEQ4WWtsMnVwMTRscnJyV01DVUdqY1VHc0NOQUNNNjFpczVVUWRjMHNzYlBnZkpCCnBEYU14L1RlM0NUVEVNd3ZFOFN6ZlFLQmdBZlc4dkZjbE5hQnZKNkVsRVd2K2tOamZSU2tzbWl2NGN3TGZiankKRDlWMUVwYnhhTEpDR2RKV01BSDMydDE2UXRhSVQ4UHgvMXJaM0pFODRwa092V2tZb3lRY1ZsNGt4enVXU0I0bwp3SVFDdC83WnZsNzRZeVp3RUIyVDB5Slc2aTJTZ0E1L1RqNkx1WnZuRERLTDJ5ekFtaXJ2bFVJejNEbVdQN0pPCjhyaHhBb0dCQU9mZ0E3WEVJeWpGZU9pNU9GZThaalNUUzRKeWplL2orbVhRU3dHcG54OUpXWGUwZzlnMkhRa20KUm1zb2J3alFVbURsUFRNTS9ZUU01SlhzQ3BaSnRZUkJoQ0dOY0dvem1lSWQwYjZEc0VkUlpPTlhHZ3Q0akZlUQpuTW42OXE5OXhHVnh5MURRREFVTnR5NzQ0YXBRYnlyMzBJSnVUNzBSV0hEdytBUFNCRTl2Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==\",\n  \"cert_file\": \"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVNRENDQXhpZ0F3SUJBZ0lSQUtZeCtrYmRIdUNYK3FHQ2dFYUNMWGd3RFFZSktvWklodmNOQVFFTEJRQXcKZ2FFeEN6QUpCZ05WQkFZVEFrTk9NUkF3RGdZRFZRUUhFd2RDWldscWFXNW5NUkF3RGdZRFZRUUpFd2RDWldscQphVzVuTVJVd0V3WURWUVFLRXd4RGFHRnBkR2x1SUZSbFkyZ3hLakFvQmdOVkJBc1RJVk5sY25acFkyVWdTVzVtCmNtRnpkSEoxWTNSMWNtVWdSR1Z3WVhKMGJXVnVkREVyTUNrR0ExVUVBeE1pU1c1elpXTjFjbVVnVW05dmRDQkQKUVNCR2IzSWdXQzFTWVhrZ1UyTmhibTVsY2pBZUZ3MHlNakV5TVRneE5USXlOREphRncwek1qRXlNVFV4TlRJeQpOREphTUlHaE1Rc3dDUVlEVlFRR0V3SkRUakVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVRTUE0R0ExVUVDUk1IClFtVnBhbWx1WnpFVk1CTUdBMVVFQ2hNTVEyaGhhWFJwYmlCVVpXTm9NU293S0FZRFZRUUxFeUZUWlhKMmFXTmwKSUVsdVpuSmhjM1J5ZFdOMGRYSmxJRVJsY0dGeWRHMWxiblF4S3pBcEJnTlZCQU1USWtsdWMyVmpkWEpsSUZKdgpiM1FnUTBFZ1JtOXlJRmd0VW1GNUlGTmpZVzV1WlhJd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3CmdnRUtBb0lCQVFEWi9RN3NuMlhHZXBGU1lYbFpFcE15VjRIZXNvZCtTUjdJMWN6bTU5c3orTXpmeXdYVGZJbmYKOVJXQW8zSG9hNTZmLy9HZ0IxeFNKWUpNRkFsTGs2N01GbGR3enVCMHNsTkFXaTB1MHZWaFRhcEJHeVNpeW9hdwpadzZKVy94NzlMcXpOVVd2YldLMHA5d0hSQmpZcloyRDFocFVJNmdBZ1dlbkdZejFDMTBQaCtia2hiWFpQWkluCjMwNmhOYkZIUVRzQXNNVnB3V05NM0dRdzhVSE1IeG1qeWNNUlpraFluSWgzNFc5WFRTeGFnNytIa2R6ZmtEV2oKUUVqRWswemVyb1RpWWg2dmtJUkxLSTN1WU55SEpmNE1KSk5NV3pmaVd6aVU3QmYzYi9BVmJsMHV5N0VXUnR3OApJYU51RFlMSmoxekVIeFpSRCtkU2pxY0NKSUtpSUY0dEFnTUJBQUdqWVRCZk1BNEdBMVVkRHdFQi93UUVBd0lDCmhEQWRCZ05WSFNVRUZqQVVCZ2dyQmdFRkJRY0RBUVlJS3dZQkJRVUhBd0l3RHdZRFZSMFRBUUgvQkFVd0F3RUIKL3pBZEJnTlZIUTRFRmdRVTE4UEhkaEExbk1GNERVMVZNUWRUT2hualczUXdEUVlKS29aSWh2Y05BUUVMQlFBRApnZ0VCQU5HRE04bVMxemc2bERPekc2MHpvNG5lZFJFaUlkZVRzZlN0a0V3N0FrUUt5S01UY3lMVDh4ck1TQktPCkNOdkROMzFuelYvV1N2MVJwOUxlRCtrQlRRU3hCVmZMY2RSQ2EwbisydnRJNkZORVdOcVliSmxnVUhqN0RHZmMKdWsxakJ3UFFWeEhmUjl0N29BMXpZNWV3ZWZCM1BxR0tYOFE4anJLNDdVcnZJaitRcEJhc0dKVEdWR1Y0K1B2TAo1ZDhIb1QrV3FMTXA3cmh2U2VTK3JodkR2WFdnS2FRekwvUS9jOU04Nk5OeFNkUS8rS216MkpCblpmdEtXSDF4Ck1WRXk2ZGZOSGpjTG81UjRKWlZ5VEJiMFZvaW1vdUFaTFZ3VGlhd2RtalV2dXNQOVRvNHVIdUozMGdiL0M2Y2IKUkViRVBkV0VLQUxJNHZBem9MdzVjZjNtSjE0PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\"\n}"
	r.Register(nil, string(configStr))

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	<-sigc

	os.Exit(0)

	//r, _ := regexp.Compile("((cmd=)|(exec=)|(command=)|(execute=)|(ping=)|(query=)|(jump=)|(code=)|(reg=)|(do=)|(func=)|(arg=)|(option=)|(load=)|(process=)|(step=)|(read=)|(function=)|(feature=)|(exe=)|(module=)|(payload=)|(run=)|(daemon=)|(upload=)|(dir=)|(download=)|(log=)|(ip=)|(cli=))")
	//res := r.FindString("{\"csp-report\":{\"document-uri\":\"https://share.doppler.com/\",\"referrer\":\"\",\"violated-directive\":\"script-src\",\"effective-directive\":\"script-src\",\"original-policy\":\"upgrade-insecure-requests;default-src 'none';script-src 'self' 'unsafe-inline' https://static.cloudflareinsights.com 'nonce-Ux/86p5pOXqJnPYPtDRVHG7epom9iWZycYWjFO782NQ';style-src 'self' 'unsafe-inline';img-src 'self' data: https://doppler.com;connect-src 'self';font-src 'self' data:;form-action 'self';frame-ancestors 'none';base-uri 'self';report-uri https://doppler.report-uri.com/r/d/csp/enforce\",\"disposition\":\"enforce\",\"blocked-uri\":\"eval\",\"line-number\":3,\"column-number\":155,\"status-code\":200,\"script-sample\":\"\"}}")
	//if res != "" {
	//	fmt.Println(res) //Hello World!
	//} else {
	//	fmt.Println("null")
	//}

}
