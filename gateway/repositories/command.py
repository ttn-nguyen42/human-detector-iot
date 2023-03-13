"""
Created: nguyen_tran

Repository layer for listening to backend's command
"""

import logging
from network.mqtt import MQTTBroker
from models.commands import CommandResponse


class ICommandRepository:
    # Subcribes to a MQTT topic that listens to commands (shutdown, start,...) between backend and gateway
    def subscribe_command(self, device_id: str, callback: callable) -> None:
        pass

    def send_command_response(self, device_id: str, res: CommandResponse) -> None:
        pass


class CommandRepository(ICommandRepository):
    # Implements ICommandRepository
    _broker: MQTTBroker = None
    _topic: str = "yolobit/command"

    def __init__(self, broker: MQTTBroker) -> None:
        self._broker = broker
        return

    # Subcribes to a MQTT topic that listens to commands (shutdown, start,...) between backend and gateway
    def subscribe_command(self, device_id: str, callback: callable) -> None:
        topic = f"{self._topic}/activity/{device_id}"
        try:
            self._broker.subscribe(topic=topic, func=callback)
        except Exception as err:
            logging.error(f"Unable to subscribe to {topic}")
            raise err
        return

    def send_command_response(self, device_id: str, res: CommandResponse):
        topic = f"{self._topic}/response/{device_id}"
        self._broker.publish(topic=topic, payload=res)
        return
