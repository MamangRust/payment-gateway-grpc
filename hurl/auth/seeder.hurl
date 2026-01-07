POST http://localhost:5000/api/auth/register
Content-Type: application/json

{
    "firstname": "John",
    "lastname": "Doe",
    "email": "john.doe@example.com",
    "password": "securepassword",
    "confirm_password": "securepassword"
}
