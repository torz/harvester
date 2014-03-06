from flask import render_template, jsonify
from Processor import app

@app.route('/')
def index():
	msg = ['Web spider', 'Processor', 'Pashbe&Szabo Corp']
	return render_template('index.html', title='Processor', message=msg)

@app.route('/urls')
def urls_to_get():
	harvest = {'harvest': [{'id': 0, 'url': 'http://bbc.co.uk'},
							{'id': 1, 'url': 'http://www.imdb.com/'},
							{'id': 2, 'url': 'http://distrowatch.com/'}]}
	return jsonify(harvest)

