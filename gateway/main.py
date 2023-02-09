#!/usr/bin/env python

#
# Entrypoint
#
# First setup serial connection to YoloBit
# Then setup a connection to AWS IoT Core MQTT broker
# Then authenticates the device with the backend through HTTP
#
import signal
import sys
import threading
from network.aws import Certs, get_certs_path, get_url_endpoint
from network.common import MQTTBroker
from network.mqtt import IotCoreMQTT
from reader.command_read import register_command_subscriber
from reader.controller_read import send_sensor_data
from repositories.command import CommandRepository, ICommandRepository
from repositories.sensor_data import ISensorDataRepository, SensorDataRepository
from services.command import CommandService, ICommandService
from services.sensor_data import ISensorDataService, SensorDataService
from utils.utils import get_device_id


def main():
    try:
        # Get environment variables/settings here
        try:
            certs: Certs = get_certs_path()
            endpoint_url: str = get_url_endpoint()
        except Exception as env_err:
            # Print out exception message
            print(env_err.args[0])
            # If environment variables are not found, simply exit the program
            sys.exit()

        # Should be authenticating to the backend
        #
        id = get_device_id()
        print(f"Device ID: {id}")
        # password = register_device(id)
        # print(f"Your device password: {password}")
        #
        # The password and device_id acts as a username and password
        # that we can use in the web app to determine which device that we want to read data from and change settings for

        # Register and connect to databases here
        aws_mqtt: MQTTBroker = IotCoreMQTT(
            endpoint_url=endpoint_url,
            ca_filepath=certs.ca,
            private_key_filepath=certs.private_key,
            cert_filepath=certs.cert,
            client_id=id
        )

        try:
            aws_mqtt.connect()
        except Exception as con_err:
            print("Unable to connect to message broker after retries")
            print(con_err.args)
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
            sensorDataRepository=sensor_data_repository
        )

        command_service: ICommandService = CommandService(
            commandRepository=command_repository
        )

        # Listen to commands
        register_command_subscriber(id, command_service)

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
