"""
Created: nguyen_tran

Update settings, device_id and perform necessary changes to the YoloBit device
"""

from services.backend import IRemoteBackendService
from repositories.settings import ILocalSettingsRepository
from utils.utils import make_device_id
from dto.commands import *
from models.settings import *


class ILocalSettingsService:
    def get_device_id(self) -> any:
        pass

    def save_data_rate(self, device_id: str, settings: DataRateCommandSettings) -> any:
        pass

    def save_settings(self, device_id: str, settings: DeviceSettings) -> any:
        pass

    def get_data_rate(self, device_id: str) -> any:
        pass


class LocalSettingsService(ILocalSettingsService):
    # Implements ILocalSettingsService

    _local_repository: ILocalSettingsRepository = None
    _backend_service: IRemoteBackendService = None

    def __init__(self, local_repository: ILocalSettingsRepository, be_service: IRemoteBackendService) -> None:
        self._local_repository = local_repository
        self._backend_service = be_service
        return

    def save_data_rate(self, device_id: str, settings: DataRateCommandSettings) -> any:
        try:
            local_settings = self._local_repository.get_settings()
        except Exception as err:
            logging.error(f"Unable to get settings {err}")
            raise err
        local_settings.data_rate = settings.rate_in_seconds
        try:
            self._local_repository.save_settings(device_id=device_id, settings=local_settings)
        except Exception as err:
            logging.error(f"Unable to save new settings {err}")
            raise err
        return

    def save_settings(self, device_id: str, settings: DeviceSettings) -> any:
        try:
            self._local_repository.save_settings(device_id, settings)
        except Exception as err:
            logging.error(f"Unable to get local settings {err}")
            raise err
        return

    def get_data_rate(self, device_id: str) -> any:
        try:
            local_settings = self._local_repository.get_settings()
        except Exception as err:
            logging.error(f"Unable to get settings {err}")
            raise err
        return local_settings.data_rate

    def get_device_id(self) -> any:
        try:
            saved_id, saved_password = self._local_repository.get_device_id()
        except Exception as err:
            logging.error(f"Unable to get device ID from database ID={err}")
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
