FROM python:3.9

ENV IS_PROD true
ENV MAIL_SERVRE 'imap.yandex.ru'
ENV MAIL_LOGIN ''
ENV MAIL_PASSWORD ''
ENV TG_BOT_TOKEN ''
ENV MAX_TEXT_MESSAGE 256
ENV TIME_OUT_EMAIL_CHECKER 5
ENV TIME_OUT_ERROR 60
ENV ALLOWED_USERS []

WORKDIR /code
ENV IS_PROD Yes
RUN apt-get update
RUN apt-get install -y wkhtmltopdf
COPY ./requirements.txt .
RUN pip install -r requirements.txt
COPY . .
CMD [ "python", "./main.py" ]

