#!/usr/bin/env python3

import argparse
import paramiko
import yaml

def run_and_print(ssh, command):
    stdin, stdout, stderr = ssh.exec_command(command)
    print(stdout.read().decode())
    print(stderr.read().decode())


parser = argparse.ArgumentParser(description='Remote deploy')

# Add the hostname, username, and command arguments
parser.add_argument('yard_version')
parser.add_argument('backyard_version')
parser.add_argument('enigmas_key')
args = parser.parse_args()

YARD_VERSION = args.yard_version
BACKYARD_VERSION = args.backyard_version
ENIGMAS_KEY = args.enigmas_key

# Load the YAML file
with open("conf.yml", "r") as file:
    config = yaml.safe_load(file)

pkey_file_path = config["private_key_file_path"]
hostname = config["host"]
username = config["username"]

# Load the private key file
private_key = paramiko.RSAKey.from_private_key_file(pkey_file_path)

# Create an SSH client object
ssh = paramiko.SSHClient()

# Automatically add the server's host key (this is insecure and should only be used for testing)
ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

# Connect to the remote server using the private key
print(f'Connecting to {hostname}')
ssh.connect(hostname=hostname, username=username, pkey=private_key)

#run_and_print(ssh, "cd /opt/basement; ./deploy.sh ")
run_and_print(ssh, f"cd /opt/basement; ./deploy.sh {YARD_VERSION} {BACKYARD_VERSION} \"{ENIGMAS_KEY}\"")

# Close the connection
ssh.close()

