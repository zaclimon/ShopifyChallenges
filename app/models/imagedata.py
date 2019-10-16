from . import db


class ImageData(db.Model):
    """Model class for an Image metadata"""
    __tablename__ = "imagedata"

    id = db.Column(db.Integer, primary_key=True)
    image_hash = db.Column(db.String(32), nullable=False)
    image_id = db.Column(db.Integer, db.ForeignKey("images.id"))
