from . import db
import datetime


class User(db.Model):
    """Model used for users"""
    __tablename__ = "users"

    id = db.Column(db.Integer, primary_key=True)
    # See: https://stackoverflow.com/a/1199245
    email = db.Column(db.String(254), unique=True)
    # See: https://security.stackexchange.com/a/39851
    password = db.Column(db.String(72))
    date_created = db.Column(db.DateTime, default=datetime.datetime.utcnow())
    images = db.relationship("Image", backref="user", lazy=True)
