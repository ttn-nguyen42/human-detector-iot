"""
Created: nguyen_tran

Keeps all common interfaces, abstract classes
and common models here
"""

import typing


"""
An interface for any MQTT broker implementations
such as AWS IoT Core MQTT or Adafruit, etc.
"""


class MQTTBroker:
    """
    Initiates a connection with the broker
    """

    def connect(self):
        pass

    """
    Publish a message to the broker
    """

    def publish(self, topic: str, payload: typing.Dict[str, any]):
        pass

    """
    Subscribe to a topic
    """

    def subscribe(self, topic: str, func: callable):
        pass
