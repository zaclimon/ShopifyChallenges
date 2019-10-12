from flask import Blueprint, flash, g, redirect, render_template, request, session, url_for
from app.models import db
from app.models.user import User
from app import flask_bcrypt

register_bp = Blueprint("register", __name__)


@register_bp.route("/register", methods=("GET", "POST"))
def register():

    if request.method == "POST":
        # Inspired from: https://flask.palletsprojects.com/en/1.1.x/tutorial/views/#the-first-view-register
        form_email = request.form["email"]
        form_password = request.form["password"]
        error = None

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
            return "Account for {} has been created".format(form_email)
    else:
        return "Hi! 👋"
