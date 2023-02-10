"""
Created: nguyen_tran

Update settings, device_id and perform necessary changes to the YoloBit device
"""

import logging
from repositories.settings import ILocalSettingsRepository, LocalSettingsRepository
from utils.utils import make_device_id


class ILocalSettingsService:
    def get_device_id(self) -> str:
        pass

class LocalSettingsService(ILocalSettingsService):
    # Implements ILocalSettingsService
    
    _local_repository: ILocalSettingsRepository = None

    def __init__(self, local_repository: LocalSettingsRepository) -> None:
        self._local_repository = local_repository
        return

    def get_device_id(self) -> str:
        try:
            saved_id = self._local_repository.get_device_id()
        except Exception as err:
            raise err
        if saved_id != "":
            return saved_id
        logging.debug(f"Found ID={saved_id}")
        saved_id = make_device_id()
        logging.debug(f"Created ID={saved_id}")
        try:
            self._local_repository.save_device_id(saved_id, "YoloBit Human Detector")
        except Exception as err:
            logging.error(err)
            raise err
        return saved_id
