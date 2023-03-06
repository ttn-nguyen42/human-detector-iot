"""
Created: nguyen_tran

Listens to data from YoloBit controller through Serial (USB)
"""

import logging
import time
from dto.sensor_data import SensorDataDto
from gateway.services.serial_data import ISerialService
from services.sensor_data import ISensorDataService


# Loops indefinitely
# Pause of a second then read serial data from the controller
# Then sends the received data to the service layer
def send_sensor_data(device_id: str, service: ISensorDataService, serial: ISerialService) -> None:
    logging.info(f"Sending sensor data to: 'yolobit/sensor/data/{device_id} ")
    while True:
        data = _read_sensor_data(serial=serial)
        data.device_id = device_id
        service.send_sensor_data(data=data)
        time.sleep(10)


def _read_sensor_data(serial: ISerialService) -> SensorDataDto:
    # Read from serial here
    res: list[str] = serial.read()
    _process_data(res=res)
    

def _process_data(res: list[str]) -> SensorDataDto:
    return SensorDataDto(
        device_id="",
        heat_level=10,
        light_level=10
    )
