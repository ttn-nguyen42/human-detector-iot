import logging
import os
from uuid import uuid4

ENV_LOG_LEVEL="LOG_LEVEL"

def get_device_id():
    return f"client-{str(uuid4())}"

def get_log_level() -> int:
    level = os.environ.get(ENV_LOG_LEVEL)
    if len(level) == 0:
        return logging.INFO
    if level == "DEBUG":
        return logging.DEBUG
    if level == "WARNING":
        return logging.WARNING
    if level == "ERROR":
        return logging.ERROR
    return logging.INFO
