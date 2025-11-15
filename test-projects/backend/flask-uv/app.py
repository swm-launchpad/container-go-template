"""
Flask test application with uv dependency manager
"""
from flask import Flask, jsonify

app = Flask(__name__)


@app.route('/')
def hello():
    return jsonify({
        'message': 'Hello from Flask with uv!',
        'dependency_manager': 'uv',
        'framework': 'Flask'
    })


@app.route('/health')
def health():
    return jsonify({'status': 'healthy'})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
