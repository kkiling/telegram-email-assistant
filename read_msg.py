import imgkit
import os
import re
import email
import imaplib
import traceback
import pathlib
from threading import Thread
from PIL import Image
from dateutil import parser
from email.header import decode_header, make_header


class MsgInfo:
    msg_id = ""
    fromEmail = ""
    toEmail = ""
    subject = ""
    date = ""
    body = []

    def __init__(self, msg_id, fromEmail, toEmail, subject, date):
        self.msg_id = msg_id
        self.fromEmail = fromEmail
        self.toEmail = toEmail
        self.subject = subject
        self.date = date


class MsgBody:
    msg_id = ""
    contentType = ""
    body = ""

    def __init__(self, msg_id, contentType, body):
        self.msg_id = msg_id
        self.contentType = contentType
        self.body = body


class MsgAttachment:
    msg_id = ""
    contentType = ""
    fileName = ""
    filePath = ""

    def __init__(self, msg_id, contentType, fileName, filePath):
        self.msg_id = msg_id
        self.contentType = contentType
        self.fileName = fileName
        self.filePath = filePath


def _msg_subject(msg):
    return str(make_header(decode_header(msg["Subject"])))


def _msg_date(msg):
    return str(make_header(decode_header(msg["Date"])))


def _msg_from(msg):
    return str(make_header(decode_header(msg["From"])))


def _msg_to(msg):
    return str(make_header(decode_header(msg["To"])))


def _regexp_getvalue(value, pattern):
    result = ""
    matches = re.finditer(pattern, value, re.MULTILINE)
    for _, match in enumerate(matches, start=1):
        if len(match.groups()) > 0:
            result = match.groups()[0]
    return result.strip()


def get_name_and_email(value):
    name = _regexp_getvalue(value, r"^(.*?)\<")
    email = _regexp_getvalue(value, r"\<(.*?)\>")
    if email == "":
        email = value
    return name, email


def html_to_png(msg_id, text_html):
    if text_html == "":
        return None

    msg_folder = f"data/{msg_id.decode('utf-8')}"
    if not os.path.isdir(msg_folder):
        os.mkdir(msg_folder)

    id = msg_id.decode('utf-8')
    filename = f'{msg_folder}/{id}.png'
    
    if os.path.exists(filename):
        return filename
        
    try:
        src_folder = pathlib.Path(__file__).parent.resolve()
        src_folder = os.path.join(src_folder, msg_folder)
        have_cid = 'src="cid:' in text_html
        text_html = text_html.replace('src="cid:', f'src="{src_folder}/')

        #html_filename = f'{msg_folder}/index.html'
        #with open(html_filename, "w", encoding='utf-8') as text_file:
        #    text_file.write(text_html)

        # Try with local files
        if have_cid:
            options = {
                'format': 'png',
                'enable-local-file-access': None,
            }
            thread = Thread(target=imgkit.from_string, args=(text_html, filename, options,))
            thread.start()
            thread.join(60)

        # Try without local files
        if not os.path.exists(filename):
            thread = Thread(target=imgkit.from_string, args=(text_html, filename))
            thread.start()
            thread.join(60)

        if os.path.exists(filename):
            picture = Image.open(filename)
            picture.save(filename)
    except:
        pass
    if os.path.exists(filename):
        return filename
    return None


def print_email_body(email: MsgInfo, max_text_message):
    text_plain = ""
    text_html = ""
    for msg in email.body:
        if type(msg) is MsgBody:
            if msg.contentType == "text/plain":
                text_plain = msg.body
            elif msg.contentType == "text/html":
                text_html = msg.body
            else:
                return f"❗ Undefined  msg content type: {msg.contentType}"
    img = None
    text = print_email(email) + "\n\n"

    #save_html_as_image = text_html != "" and (len(text_plain) > max_text_message or text_plain == "")
    #if save_html_as_image or 'src="cid:' in text_html:
    if text_html != "":
        img = html_to_png(email.msg_id, text_html)
    else:
        text += text_plain

    # Attachment
    attachment = []
    for msg in email.body:
        if type(msg) is MsgAttachment:
            text += f"\n📎 {msg.fileName}"
            attachment.append(msg.filePath)

    return text, img, attachment


