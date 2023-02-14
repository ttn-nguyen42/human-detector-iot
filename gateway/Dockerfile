FROM python:3.8-slim-buster

WORKDIR /gateway

COPY ./ /gateway/

# Install dependencies
RUN pip install --upgrade setuptools
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

# Environment variables
ENV AWS_IOT_CORE_CERT=
ENV AWS_IOT_CORE_PRIVATE=
ENV AWS_IOT_CORE_PUBLIC=
ENV AWS_IOT_CORE_ENDPOINT=
ENV AWS_IOT_CORE_ROOT_CA=

ENV LOG_LEVEL=

ENV SQLITE_DATABASE=

ENV URL_BACKEND=

RUN mkdir /gateway/db

CMD [ "python3", "main.py" ]
