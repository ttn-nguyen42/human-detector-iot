"""
Created: nguyen_tran

Perform settings related actions
"""

from ast import Tuple
import logging
from services.backend import IRemoteBackendService
from services.settings import ILocalSettingsService
from utils.utils import make_device_id

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
