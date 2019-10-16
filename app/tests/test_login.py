from app import create_app

register_url = "api/register"
login_url = "api/login"
app = create_app()
test_client = app.test_client()


def test_valid_login():
    """Login a valid registered user"""
    register_request_body = {"email": "test2@test.com", "password": "123456"}
    test_client.post(register_url, json=register_request_body)
    response = test_client.post(login_url, json=register_request_body)
    assert response.status_code == 200


def test_invalid_email_login():
    """Errors when trying to login with an unregistered user"""
    request_body = {"email": "test3@test.com", "password": "123456"}
    response = test_client.post(login_url, json=request_body)
    assert response.status_code == 404


def test_invalid_password_login():
    """Errors when trying to login with a registered user with the wrong password"""
    request_body = {"email": "test2@test.com", "password": "1234567890"}
    response = test_client.post(login_url, json=request_body)
    assert response.status_code == 403
