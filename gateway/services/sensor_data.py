"""
Created: nguyen_tran

Keeps business logic that interacts with other services or with a repository
"""

import time
from dto.sensor_data import SensorDataDto
from models.sensor_data import SensorData
from repositories.sensor_data import ISensorDataRepository


class ISensorDataService:
    def send_sensor_data(self, data: SensorDataDto) -> None:
        pass


class SensorDataService(ISensorDataService):
    # Implements ISensorDataService
    
    _sensor_data_repostory: ISensorDataRepository = None

    def __init__(self, sensor_data_repository: ISensorDataRepository) -> None:
        # Dependency injection
        self._sensor_data_repostory = sensor_data_repository

    def send_sensor_data(self, data: SensorDataDto) -> None:
        # The service layer implements all business logic
        # For example, this sensor data needs to be sent to somewhere
        # That somewhere will be handled by the repository, which is the database layer
        # The database layer takes in that data and determine where to send it to
        #
        # If we were to change out AWS IoT Core for Adafruit, etc.
        # then we only needs to update the repository
        #
        # This is the layered architecture pattern, it offers "seperation of concern"
        entity = SensorData(
            device_id=data.device_id,
            heat_level=data.heat_level,
            light_level=data.light_level,
            # Additional business logic added here that differs entity (SensorData) from dto (SensorDataDto)
            timestamp=time.time()
        )
        self._sensor_data_repostory.send_sensor_packet(entity)
        return
