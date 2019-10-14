from . import ma
from marshmallow import fields


class ImageDataSchema(ma.Schema):
    image_hash = fields.String()


class ImageSchema(ma.Schema):
    id = fields.Integer()
    filename = fields.String()
    url = fields.String()
    size = fields.Integer()
    image_data = fields.Nested(ImageDataSchema)
    user_id = fields.Integer()
