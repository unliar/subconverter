[General]
loglevel = notify
bypass-system = true
skip-proxy = 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12, localhost, *.local, captive.apple.com
bypass-tun = 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12
dns-server = system

[Proxy]
DIRECT = direct
{{- range .Proxies }}
{{ formatProxy . "surge" }}
{{- end }}

[Proxy Group]
PROXY = select, DIRECT{{range .Proxies}}, {{ .Remark }}{{end}}

[Rule]
DOMAIN-SUFFIX,local,DIRECT
GEOIP,CN,DIRECT
FINAL,PROXY
