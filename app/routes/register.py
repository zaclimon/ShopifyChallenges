from flask import Blueprint

register_bp = Blueprint("register", __name__)


@register_bp.route("/register")
def register():
    return "Hello World!"
