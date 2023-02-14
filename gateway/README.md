# IoT Gateway
## Cách setup dự án
### Yêu cầu
1. Tạo một tài khoản AWS (cần có thẻ VISA/MasterCard, không bị trừ tiền vì có Free Tier)
2. Có Docker trên máy (Windows/WSL/MacOS dùng Docker Desktop, Linux có thể dùng Docker CLI)
3. Có Python (>= 3.8) để có syntax highlighting trong editor
### Hướng dẫn setup source code
Trong Dockerfile, có yêu cầu phải set đủ các environment variables (biến môi trường).
```
ENV AWS_IOT_CORE_CERT=
ENV AWS_IOT_CORE_PRIVATE=
ENV AWS_IOT_CORE_PUBLIC=
ENV AWS_IOT_CORE_ENDPOINT=
ENV AWS_IOT_CORE_ROOT_CA=

ENV LOG_LEVEL=
ENV SQLITE_DATABASE=
```
Các biến này được cài đặt như sau:
1. Tạo một file `.env`. File này bị `.gitignore` nên sẽ không bị đẩy lên GitHub
2. Copy đoạn code sau
```
AWS_IOT_CORE_CERT=
AWS_IOT_CORE_PRIVATE=
AWS_IOT_CORE_PUBLIC=
AWS_IOT_CORE_ENDPOINT=
AWS_IOT_CORE_ROOT_CA=./certs/root-CA.crt
LOG_LEVEL=DEBUG
SQLITE_DATABASE=./db/persistence.db
```
1. Để lấy được giá trị cho các biến còn thiếu, xem phần phía dưới để biết cách setup AWS
2. Sau khi đã chuẩn bị AWS xong, set log level. `DEBUG` sẽ in nhiều log hơn `INFO`
3. Sau đó, set location cho database của SQLite. Ví dụ: `./db/persistence.db`, khi chạy thì dữ liệu local sẽ nằm trong file database này
4. Chạy Docker, mount vị trí của file database (folder `db/` trong ví dụ trên)
- `docker build -t hd_iot/gateway:latest -f ./Dockerfile .`
- `docker run -it -v "$(pwd)/db":"/gateway/db" --env-file=./.env hd_iot/gateway:latest`
`--env-file=./.env` lấy file `.env` vừa tạo để set các environment variable cho dự án.

Lần đầu chạy:
- Khi IoT Gateway chạy lần đầu, phần mềm sẽ tạo một `device_id` ngẫu nhiên, gửi ID đó tới backend qua `POST /api/backend/register_device`
- Gateway sẽ nhận về `password` ngẫu nhiên từ backend trong response
- Lưu chúng vào SQLite database tại `./db/persistence.db`
- Các lần chạy tiếp theo sẽ sử dụng cùng `device_id`, `password` trên
Việc mount database trên sẽ giúp `device_id` và `password` được lưu lại và trong khi testing, developer sẽ không cần thiết phải vào console để đọc log, lấy ID và `password` nữa
vì chúng sẽ được Gateway tái sử dụng

### Hướng dẫn setup AWS
Sau khi đăng nhập vào tài khoản, vào thanh tìm kiếm, tìm `AWS IoT Core`.
1. Bấm vào `Connect device`
2. Kéo xuống dưới, tại mục số 4, copy lại `[URL]` ở đấy
`ping [URL]`. Copy URL này vào file `.env` cho `AWS_IOT_CORE_ENDPOINT`
3. Bấm `Next`, chọn `Create a new thing`, đặt tên cho project, có thể là `human_detector_dev`
4. Bấm `Next`, chọn hệ điều hành đang sử dụng, và chọn Python ở dưới
5. Bấm `Next`, bấm `Download connection kit`. Extract và copy những file sau vào folder dự án, tại `./certs`, cùng folder với `root-CA.crt`
```
human_detector_dev.cert.pem
human_detector_dev.private.key
human_detector_dev.public.key
```
6. Bên trong file `./.env`, điền vào như sau:
```
AWS_IOT_CORE_CERT=./certs/human_detector_dev.cert.pem
AWS_IOT_CORE_PRIVATE=./certs/human_detector_dev.private.key
AWS_IOT_CORE_PUBLIC=./certs/human_detector_dev.public.key
AWS_IOT_CORE_ENDPOINT=cai-gi-do-tai-day.ap-southeast-1.amazonaws.com
AWS_IOT_CORE_ROOT_CA=./certs/root-CA.crt
LOG_LEVEL=DEBUG
SQLITE_DATABASE=db/persistence.database
```
7. Quay lại với AWS. Ở bên trái, bấm vào `All devices`, chọn `Things`. Chọn vào dự án vừa tạo
8. Ở dưới, bấm vào `Certificates`, chọn `Policies`, chọn cái đầu tiên
![AWS Certificates](/docs/images/AWS_Certs.PNG "AWS Certs")
1.  Ở `All versions`, bấm vào hàng đầu tiên, chọn `Edit version`
![AWS Certificates](/docs/images/AWS_Certs_Policy.PNG "AWS Certs")
1.  Tại các trường `Policy resource`, xóa hết và đổi lại thành `*`. Bấm `Save as new version`
![Policies](/docs/images/AWS_Certs_Policies.PNG "Policies")

### Kiểm tra dữ liệu trên AWS
Ứng dụng sẽ gửi dữ liệu lên AWS, để kiểm tra, có thể vào `AWS IoT Core > MQTT test client`.
Tại `Topic filter`, chọn topic muốn lắng nghe. Muốn biết có những topic nào thì đọc code, topic name thường nằm trong `repositories/` hoặc qua log
![IoT Gateway logs](/docs/images/Log_Gateway.PNG "AWS Gateway logs"){align=center}
Mặc định IoT Gateway sẽ lắng nghe (subscribe) 2 topic:
- `yolobit/command/activity/{device_id}`: nhận lệnh `Shutdown`, `Restart`,... từ backend
- `yolobit/command/settings/{device_id}`: nhận thông tin thay đổi về settings từ backend

và publish tới topic:
- `yolobit/sensor/data/{device_id}`: gửi thông tin về sensor data tới backend`

## References
1. [Hướng dẫn sử dụng AWS IoT Core SDK](https://aws.amazon.com/premiumsupport/knowledge-center/iot-core-publish-mqtt-messages-python/)
2. [Cách publish dữ liệu lên MQTT broker của AWS IoT Core SDK](https://dev.to/aws-builders/aws-iot-pubsub-over-mqtt-1oig)
3. [Get started with AWS IoT Core Quick Connect](https://www.youtube.com/watch?v=6w9a6y_-T2o)
4. [MQTT Overview](https://www.youtube.com/watch?v=EIxdz-2rhLs)