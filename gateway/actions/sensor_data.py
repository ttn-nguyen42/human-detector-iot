"""
Created: nguyen_tran

Listens to data from YoloBit controller through Serial (USB)
"""

import logging
import time
from typing import List
from dto.sensor_data import SensorDataDto
from services.serial_data import ISerialService
from services.sensor_data import ISensorDataService


# Loops indefinitely
# Pause of a second then read serial data from the controller
# Then sends the received data to the service layer
def send_sensor_data(device_id: str, service: ISensorDataService, serial: ISerialService) -> None:
    logging.info(f"Sending sensor data to: 'yolobit/sensor/data/{device_id}")
    while True:
        logging.info("Read")
        # Read from serial here
        res: list[str] = serial.read()
        for dat in res:
            try:
                model = _process_data(res=dat)
            except Exception as e:
                logging.info("{0} - {1}".format(e, dat))
                continue
            model.device_id = device_id
            service.send_sensor_data(data=model)
        time.sleep(2)
    return

def _process_data(res: str) -> SensorDataDto:
    try:
        model = SensorDataDto(res)
    except Exception as e:
        raise e
    return model