def print_email(email: MsgInfo):
    to_name, to_email = get_name_and_email(email.toEmail)
    from_name, from_email = get_name_and_email(email.fromEmail)

    result = ""
    if to_name != "":
        result += f"<b>📫 {to_name}</b>\t"
        result += f"({to_email})\n\n"
    else:
        result += f"<b>📫 {to_email}</b>\n\n"

    if from_name != "":
        result += f"<b>📨 {from_name}</b>\t"
        result += f"({from_email})\n\n"
    else:
        result += f"<b>📨 {from_email}</b>\n\n"
    date_time = parser.parse(email.date)
    d = date_time.strftime("%d %B %Y, %H:%M")
    result += f"⏰ <b>{d}</b>\n\n"
    result += f"📝 <b>{email.subject}</b>"
    return result


def read_unseen_emails(server, login, password):
    result = []

    imap = imaplib.IMAP4_SSL(server)
    imap.login(login, password)
    try:
        status, _ = imap.select('INBOX')
        if status != "OK":
            print("Error when selecting INBOX")
            return [], True

        # nmessages = select_data[0].decode('utf-8')
        status, search_data = imap.search(None, 'UnSeen')
        if status != "OK":
            print("Error when reading UnSeen messages")
            return [], True

        for msg_id in search_data[0].split():
            # msg_id_str = msg_id.decode('utf-8')
            # print("Fetching message {} of {}".format(msg_id_str, nmessages, policy=default_policy))
            _, msg_data = imap.fetch(msg_id, '(RFC822)')
            msg_raw = msg_data[0][1]
            msg = email.message_from_bytes(
                msg_raw, _class=email.message.EmailMessage)

            from_ = _msg_from(msg)
            to = _msg_to(msg)
            subject = _msg_subject(msg)
            date = _msg_date(msg)

            result.append(MsgInfo(msg_id, from_, to, subject, date))

            print("From: ", from_)
            print("To: ", to)
            print("Subject: ", subject)
            print("Date: ", date)
            print("="*100)

        return result, False
    except Exception:
        print("Error while reading messages:")
        print(traceback.format_exc())
        return [], True
    finally:
        imap.close()
        imap.logout()
        
def _save_attachment_file(msg_folder, part):
    filename, encoding = decode_header(part.get_filename())[0]
    if(encoding is None):
        open(msg_folder + filename, 'wb').write(part.get_payload(decode=True))
        return filename
    else:
        filename = filename.decode(encoding)
        open(msg_folder + filename, 'wb').write(part.get_payload(decode=True))
        return filename

def _save_inline_file(msg_folder, part):
    content_id = str(part.get("Content-ID"))
    content_id = _regexp_getvalue(content_id, r"\<(.*?)\>")
    #content_id = f"cid:{content_id}"
    open(msg_folder + content_id, 'wb').write(part.get_payload(decode=True))

def read_email_body(server, login, password, msg_id):
    msg_folder = f"data/{msg_id.decode('utf-8')}/"
    if not os.path.isdir(msg_folder):
        os.mkdir(msg_folder)

    imap = imaplib.IMAP4_SSL(server)
    imap.login(login, password)
    try:
        status, _ = imap.select('INBOX')
        if status != "OK":
            print("Error when selecting INBOX")
            return None, True

        _, msg_data = imap.fetch(msg_id, '(RFC822)')
        msg_raw = msg_data[0][1]
        msg = email.message_from_bytes(
            msg_raw, _class=email.message.EmailMessage)

        from_ = _msg_from(msg)
        to = _msg_to(msg)
        subject = _msg_subject(msg)
        date = _msg_date(msg)
        result = MsgInfo(msg_id, from_, to, subject, date)

        if msg.is_multipart():
            result.body = []
            for part in msg.walk():
                content_type = part.get_content_type()
                content_disposition = part.get_content_disposition()
                
                body = None
                try:
                    body = part.get_payload(decode=True).decode()
                except:
                    pass
                part.get_content_disposition
                if content_disposition != None and "attachment" in content_disposition:

                    filename = _save_attachment_file(msg_folder, part)
                    result.body.append(MsgAttachment(msg_id, content_type, filename, msg_folder + filename))

                elif content_disposition != None and "inline" in content_disposition:
                    _save_inline_file(msg_folder, part)

                    # Save inline file as attachment
                    filename = _save_attachment_file(msg_folder, part)
                    result.body.append(MsgAttachment(msg_id, content_type, filename, msg_folder + filename))

                elif body != None:
                    result.body.append(MsgBody(msg_id, content_type, body))
        else:
            content_type = msg.get_content_type()
            body = msg.get_payload(decode=True).decode()
            result.body.append(MsgBody(msg_id, content_type, body))

        return result, False
    except Exception:
        print("Error while reading messages:")
        print(traceback.format_exc())
        return None, True
    finally:
        imap.close()
        imap.logout()
