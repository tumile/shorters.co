[Unit]
Description=Shorters
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=1
StartLimitBurst=5
StartLimitIntervalSec=10
User=root
ExecStart=_exe_

[Install]
WantedBy=multi-user.target
