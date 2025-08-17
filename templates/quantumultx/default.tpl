[general]
network_check_url=http://www.apple.com/
server_check_url=http://www.google.com/generate_204

[server_local]
{{- range .Proxies }}
{{ formatProxy . "quantumultx" }}
{{- end }}

[filter_local]
host-suffix, local, direct
geoip, cn, direct
final, proxy

[rewrite_local]

[task_local]

[http_backend]

[mitm]
