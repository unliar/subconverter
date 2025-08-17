port: 7890
socks-port: 7891
allow-lan: false
mode: rule
log-level: info
external-controller: 127.0.0.1:9090

proxies:
{{- range .Proxies }}
  - name: "{{ .Remark }}"
    type: {{ .Type.String }}
    server: {{ .Hostname }}
    port: {{ .Port }}
    {{- if eq .Type.String "ss" }}
    cipher: {{ .EncryptMethod }}
    password: {{ .Password }}
    {{- else if eq .Type.String "vmess" }}
    uuid: {{ .UserID }}
    alterId: {{ .AlterID }}
    cipher: {{ default "auto" .EncryptMethod }}
    {{- else if eq .Type.String "trojan" }}
    password: {{ .Password }}
    {{- end }}
{{- end }}

proxy-groups:
  - name: "PROXY"
    type: select
    proxies:
      - "Auto"
{{- range .Proxies }}
      - "{{ .Remark }}"
{{- end }}
  
  - name: "Auto"
    type: url-test
    proxies:
{{- range .Proxies }}
      - "{{ .Remark }}"
{{- end }}
    url: 'http://www.gstatic.com/generate_204'
    interval: 300

rules:
  - DOMAIN-SUFFIX,local,DIRECT
  - IP-CIDR,127.0.0.0/8,DIRECT
  - IP-CIDR,172.16.0.0/12,DIRECT
  - IP-CIDR,192.168.0.0/16,DIRECT
  - IP-CIDR,10.0.0.0/8,DIRECT
  - GEOIP,CN,DIRECT
  - MATCH,PROXY
