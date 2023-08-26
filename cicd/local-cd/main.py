import json
import sys
import os
import base64

from docker_compose_helper import update_image
from http.server import HTTPServer, BaseHTTPRequestHandler


DIR_NAME, _ = os.path.split(os.path.abspath(__file__))
FILE_NAME = "docker-compose.yaml"
USERNAME = "root"
PASSWORD = ""

TOKEN = base64.b64encode((USERNAME + ":" + PASSWORD).encode("utf-8")).decode("utf-8")

class Handler(BaseHTTPRequestHandler):

    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'Hello, world!')

    def do_POST(self):
        auth = self.headers['Authorization']
        if auth != "Basic " + TOKEN:
            self.send_response(401)
            self.end_headers()
            self.wfile.write(b'Unauthorized')
            return
        
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        image_info = json.loads(post_data)
        print(image_info)
        if "image" not in image_info or "tag" not in image_info:
            self.send_response(404)
            self.end_headers()
            self.wfile.write(("check data: %s" % post_data).encode("utf-8"))
            return
        
        # exec update
        env = {"IMAGE": image_info["image"], "TAG": image_info["tag"]}
        def show(output):
            for line in iter(output.readline, ""):
                self.wfile.write(line.encode("utf-8"))
                if line == "":
                    break

        code = update_image(DIR_NAME, FILE_NAME, env, show)

        # return result
        print(code)
        if code:
            self.send_response(500)
        else:
            self.send_response(200)
        self.end_headers()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 main.py <port>")
        exit(1)
    port = int(sys.argv[1])
    httpd = HTTPServer(('localhost', port), Handler)
    httpd.serve_forever()
