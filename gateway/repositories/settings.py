from ast import Tuple
import logging
from database.sqlite import SQLDatabase
from models.settings import DeviceSettings


class ILocalSettingsRepository:
    # Create necessary tables in the database
    def initialize(self) -> any:
        pass

    # Save the device_id into the table for it
    def save_device_id(self, device_id: str, password: str, model: str) -> None:
        pass

    # Save the settings into the table for it
    def save_settings(self, device_id: str, settings: DeviceSettings) -> None:
        pass

    # Get the device_id from the database
    def get_device_id(self) -> any:
        pass

    # Get the settings from the database
    def get_settings(self) -> DeviceSettings:
        pass


class LocalSettingsRepository(ILocalSettingsRepository):
    # Interacts with the local SQL database
    # for storing settings and other information

    _database: SQLDatabase = None
    _initialized: bool = False
    _device_info_db: str = ""
    _settings_db: str = ""

    def __init__(self, database: SQLDatabase) -> None:
        self._database = database
        self._device_info_db = "device_info"
        self._settings_db = "device_settings"
        return

    def initialize(self) -> any:
        # Create a table for saving device-specific info
        try:
            self._database.execute(f"""CREATE TABLE IF NOT EXISTS {self._device_info_db} (
                                id              INTEGER PRIMARY KEY NOT NULL,
                                device_id       TEXT NOT NULL,
                                password        TEXT NOT NULL,
                                model           TEXT        
                               );                      
                               """)
        except Exception as err:
            raise err

        logging.debug(f"Created {self._device_info_db}")

        # Create a table for saving settings
        try:
            self._database.execute(f"""CREATE TABLE IF NOT EXISTS {self._settings_db}(
                                id                      INTEGER PRIMARY KEY NOT NULL,
                                device_id               TEXT NOT NULL,
                                data_rate               INTEGER
                               );
                               """)
        except Exception as err:
            raise err

        logging.debug(f"Created {self._settings_db}")
        self._initialized = True
        return self

    def save_device_id(self, device_id: str, password: str, model: str) -> None:
        if self._initialized == False:
            raise Exception("Table has not been initialized")
        try:
            self._database.execute(f"""INSERT OR REPLACE INTO {self._device_info_db}
                               (id, device_id, password, model)
                               VALUES
                               ({1}, "{device_id}", "{password}", "{model}");
                               """)
            self._database.commit()
        except Exception as err:
            raise err
        logging.debug(
            f"Added device info to database table={self._device_info_db}")
        return

    def save_settings(self, device_id: str, settings: DeviceSettings) -> None:
        if self._initialized == False:
            raise Exception("Table has not been initialized")
        try:
            self._database.execute(f"""INSERT OR REPLACE INTO {self._settings_db}
                               (id, device_id, data_rate)
                               VALUES
                               ({1}, "{device_id}", {settings.data_rate})
                               """)
            self._database.commit()
        except Exception as err:
            raise err
        logging.debug(
            f"Settings just got saved into the database table={self._settings_db}")
        return

    def get_device_id(self) -> any:
        if self._initialized == False:
            raise Exception("Table has not been initialized")
        try:
            cursor = self._database.execute(f"""
                                     SELECT DISTINCT * FROM {self._device_info_db}
                                     LIMIT 1
                                     """)
        except Exception as err:
            raise err
        device_id: str = ""
        password: str = ""
        logging.debug(f"Received length={cursor.arraysize}")
        for row in cursor:
            logging.info(
                f"Table={self._device_info_db} device_id={row[1]} model={row[3]}")
            device_id = row[1]
            password = row[2]

        return device_id, password

    def get_settings(self) -> DeviceSettings:
        if self._initialized == False:
            raise Exception("Table has not been initialized")

        try:
            cursor = self._database.execute(f"""
                                        SELECT DISTINCT * FROM {self._settings_db}
                                        LIMIT 1
                                        """)
        except Exception as err:
            raise err
        settings: DeviceSettings = DeviceSettings()
        for row in cursor:
            logging.info(
                f"Table={self._settings_db} device_id={row[1]} data_rate={row[2]}")
            settings.data_rate = row[2]
        return settings
