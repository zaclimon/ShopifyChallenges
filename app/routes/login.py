from flask import Blueprint

login_bp = Blueprint("login", __name__)


@login_bp.route("/login")
def register():
    return "Hey there!"
