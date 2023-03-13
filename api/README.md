# API Documentation
## Endpoints
### `POST /api/backend/register_device`
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
### `POST /api/backend/login`
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
### `GET /api/backend/data`
**Quan trọng**: Sử dụng `text/event-stream` ([HTTP SSE](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)) thay cho WebSocket để vận chuyển real-time data.

Lấy dữ liệu trực tiếp từ AWS IoT Core, real-time.
#### **Header**
Yêu cầu header:
- `Content-Type: application/json`
- `Authorization: Bearer {token}`: `{token}` lấy từ response của API `POST /api/backend/login`

#### **Body**
Không yêu cầu body cho các `GET` requests.
#### **Response**
Trả về dữ liệu từ sensor data của thiết bị, liên tục, theo interval được set thông qua settings.
```
{
	"heat_level": 10,
	"light_level": 10,
	"device_id": "ID_that_Gateway_creates_on_first_run",
	"timestamp": 1676347580.9927943
}
```
### `GET /api/backend/check_active`
Kiểm tra xem thiết bị có đang hoạt động hay không. Phải được chạy trước khi hiển thị dashboard.
#### **Header**
Yêu cầu header:
- `Content-Type: application/json`
- `Authorization: Bearer {token}`: `{token}` lấy từ response của API `POST /api/backend/login`
#### **Body**
Không yêu cầu body cho các `GET` requests.
#### **Response**
Trả lại `502 Bad Gateway` nếu gateway không hoạt động và
`503 Service Unavailable` nếu controller không hoạt động và `200 OK` nếu có. Thời gian timeout là 3 giây.
```
{
    "id": "ID_that_Gateway_creates_on_first_run"
}
```
### `POST /api/backend/settings/data_rate`
Thay đổi data rate của thiết bị (bao nhiêu giây gửi dữ liệu một lần). Chỉ được sử dụng request này khi thiết bị đang chạy.
#### **Header**
Yêu cầu header:
- `Content-Type: application/json`
- `Authorization: Bearer {token}`: `{token}` lấy từ response của API `POST /api/backend/login`

#### **Body**

```
{
    "rate_in_seconds": 10
}
```
#### **Response**
Trả lại `503 Service Unavailable` khi thiết bị không được kết nối, `504 Gateway Timeout` khi gateway không hoạt động, `200 OK` khi settings đã được thay đổi.
```
{
    "id: "ID_that_Gateway_creates_on_first_run"
}
```
### `GET /api/backend/settings`
Lấy các settings, ví dụ như data rate. Trong trường hợp user chưa có settings, API này sẽ tạo một default settings cho user đó.
#### **Header**
Yêu cầu header:
- `Content-Type: application/json`
- `Authorization: Bearer {token}`: `{token}` lấy từ response của API `POST /api/backend/login`

#### **Body**
Không yêu body trong các `GET` requests.

#### **Response**
Trả lại `200 OK` kèm settings. Ví dụ default settings:
```
{
    "device_id": "ID_that_Gateway_creates_on_first_run",
    "data_rate": 3
}
```
## Notes
Tất cả các endpoints sẽ trả lại `500 Internal Server Error` trong trường hợp database, MQTT server bị mất kết nối hoặc lỗi trong code của API. Khi này, response sẽ như sau:
```
{
    "message": "Some error message"
}
```

 `400 Bad Request` sẽ được trả lại khi body không đúng format hoặc giá trị. Tương tự như vậy, response được trả lại sẽ có format như trên. 
