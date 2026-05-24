[Unit]
Description=api_{{.AppName}}

[Service]
Type=simple
WorkingDirectory=/opt/{{.AppName}}/configs/
ExecStart=/usr/local/bin/{{.AppName}} -runtype=service -config=/opt/{{.AppName}}/configs/local/config.yaml -log=/var/log/{{.AppName}}/api.log

[Install]
WantedBy=multy-user.target