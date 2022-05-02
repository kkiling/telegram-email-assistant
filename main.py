import json
import os
import telebot
from time import sleep
from data_storage import add_from_email, get_from_email, load_users, save_users
from read_msg import get_name_and_email, print_email, print_email_body, read_email_body, read_unseen_emails
from threading import Thread

IS_PROD = os.environ.get('IS_PROD', False)
if not IS_PROD:
    from dotenv import load_dotenv
    load_dotenv(".env")

# pip install python-dotenv
server = os.environ.get('MAIL_SERVRE', 'imap.yandex.ru')
login = os.environ.get('MAIL_LOGIN', '')
password = os.environ.get('MAIL_PASSWORD', '')
bot_token = os.environ.get('TG_BOT_TOKEN', '')
max_text_message = int(os.environ.get('MAX_TEXT_MESSAGE', 256))
time_out_email_checker = int(os.environ.get('TIME_OUT_EMAIL_CHECKER', 5))
time_out_error = int(os.environ.get('TIME_OUT_ERROR', 60))
allowed_users = set(json.loads(os.environ.get('ALLOWED_USERS_ID', '[]')))
#
users_id = set()
load_users(users_id)
users_id = set(users_id) & set(allowed_users)
#
bot = telebot.TeleBot(bot_token)
run_reading_progress = {}

if not os.path.isdir("data"):
    os.mkdir("data")

def read_email():
    while True:
        emails, isError = read_unseen_emails(server, login, password)
        if isError:
            for id in users_id:
                bot.send_message(id, '❗ An error occurred while reading email')
            sleep(time_out_error)
        else:
            for id in users_id:
                for email in emails:
                    msg_id = email.msg_id.decode('utf-8')
                    callback_data = msg_id
                    markup = telebot.types.InlineKeyboardMarkup()
                    markup.add(telebot.types.InlineKeyboardButton(
                        text='Read email', callback_data=callback_data))
                    bot.send_message(id, print_email(email),
                                     parse_mode='HTML', reply_markup=markup)

                    _, from_email = get_name_and_email(email.fromEmail)
                    add_from_email(msg_id, from_email)
            sleep(time_out_email_checker)


def read_mail_time(chatId, msg_id):
    fromEmail = get_from_email(msg_id.decode('utf-8'))
    message = bot.send_message(chatId, f'⌛ Reading a mail from {fromEmail}')
    index = 0
    seconds = 0
    while msg_id in run_reading_progress and run_reading_progress[msg_id] == True:
        sleep(0.1)
        index += 1
        text = ''
        if index % 10 == 0:
            seconds += 1
            text = f'⏳ Reading a mail from {fromEmail} ({seconds} sec)'
        if index % 20 == 0:
            seconds += 1
            text = f'⌛ Reading a mail from {fromEmail} ({seconds} sec)'
        if text != "":
            bot.edit_message_text(
                chat_id=chatId, message_id=message.id, text=text)
    bot.delete_message(chat_id=chatId, message_id=message.id)


@bot.message_handler(commands=["start"])
def start(m, res=False):
    if not m.chat.id in allowed_users:
        bot.send_message(
            m.chat.id, f'❗ Access is denied: your id #{m.chat.id}')
        return

    users_id.add(m.chat.id)
    save_users(users_id)

    msg = f'✌ Hey! I am your personal email assistant.\n'
    msg += f'📧 I will send notifications of new email in your mailbox: {login}'
    bot.send_message(m.chat.id, msg)


@bot.message_handler(func=lambda message: True)
def echo_message(m):
    bot.send_message(
        m.chat.id, '🚫 I am not trained to respond to messages or commands')


@bot.callback_query_handler(func=lambda call: True)
def query_handler(call):
    # Прерываем чтение всех остальных потоков...
    global run_reading_progress
    run_reading_progress = {}

    id = call.message.chat.id
    msg_id = call.data.encode('utf-8')

    # Progress
    run_reading_progress = {msg_id: True}
    thread = Thread(target=read_mail_time, args=(id, msg_id,))
    thread.start()
    # Read
    result, isError = read_email_body(server, login, password, msg_id)

    if msg_id in run_reading_progress and run_reading_progress[msg_id] == True:

        if isError:
            bot.send_message(id, '❗ An error occurred while reading email')
            return
        text, img, attachment = print_email_body(result, max_text_message)

        if img != None:
            img_file = open(img, 'rb')
            bot.send_photo(id, img_file, caption=text, parse_mode='HTML')
        else:
            bot.send_message(id, text, parse_mode='HTML')

        for at in attachment:
            doc_file = open(at, 'rb')
            bot.send_document(id, doc_file)

    run_reading_progress[msg_id] = False


#
thread = Thread(target=read_email)
thread.start()
#
bot.infinity_polling()