from flask import Blueprint, request, jsonify
from marshmallow import ValidationError
from app import flask_bcrypt
from app.schemas.userschema import UserSchema
from app.models import db
from app.models.user import User

register_bp = Blueprint("register", __name__)


@register_bp.route("/register", methods=("GET", "POST"))
def register():
    """Serves as an endpoint so the a person can register him/herself to upload pictures."""
    if request.method == "POST":
        user_schema = UserSchema()

        try:
            # Encode the password so it can be easier to store in the DB.
            # See: https://github.com/maxcountryman/flask-bcrypt/issues/38#issuecomment-247513357
            user_schema.load(request.json)

            json_email = request.json["email"]
            json_password = request.json["password"]

            if User.query.filter_by(email=json_email).first() is not None:
                return jsonify(error="The email {} is already used".format(json_email)), 409

            hashed_password = flask_bcrypt.generate_password_hash(json_password).decode('utf-8')
            new_user = User(email=json_email, password=hashed_password)
            db.session.add(new_user)
            db.session.commit()
            return user_schema.dump(new_user)
        except ValidationError as err:
            return jsonify(error=err.messages), 400
    else:
        return "Hi! 👋"
