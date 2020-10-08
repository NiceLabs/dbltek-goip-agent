#!/usr/bin/env python3
import json
import sys

import requests

API_KEY = "your mailgun api key"
DOMAIN = "no-reply.your-domain.tld"
SENDER = "sms-notify"
payload = json.load(sys.stdin)

requests.post(
    "https://api.mailgun.net/v3/%s/messages" % DOMAIN,
    auth=("api", API_KEY),
    data={
        "from": "%s <%s@%s>" % (payload["phone"], SENDER, DOMAIN),
        "to": payload["to"],
        "subject": payload["sender"],
        "text": payload["message"].replace("\\n\\n", "\n").replace("\n ", "\n"),
    },
)
