"""
Created: nguyen_tran

Listens to data from YoloBit controller through Serial (USB)
"""

import time
from dto.sensor_data import SensorDataDto
from services.sensor_data import ISensorDataService

# Loops indefinitely
# Pause of a second then read serial data from the controller
# Then sends the received data to the service layer
def send_sensor_data(device_id: str, service: ISensorDataService) -> None:
    while True:
        data = _read_sensor_data()
        data.device_id = device_id
        service.send_sensor_data(data=data)
        time.sleep(2)
    return

def _read_sensor_data() -> SensorDataDto:
    # Read from serial here
    return SensorDataDto(
        device_id="",
        heat_level=10,
        light_level=10
    )

def listen_settings_data():
    pass