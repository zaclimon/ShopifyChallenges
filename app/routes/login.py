from flask import Blueprint, request, jsonify, current_app as app
from marshmallow import ValidationError
from itsdangerous import JSONWebSignatureSerializer
from app import flask_bcrypt
from app.models.user import User
from app.schemas.userschema import UserSchema

login_bp = Blueprint("login", __name__)


@login_bp.route("/login", methods=("GET", "POST"))
def login():
    """Serves as an endpoint for authenticating a user so he/she can upload or search pictures."""
    if request.method == "POST":
        user_schema = UserSchema()

        try:
            user_json = user_schema.load(request.json)
            json_password = user_json["password"]
            saved_user = User.query.filter_by(email=user_json["email"]).first()

            if saved_user is None:
                return jsonify(error="User not found"), 404
            elif not flask_bcrypt.check_password_hash(saved_user.password, json_password):
                return jsonify("Invalid password"), 403

            # Create token for user
            serializer = JSONWebSignatureSerializer(app.config["SECRET_KEY"])
            token = serializer.dumps({'id': saved_user.id})
            return jsonify(access_token=token.decode("utf-8"))

        except ValidationError as err:
            return jsonify(error=err.messages), 400
