import smtplib
import sys
import os
import json

def send_email(subject, sender, recipient, body):
    # SMPT server config
    SMTP_SERVER = os.environ.get("SMTP_SERVER", "mail.privateemail.com")
    SMTP_PORT = int(os.environ.get("SMTP_PORT", 587))
    SMTP_USERNAME = os.environ.get("SMTP_USERNAME", "info@hospedate.app")
    SMTP_PASSWORD = os.environ.get("SMTP_PASSWORD", "xxx")
    headers = [
        f"From: {sender}",
        f"To: {recipient}",
        f"Subject: {subject}",
        "MIME-Version: 1.0",
        "Content-Type: text/html"
    ]
    headers = "\r\n".join(headers)
    message = (headers + "\r\n\r\n" + body).encode('utf8')

    try:
        # server connect and send email
        server = smtplib.SMTP(SMTP_SERVER, SMTP_PORT)
        server.starttls()
        server.login(SMTP_USERNAME, SMTP_PASSWORD)
        server.sendmail(sender, recipient, message)
        server.quit()
        print(json.dumps({}))
    except Exception as e:
        print(json.dumps({"error": f"{e}"}))
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 5:
        print("Usage: python main.py <subject> <sender> <recipient> <body>")
        sys.exit(1)
   
    _, subject, sender, recipient, body = sys.argv
    send_email(subject, sender, recipient, body)