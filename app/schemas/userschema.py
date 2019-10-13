from . import ma
from marshmallow import fields, validate


class UserSchema(ma.Schema):
    email = fields.Email(required=True)
    password = fields.String(required=True, load_only=True, validate=validate.Length(min=6))
    date_created = fields.DateTime(dump_only=True)
