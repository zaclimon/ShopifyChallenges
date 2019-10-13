from . import ma


class UserSchema(ma.Schema):
    class Meta:
        fields = ("email", "date_created")
        ordered = True
