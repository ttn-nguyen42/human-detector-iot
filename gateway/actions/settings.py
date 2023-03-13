"""
Created: nguyen_tran

Perform settings related actions
"""

from services.backend import IRemoteBackendService
from services.settings import ILocalSettingsService
from models.settings import *

DEVICE_MODEL = "YoloBit Human Detector"


def authenticate(service: ILocalSettingsService) -> any:
    # Generate a device ID on first run
    # Register itself with the server
    # Provide a password on return 
    try:
        saved = service.get_device_id()
    except Exception as err:
        raise err
    return saved[0], saved[1]


def get_token(service: IRemoteBackendService, device_id: str, password: str) -> any:
    # Retrieve a JWT token on startup
    try:
        token = service.get_token(device_id, password)
    except Exception as err:
        raise err
    return token


def get_setings(service: IRemoteBackendService, local_service: ILocalSettingsService, device_id: str) -> any:
    # Retrieve settings from the backend
    try:
        settings = service.get_settings()
        local_service.save_settings(device_id, settings=DeviceSettings(
            data_rate=settings.data_rate
        ))
    except Exception as err:
        raise err
    return settings
