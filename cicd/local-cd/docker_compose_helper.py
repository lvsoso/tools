import subprocess
from subprocess import PIPE
import os
from datetime import datetime

"""
ref: https://stackoverflow.com/questions/61973128/python-subprocess-for-docker-compose
"""

class DockercomposeRun:

    def __init__(self):
        self.root_command = "/usr/local/bin/docker-compose"

    def run(self, commands: list, env: dict):
        popen = subprocess.Popen(commands, env=env, stdin=PIPE, stdout=PIPE, stderr=PIPE, universal_newlines=True)
        return popen

    def up(self, file_path, env: dict):
        commands = [self.root_command, "-f", file_path, "up", "-d"]
        return self.run(commands, env)

    def down(self, file_path, env: dict):
        commands = [self.root_command, "-f", file_path, "down"]
        return self.run(commands, env)
    
    def force_recreate(self, file_path, env: dict):
        commands = [self.root_command, "-f", file_path, "up", "-d", "--force-recreate"]
        return self.run(commands, env)

def update_image( file_path, env:dict, logger_func):
    dcr = DockercomposeRun()
    try:
        rc = dcr.force_recreate(file_path, env)
        logger_func(rc.stdout)
        rc.stdout.close()
        logger_func(rc.stderr)
        rc.stderr.close()

        # rc = dcr.down(filename,{})
        # logger_func(rc.stdout)
        # rc.stdout.close()
        # logger_func(rc.stderr)
        # rc.stderr.close()

        # rc = dcr.up(filename, env)
        # logger_func(rc.stdout)
        # rc.stdout.close()
        # logger_func(rc.stderr)
        # rc.stderr.close()
        
        return rc.returncode
    except Exception as e:
        print(e)
        return -1

def show(output):
    for line in iter(output.readline, ""):
        print(line, end="")
        if line == "":
            break

if __name__ == "__main__":
    dir_name, _ = os.path.split(os.path.abspath(__file__))
    filename = "docker-compose.yaml"
    file_path = os.path.join(dir_name, filename)
    env = {"IMAGE": "nginx", "TAG": "1.17"}
    print(update_image(file_path, env, show))
