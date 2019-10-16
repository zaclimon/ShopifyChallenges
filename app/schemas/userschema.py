from . import ma
from marshmallow import fields, validate


class UserSchema(ma.Schema):
    """Marshmallow schema used for representing and validating a user model."""
    email = fields.Email(required=True)
    password = fields.String(required=True, load_only=True, validate=validate.Length(min=6))
    date_created = fields.DateTime(dump_only=True)
