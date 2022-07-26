# TelegramEmailAssistant

Telegram bot that sends notifications of new emails. With the ability to read the mail
.
<img src="https://raw.githubusercontent.com/kiling91/telegram-email-assistant/main/preview.png" height="640" />

# Configuration

The configuration file template is located along the path /configs/config.template.yml

| Parameter                 |                      Description                       |
| ------------------------- | :----------------------------------------------------: |
| app.file_directory        |      Path for saving emails and attachment files       |
| app.store_db              |                  Path to sqlite file                   |
| app.max_text_message_size |  The maximum size of the mail message output as text   |
| app.mail_check_timeout    |        Interval for checking new messages (sec)        |
| imap.login                |                   Imap server login                    |
| imap.imap_server          |                  Imap server address                   |
| imap.password             |                  Password imap server                  |
| telegram.bot_token        |           Telegram bot token ( @BotFather )            |
| telegram.users.user_id    |                   Telegram bot user                    |
| telegram.users.imap_login | Mailbox that is associated with telegram.users.user_id |

# Build docker image

```
docker build -t kiling91/telegram-email-assistant:0.1.20 .
docker push kiling91/telegram-email-assistant:0.1.20
```

# Run docker

```
docker-compose up
```
