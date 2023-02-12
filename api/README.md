# API Documentation
## Endpoints
### `POST api/backend/register_device`
Đăng kí thiết bị tới backend và nhận lại password tạo bởi backend. Password này sẽ được sử dụng để đăng nhập vào web application.
IoT Gateway sẽ thực hiện HTTP request tới endpoint khi ở lần chạy đầu tiên. Ở phía IoT Gateway, password khi nhận lại sẽ được in ra console.
#### **Header**
Yêu cầu header:
1. `Content-Type: application/json`
#### **Body**
```
{
	"device_id": "ID_that_Gateway_creates_on_first_run"
}
```
#### **Response**
Password chỉ được trả lại vào request đầu tiên, các lần tiếp theo sẽ trả lại status code `409 Conflict`.
```
{
	"password": "raw_password_no_hash",
}
```
### `POST api/backend/login`
Đăng nhập bằng `device_id` và `password` lấy được từ endpoint trên và từ console của IoT Gateway.
#### **Header**
Yêu cầu header:
1. `Content-Type: application/json`
#### **Body**
```
{
	"device_id": "ID_that_Gateway_creates_on_first_run",
	"password": "raw_password_no_hash"
}
```
#### **Response**
Thành công, trả lại `200 OK` kèm một JWT token.
JWT token này được sử dụng ở các request yêu cầu bảo mật để nhận diện người dùng. Web application nên lưu token này vào Local Storage để sử dụng cho các request cần thiết.
Không thành công (`password` sai, `device_id` không tồn tại,...) , trả lại status code `401 Unauthorized`.
```
{
	"token": "eyJWT.payload.signature"
}
```
### `GET api/backend/data`
Unimplemented. Will use `text/event-stream` for real-time data.
### `GET api/backend/settings/data_rate`
Unimplemented