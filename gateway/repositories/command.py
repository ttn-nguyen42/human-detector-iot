"""
Created: nguyen_tran

Repository layer for listening to backend's command
"""


import logging
from network.common import MQTTBroker


class ICommandRepository:
    # Subcribes to a MQTT topic that listens to commands (shutdown, start,...) between backend and gateway
    def subscribe_command(self, device_id: str, callback: callable) -> None:
        pass
    
    # Subscribes to a MQTT topic that transfers settings data between backend and gateway
    def subscribe_settings(self, device_id: str, callback: callable) -> None:
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
        
    # Subscribes to a MQTT topic that transfers settings data between backend and gateway
    def subscribe_settings(self, device_id: str, callback: callable) -> None:
        topic = f"{self._topic}/settings/{device_id}"
        try:
            self._broker.subscribe(topic=topic, func=callback)
        except Exception as err:
            logging.error(f"Unable to subscribe to {topic}")
            raise err
        return
        
