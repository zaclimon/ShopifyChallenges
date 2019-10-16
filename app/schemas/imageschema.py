from . import ma
from marshmallow import fields


class ImageDataSchema(ma.Schema):
    """Marshmallow schema used for representing an image's metadata model"""
    image_hash = fields.String()


class ImageSchema(ma.Schema):
    """Marshmallow schema used for representing an image model."""
    id = fields.Integer()
    filename = fields.String()
    url = fields.String()
    size = fields.Integer()
    image_data = fields.Nested(ImageDataSchema)
    user_id = fields.Integer()
