from .register import register_bp
from .login import login_bp
from .upload import upload_bp
from .search import search_bp


def init_app(app):
    app.register_blueprint(register_bp, url_prefix="/api")
    app.register_blueprint(login_bp, url_prefix="/api")
    app.register_blueprint(upload_bp, url_prefix="/api")
    app.register_blueprint(search_bp, url_prefix="/api")
