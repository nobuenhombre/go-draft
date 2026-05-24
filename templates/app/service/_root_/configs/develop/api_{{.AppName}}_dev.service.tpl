[Unit]
Description=api_{{.AppName}}_dev

[Service]
Type=simple
WorkingDirectory=/opt/{{.AppName}}.dev/configs/
ExecStart=/usr/local/bin/{{.AppName}}-dev -runtype=service -config=/opt/{{.AppName}}.dev/configs/develop/config.yaml -log=/var/log/{{.AppName}}/api.dev.log

[Install]
WantedBy=multy-user.target