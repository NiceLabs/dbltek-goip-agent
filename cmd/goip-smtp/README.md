# GoIP SMTP Agent

## Configuration

Filename: `configure.json`

```jsonc
{
    "smtp_username": "", // SMTP Server Username
    "smtp_password": "", // SMTP Server Password
    "hook": "send-message-program.py", // Local Hook program file
    "phones": {
        "[Device - Serial Number]": {
            "[Channel Number]": "[Phone number]"
        }
    }
}
```

## Program file - Standard input

\* **all fields is string type**

```jsonc
{
    "to": "your GoIP settings send target address",

    "sn": "the device serial number",
    "channel": "the channel number in device",

    "phone": "in configure file",
    "date": "mm-dd HH:MM:SS",
    "sender": "the sender phone number",
    "message": "sms message content"
}
```

## [example] File tree

```plain
/opt/goip-smtp/
├── goip-smtp
├── configure.json
└── send-via-mailgun.py
```

## [example] Systemd service

```ini
[Unit]
Description=GoIP SMTP Agent

[Service]
Type=simple
ExecStart=/opt/goip-smtp/goip-smtp
Restart=always
StandardOutput=file:/var/log/goip-smtp.log

[Install]
WantedBy=multi-user.target
```
