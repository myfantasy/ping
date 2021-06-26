# ping

Call http service. When not 200 ok then sent to telegram.

## ping.settings.json
``` json
{
    "token":"1 YOUR SUPER SECRET TOKEN 0",
    "chats":[121,1222],
    "urls":[{
		"url": "https://yandex.ru/",
		"name": "yandex"
		},
		{
		"url": "http://localhost:8676/ping",
		"name": "capserver"
		}
	]
}
```
## How to get chat id
1. run ping
1. Write something to bot
1. get your chat id from console
