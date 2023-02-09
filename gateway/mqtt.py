"""
Created: nguyen_tran

Handles MQTT traffic between this gateway and AWS IoT Core
"""

from concurrent.futures import CancelledError
import json
from threading import Lock
import typing
from gateway.common import MQTTBroker
from awscrt import io, mqtt, auth, http
from awsiot import mqtt_connection_builder


class IotCoreMQTT(MQTTBroker):

    _connection: mqtt.Connection = None
    _is_connected: bool = False
    _lock: Lock = Lock()

    def __init__(self, endpoint_url: str, cert_filepath: str, private_key_filepath: str, ca_filepath: str, client_id: str) -> None:
        super().__init__()
        self._is_connected = False
        self._connection = mqtt_connection_builder.mtls_from_path(
            endpoint=endpoint_url,
            cert_filepath=cert_filepath,
            pri_key_filepath=private_key_filepath,
            ca_filepath=ca_filepath,
            client_id=client_id
        )

    """
    Initiates a connection with the broker
    """

    def connect(self):
        if self._is_connected == True:
            print("Already connected to AWS IoT Core MQTT broker")
            return
        connect_future = self._connection.connect()
        try:
            connect_future.result(5)
            self._is_connected = True
        except CancelledError:
            print("Connection to AWS IoT Core cancelled")
            return
        except TimeoutError:
            print("Connection to AWS IoT Core timed out")
            return
        else:
            print("Connection to AWS IoT Core failed for unknown reason")
            return

    """
    Publish a message to the broker
    """

    def publish(self, topic: str, payload: typing.Dict[str, any]):
        self._connection.publish(
            topic=topic,
            payload=json.dumps(payload),
            qos=mqtt.QoS.AT_LEAST_ONCE
        )
        print("Message published")
        return

    """
    Subscribe to a topic
    """

    def subscribe(self, topic: str, func: callable):
        subscribe_future, packet_id = self._connection.subscribe(
            topic=topic,
            qos=mqtt.QoS.AT_LEAST_ONCE,
            callback=func
        )
        subscribe_result = subscribe_future.result()
        print(f"Received from {topic}, packet ID {packet_id}")
        return
