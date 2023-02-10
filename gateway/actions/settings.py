"""
Created: nguyen_tran

Perform settings related actions
"""

import logging
from services.backend import IRemoteBackendService
from services.settings import ILocalSettingsService
from utils.utils import make_device_id

DEVICE_MODEL = "YoloBit Human Detector"


def authenticate(service: ILocalSettingsService, backend: IRemoteBackendService) -> any:
    try:
        saved_id: str = service.get_device_id()
    except Exception as err:
        logging.error(err)
        raise err
    # Should be authenticating against the back-end here
    # POST /api/v1/register_device
    # { "device_id": "something". "model": "YoloBit" }
    try:
        saved_password: str = backend.authenticate(saved_id, DEVICE_MODEL)
    except Exception as err:
        logging.error(err)
        raise err
    return (saved_id, saved_password)
