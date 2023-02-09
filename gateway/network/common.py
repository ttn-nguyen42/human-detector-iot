"""
Created: nguyen_tran

Keeps all common interfaces, abstract classes
and common models here
"""

import typing


class MQTTBroker:
    # An interface for any MQTT broker implementations
    # such as AWS IoT Core MQTT or Adafruit, etc.

    # Initiates a connection with the MQTT broker
    def connect(self) -> None:
        pass
    
    def disconnect(self) -> None:
        pass

    # Publish a message to the broker
    def publish(self, topic: str, payload: typing.Dict[str, any]) -> None:
        pass

    # Subscribe to a topic
    # Returns subscription settings
    # Data can be access through func() as a callback
    def subscribe(self, topic: str, func: callable) -> None:
        pass
    
