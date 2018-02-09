import datetime
import os
import os.path

from flask import Flask, jsonify

app = Flask(__name__)

LAST_RESTART = datetime.datetime.now(datetime.timezone.utc)


@app.route('/', defaults={'path': ''})
@app.route('/<path:path>')
def index(path):
    file_list = []
    for root, dirs, files in os.walk('.'):
        for fname in files:
            fpath = os.path.join(root, fname)
            file_list.append({
                "name": fpath[2:],
                "mtime": os.stat(fpath).st_mtime,
            })

    return jsonify({
        "restart": LAST_RESTART,
        "pod": os.environ.get('POD_NAME'),
        "files": file_list,
    })


@app.route('/demo')
def demo():
    return "demo - v1"
