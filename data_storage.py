import os


def save_users(user_ids):
    with open("user_ids", "w") as file:
        for line in user_ids:
            file.write(str(line) + "\n")


def load_users(users_id):
    if not os.path.exists("user_ids"):
        return
    with open("user_ids", "r") as file:
        lines = file.readlines()
        for line in lines:
            users_id.add(int(line))


def add_from_email(msg_id, email):
    with open("from_email", "a") as file:
        file.write(f"{msg_id}:{email}\n")


def get_from_email(msg_id):
    if not os.path.exists("from_email"):
        return
    with open("from_email", "r") as file:
        lines = file.readlines()
        for line in lines:
            a = line.split(":")
            if a[0] == msg_id:
                return a[1]
    return ""
