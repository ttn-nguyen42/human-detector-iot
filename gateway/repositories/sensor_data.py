"""
Created: nguyen_tran

Repository layer for communicating with a back-end or database/message broker
"""

from models.sensor_data import SensorData
from network.common import MQTTBroker


class ISensorDataRepository:

    # Takes in a sensor data then transfer that to a message broker
    def send_sensor_packet(self, packet: SensorData) -> None:
        pass


class SensorDataRepository(ISensorDataRepository):
    # Implements ISensorDataRepository

    _broker: MQTTBroker = None
    _topic: str = ""

    def __init__(self, broker: MQTTBroker) -> None:
        self._broker = broker
        self._topic = "yolobit/data/sensor"
        return

    def send_sensor_packet(self, packet: SensorData) -> None:
        topic = f"{self._topic}/{packet.device_id}"
        self._broker.publish(topic, packet)
        return
