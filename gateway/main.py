#!/usr/bin/env python

#
# Entrypoint
#
# First setup serial connection to YoloBit
# Then setup a connection to AWS IoT Core MQTT broker
# Then authenticates the device with the backend through HTTP
#
import logging
import signal
import sys
import threading
from network.aws import Certs, get_certs_path, get_url_endpoint
from network.common import MQTTBroker
from network.mqtt import IotCoreMQTT, PahoMQTT
from reader.command import register_command_subscriber, register_settings_subcriber
from reader.sensor_data import send_sensor_data
from repositories.command import CommandRepository, ICommandRepository
from repositories.sensor_data import ISensorDataRepository, SensorDataRepository
from services.command import CommandService, ICommandService
from services.sensor_data import ISensorDataService, SensorDataService
from utils.utils import get_device_id, get_log_level


def main():
    try:
        logging.basicConfig(level=get_log_level())
        # Get environment variables/settings here
        try:
            certs: Certs = get_certs_path()
            endpoint_url: str = get_url_endpoint()
        except Exception as env_err:
            # Print out exception message
            logging.error(env_err)
            logging.fatal("Missing environment variables")
            # If environment variables are not found, simply exit the program
            sys.exit()

        # Should be authenticating to the backend
        id = get_device_id()
        # Saves this to local database
        logging.info(f"Device ID: {id}")
        # password = register_device(id)
        # print(f"Your device password: {password}")
        #
        # The password and device_id acts as a username and password
        # that we can use in the web app to determine which device that we want to read data from and change settings for

        # Register and connect to databases here
        # aws_mqtt: MQTTBroker = IotCoreMQTT(
        #     endpoint_url=endpoint_url,
        #     ca_filepath=certs.ca,
        #     private_key_filepath=certs.private_key,
        #     cert_filepath=certs.cert,
        #     client_id=id
        # )
        aws_mqtt: MQTTBroker = PahoMQTT(
            ca_path=certs.ca,
            cert_path=certs.cert,
            priv_key_path=certs.private_key,
            endpoint=endpoint_url
        )

        try:
            aws_mqtt.connect()
        except Exception as con_err:
            logging.error(con_err)
            logging.fatal("Unable to connect to message broker after retries")
            sys.exit()

        # Initializes repository layers here
        sensor_data_repository: ISensorDataRepository = SensorDataRepository(
            broker=aws_mqtt
        )

        command_repository: ICommandRepository = CommandRepository(
            broker=aws_mqtt
        )

        # Intializes service layers here
        sensor_data_service: ISensorDataService = SensorDataService(
            sensor_data_repository=sensor_data_repository
        )

        command_service: ICommandService = CommandService(
            command_repository=command_repository
        )

        # Listen to commands
        try:
            register_command_subscriber(id, command_service)
            register_settings_subcriber(id, command_service)
        except Exception as err:
            raise KeyboardInterrupt

        # Listen to sensor data from YoloBit
        # Then sends that to the service layer
        # This acts as a controller layer to complete the MVC architecture
        # Not implemented
        send_sensor_data(id, sensor_data_service)
        return 0
    except KeyboardInterrupt:
        exit()
    finally:
        del aws_mqtt


if __name__ == "__main__":
    main()
