from flask import Flask
from . import models, routes


def create_app():
    app = Flask(__name__)
    models.init_app(app)
    routes.init_app(app)
    return app
