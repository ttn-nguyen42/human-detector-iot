"""
Created: nguyen_tran

Update settings, device_id and perform necessary changes to the YoloBit device
"""

from ast import Tuple
import logging
from services.backend import IRemoteBackendService
from repositories.settings import ILocalSettingsRepository, LocalSettingsRepository
from utils.utils import make_device_id


class ILocalSettingsService:
    def get_device_id(self) -> any:
        pass


class LocalSettingsService(ILocalSettingsService):
    # Implements ILocalSettingsService

    _local_repository: ILocalSettingsRepository = None
    _backend_service: IRemoteBackendService = None

    def __init__(self, local_repository: LocalSettingsRepository, be_service: IRemoteBackendService) -> None:
        self._local_repository = local_repository
        self._backend_service = be_service
        return

    def get_device_id(self) -> any:
        try:
            saved_id, saved_password = self._local_repository.get_device_id()
        except Exception as err:
            raise err
        if saved_id != "":
            return saved_id, saved_password
        logging.debug(f"Found ID={saved_id}")
        saved_id = make_device_id()
        logging.debug(f"Created ID={saved_id}")
        try:
            received_pass = self._backend_service.authenticate(saved_id)
        except Exception as err:
            raise err
        try:
            self._local_repository.save_device_id(
                saved_id, received_pass, "YoloBit Human Detector")
        except Exception as err:
            raise err
        return saved_id, received_pass
