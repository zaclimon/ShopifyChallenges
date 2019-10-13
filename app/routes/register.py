from flask import Blueprint, request, jsonify
from app import flask_bcrypt
from app.schemas.userschema import UserSchema
from app.models import db
from app.models.user import User

register_bp = Blueprint("register", __name__)


@register_bp.route("/register", methods=("GET", "POST"))
def register():

    if request.method == "POST":
        # Inspired from: https://flask.palletsprojects.com/en/1.1.x/tutorial/views/#the-first-view-register
        form_email = request.form["email"]
        form_password = request.form["password"]
        user_schema = UserSchema()

        if not form_email:
            error = "An email is required"
        elif not form_password:
            error = "A password is required"
        elif User.query.filter_by(email=form_email).first() is not None:
            error = "The email {} is already registered!".format(form_email)
        else:
            # Encode the password so it can be readable.
            # See: https://github.com/maxcountryman/flask-bcrypt/issues/38#issuecomment-247513357
            hashed_password = flask_bcrypt.generate_password_hash(form_password).decode('utf-8')
            new_user = User(email=form_email, password=hashed_password)
            db.session.add(new_user)
            db.session.commit()
            return user_schema.dump(new_user)

        error_data = {"error": error}
        return jsonify(error_data)
    else:
        return "Hi! 👋"
