from flask import Flask
from flask_bcrypt import Bcrypt

flask_bcrypt = Bcrypt()


def create_app():
    from . import models, routes, schemas
    app = Flask(__name__)
    app.config['JSON_SORT_KEYS'] = False
    flask_bcrypt.init_app(app)
    models.init_app(app)
    routes.init_app(app)
    schemas.init_app(app)
    return app
