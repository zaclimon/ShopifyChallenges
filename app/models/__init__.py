from flask_sqlalchemy import SQLAlchemy
db = SQLAlchemy()


# Initialize the database
def init_app(app):
    """Initialize the models for the app

    :param app The Flask instance
    """
    db.init_app(app)

    with app.app_context():
        # Import the models beforehand so we can create the tables
        from . import user, image, imagedata
        db.create_all()
