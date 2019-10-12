from flask_sqlalchemy import SQLAlchemy
db = SQLAlchemy()


# Initialize the database
def init_app(app):
    app.config['SQLALCHEMY_DATABASE_URI'] = 'postgresql://postgres@localhost'
    db.init_app(app)

    with app.app_context():
        # Import the models beforehand so we can create the tables
        from . import user, image, imagedata
        db.create_all()
