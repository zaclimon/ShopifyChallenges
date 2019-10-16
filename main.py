from app import create_app

# Initialize the app using an Application Factory function.
app = create_app()
# Push the application context so we can deal with the database later on
app.app_context().push()
