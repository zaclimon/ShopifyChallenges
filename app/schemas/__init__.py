from flask_marshmallow import Marshmallow

ma = Marshmallow()


def init_app(app):
    """Initialize the schemas

    :param app The Flask instance
    """
    ma.init_app(app)
