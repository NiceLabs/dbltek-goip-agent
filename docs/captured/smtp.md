# SMS Receive via SMTP

1. The SMTP NOT SECURE, its work in 25 port
2. Plain Login ONLY
3. Playload using base64 encoding

```plain
220  ESMTP Service Ready
EHLO sms-agent.lan
250-Hello sms-agent.lan
250-PIPELINING
250-8BITMIME
250-ENHANCEDSTATUSCODES
250 AUTH PLAIN LOGIN
AUTH LOGIN
334 VXNlcm5hbWU6
dFBQYVRrelhabUhlaEFzcw==
334 UGFzc3dvcmQ6
eUhuUWpNdlZRNnd0R2NQQw==
235 2.0.0 Authentication succeeded
MAIL FROM: <tPPaTkzXZmHehAss>
250 2.0.0 Roger, accepting mail from <tPPaTkzXZmHehAss>
RCPT TO: <sample@example.com>
250 2.0.0 I'll make sure <sample@example.com> gets this
DATA
354 2.0.0 Go ahead. End your data with <CR><LF>.<CR><LF>
From: tPPaTkzXZmHehAss
To: sample@example.com
Subject: SMS
Date: Sun, 12 Apr 2020 04:16:13+0
X-Mailer: smail 0.1
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: base64

U046OE1DRFJNMTgwNDcyNDcgQ2hhbm5lbDo4IFNlbmRlcjowNC0xMiAxMjoxNjoxMyxTbWFyVG9uZSznj77mmYLkvaDnmoTlhLLlgLzph5HpoY3ngro6JDI0Mi4wNCDmnInmlYjmnJ/oh7MyMy8wNS8yMDIwIDIzOjU5Cg==
.
250 2.0.0 OK: queued
QUIT
221 2.0.0 Goodnight and good luck
```

## Payload structure

```plain
SN:8MCDRM18047247 Channel:8 Sender:04-12 12:16:13,SmarTone,現時你的儲值金額為:$242.04 有效期至23/05/2020 23:59\n
```
