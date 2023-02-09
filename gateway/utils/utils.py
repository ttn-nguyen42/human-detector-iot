from uuid import uuid4


def get_device_id():
    return f"client{str(uuid4())}"