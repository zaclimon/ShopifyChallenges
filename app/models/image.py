from . import db


class Image(db.Model):

    __tablename__ = "images"

    id = db.Column(db.Integer, primary_key=True)
    # See: https://stackoverflow.com/a/265782
    filename = db.Column(db.String(255), nullable=False)
    # See: https://stackoverflow.com/a/417184
    url = db.Column(db.String(2000), nullable=False)
    # Size in Bytes
    size = db.Column(db.BigInteger)
    user_id = db.Column(db.Integer, db.ForeignKey("users.id"))
    image_data_id = db.relationship("Image", backref="image", lazy=True, uselist=False)
