from app import create_app

url = "api/register"
app = create_app()
test_client = app.test_client()


def test_valid_email():
    """Creates a user using a valid email and password"""
    request_body = {"email": "test@test.com", "password": "123456"}
    response = test_client.post(url, json=request_body)
    assert response.status_code == 200


def test_existing_email():
    """Returns an error if the user already exists"""
    request_body = {"email": "test@test.com", "password": "123456"}
    response = test_client.post(url, json=request_body)
    assert response.status_code == 409


def test_invalid_email():
    """Returns an error for a user when the email is invalid"""
    request_body = {"email": "test", "password": "123456"}
    response = test_client.post(url, json=request_body)
    assert response.status_code == 400


def test_invalid_password():
    """Returns an error when the password is invalid (Has less than 6 characters in this case"""
    request_body = {"email": "test@test.com", "password": "1234"}
    response = test_client.post(url, json=request_body)
    assert response.status_code == 400