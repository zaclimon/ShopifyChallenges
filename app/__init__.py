from flask import Flask
from flask_bcrypt import Bcrypt
from dotenv import load_dotenv
import os, warnings

flask_bcrypt = Bcrypt()


def create_app():
    from . import models, routes, schemas
    load_dotenv()
    verify_env_variables()
    app = Flask(__name__)
    set_app_variables(app)
    flask_bcrypt.init_app(app)
    models.init_app(app)
    routes.init_app(app)
    schemas.init_app(app)
    return app


def set_app_variables(app):
    app.config["SQLALCHEMY_DATABASE_URI"] = os.getenv("SQLALCHEMY_DATABASE_URI")
    app.config["SECRET_KEY"] = os.getenv("SECRET_KEY")
    app.config["UPLOAD_FOLDER"] = os.getenv("UPLOAD_FOLDER")


def verify_env_variables():
    if os.getenv("SQLALCHEMY_DATABASE_URI") is None:
        raise ValueError("Please set SQLALCHEMY_DATA_BASE in your environment variables")

    if os.getenv("SECRET_KEY") is None:
        raise ValueError("Please set SECRET_KEY in your environment variables")

    if os.getenv("GOOGLE_APPLICATION_CREDENTIALS") is None:
        warnings.warn("GOOGLE_APPLICATION_CREDENTIALS variable is not set. Upload will not be possible!")

    if os.getenv("UPLOAD_FOLDER") is None:
        warnings.warn("UPLOAD_FOLDER variable is not set. Upload will not be possible!")
