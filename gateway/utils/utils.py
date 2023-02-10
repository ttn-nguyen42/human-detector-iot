import logging
import os
from uuid import uuid4

ENV_LOG_LEVEL="LOG_LEVEL"

def make_device_id():
    # Generates a random device ID
    return f"client_{str(uuid4())}"

def get_log_level() -> int:
    # Reads LOG_LEVEL environment variable
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
