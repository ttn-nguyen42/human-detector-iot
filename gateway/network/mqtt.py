"""
Created: nguyen_tran

Handles MQTT traffic between this gateway and AWS IoT Core
"""

from concurrent.futures import CancelledError
import json
import logging
import ssl
from threading import Lock
import typing

import paho.mqtt.client as paho
from awscrt import mqtt
from awsiot import mqtt_connection_builder


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


def _on_publish(_, __, message_id) -> None:
    logging.debug("Published message ID={0} success".format(message_id))
    return


def _on_subscribe(_, __, mid, granted_qos) -> None:
    logging.debug("Subscribed ID={0}, QoS={1}".format(mid, granted_qos[0]))
    return


class PahoMQTT(MQTTBroker):
    # Implements MQTTBroker
    _is_connected: bool = False
    _client: paho.Client = None
    _endpoint: str = ""
    _port: int = 0

    def __init__(self, ca_path: str, priv_key_path: str, cert_path: str, endpoint: str) -> None:
        self._is_connected = False
        self._client = paho.Client()
        self._client.on_connect = self._on_connected
        self._client.on_disconnect = self._on_disconnect
        self._endpoint = endpoint
        self._port = 8883
        self._client.tls_set(
            ca_certs=ca_path,
            keyfile=priv_key_path,
            certfile=cert_path,
            cert_reqs=ssl.CERT_REQUIRED,
            tls_version=ssl.PROTOCOL_TLSv1_2,
            ciphers=None
        )
        return

    def _on_connected(self, _, __, flags, rc) -> None:
        self._is_connected = True
        logging.info("Connected to AWS IoT Core")
        logging.debug("Connected result={0}, flags={1}".format(rc, flags))
        return

    def _on_disconnect(self, _, rc) -> None:
        self._is_connected = False
        logging.info("Disconnected from AWS IoT Core")
        logging.debug("Disconnected result={0}".format(rc))
        return

    # Initiates a connection with the MQTT broker
    def connect(self) -> None:
        result = self._client.connect(
            self._endpoint, self._port, keepalive=120)
        if result != 0:
            raise ConnectionRefusedError("Cannot connect to AWS IoT Core")
        self._client.loop_start()
        return

    def disconnect(self) -> None:
        result = self._client.disconnect()
        logging.info("Disconnected MQTT broker")
        if result != 0:
            raise ConnectionAbortedError("Cannot disconnect from AWS IoT Core")
        return

    # Publish a message to the broker
    def publish(self, topic: str, payload: typing.Dict[str, any]) -> None:
        self._client.on_publish = _on_publish
        result = None
        try:
            pl = json.dumps(payload)
            result = self._client.publish(
                topic=topic, payload=pl)
        except Exception as err:
            logging.error(f"Unable to publish err={err}")
        try:
            result.wait_for_publish(timeout=5)
        except Exception as err:
            logging.debug(f"Publish topic={topic} not send error={err}")
        logging.info(f"Sent message ID={result.mid} to topic={topic}")
        return

    # Subscribe to a topic
    # Returns subscription settings
    # Data can be access through func() as a casllback

    def subscribe(self, topic: str, func: callable) -> None:
        self._client.message_callback_add(topic, func)
        self._client.on_subscribe = _on_subscribe
        result, _ = self._client.subscribe(topic=topic)
        if result == paho.MQTT_ERR_SUCCESS:
            logging.info(f"Subscribed to topic {topic}")
            return
        raise Exception(f"Cannot subscribe message to {topic}")

    def __del__(self) -> None:
        self.disconnect()
        return


class IotCoreMQTT(MQTTBroker):
    # Implements MQTTBroker

    #
    # Not used
    #
    _connection: mqtt.Connection = None
    _is_connected: bool = False
    _lock: Lock = Lock()

    def __init__(self, endpoint_url: str,
                 cert_filepath: str,
                 private_key_filepath: str,
                 ca_filepath: str,
                 client_id: str) -> None:
        self._is_connected = False
        self._connection = mqtt_connection_builder.mtls_from_path(
            endpoint=endpoint_url,
            cert_filepath=cert_filepath,
            pri_key_filepath=private_key_filepath,
            ca_filepath=ca_filepath,
            client_id=client_id
        )
        return

    def __del__(self) -> None:
        self.disconnect()
        return

    def connect(self) -> None:
        if self._is_connected:
            logging.error("Already connected to AWS IoT Core MQTT broker")
            return
        connect_future = self._connection.connect()
        try:
            connect_future.result(5)
            self._is_connected = True
            logging.info("Connected to AWS IoT Core")
        except CancelledError as err:
            logging.error("Connection to AWS IoT Core cancelled")
            raise err
        except TimeoutError as err:
            logging.error("Connection to AWS IoT Core timed out")
            raise err
        except Exception as err:
            logging.error(
                "Connection to AWS IoT Core failed for unknown reason")
            raise err

    def disconnect(self) -> None:
        dist_future = self._connection.disconnect()
        dist_future.result()
        return

    def publish(self, topic: str, payload: typing.Dict[str, any]) -> None:
        self._connection.publish(
            topic=topic,
            payload=json.dumps(payload),
            qos=mqtt.QoS.AT_LEAST_ONCE
        )
        logging.info(f"Message published to {topic}")
        return

    def subscribe(self, topic: str, func: callable) -> None:
        subscribe_future, packet_id = self._connection.subscribe(
            topic=topic,
            qos=mqtt.QoS.AT_LEAST_ONCE,
            callback=func
        )
        try:
            subscribe_result = subscribe_future.result()
            logging.debug("Result: {}".format(str(subscribe_result['qos'])))
        except Exception as err:
            # print(err)
            raise err
        logging.info(f"Subcribed to {topic}, packet ID {packet_id}")
        return
