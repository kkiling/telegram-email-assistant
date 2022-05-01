# TelegramEmailAssistant

# Build docker image

```
docker build -t kiling91/telegram-email-assistant:latest .
docker push kiling91/telegram-email-assistant:latest
```

# Run docker

```
docker run -d --name telegram-email-assistant \
	-e MAIL_SERVRE='imap.yandex.ru' \
	-e MAIL_LOGIN='' \
    -e MAIL_PASSWORD='' \
    -e TG_BOT_TOKEN='' \
    -e ALLOWED_USERS_ID='[]' \
	kiling91/telegram-email-assistant:latest
```
