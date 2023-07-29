from flask import Flask, request
from prometheus_flask_exporter import PrometheusMetrics

app = Flask(__name__)

@app.route("/hello")
def hello():
    return "hi"

@app.route("/hi")
def hi():
    return "hello"

def get_outcome(resp):
    if resp.status_code >=400:
        return 'fail'
    return 'success'

def get_status(resp):
    return resp.status_code

metrics = PrometheusMetrics(app)

metrics.register_default(
    metrics.histogram(
        'http_server_requests_seconds', 
        'How long it took to process the request, partitioned by status code and HTTP method.',
        labels={
            "application":"py-app", 
			"status":get_status, 
			"uri":lambda:request.path, 
			"method":lambda:request.method,
            "outcome":get_outcome
		}
    )
)
