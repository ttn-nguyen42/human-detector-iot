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
```
3. Để lấy được giá trị cho các biến còn thiếu, xem phần phía dưới để biết cách setup AWS
4. Sau khi đã chuẩn bị AWS xong, chạy project bằng Docker như sau:
1. `docker build -t project/iot_gateway -f ./Dockerfile .`
2. `docker run -it --env-file=./.env project/iot_gateway`
`--env-file=./.env` lấy file `.env` vừa tạo để set các environment variable cho dự án.

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
```
7. Quay lại với AWS. Ở bên trái, bấm vào `All devices`, chọn `Things`. Chọn vào dự án vừa tạo
8. Ở dưới, bấm vào `Certificates`, chọn `Policies`, chọn cái đầu tiên
9. Ở `All versions`, bấm vào hàng đầu tiên, chọn `Edit version`
10. Tại các trường `Policy resource", xóa hết và đổi lại thành `*`. Bấm `Save as new version`