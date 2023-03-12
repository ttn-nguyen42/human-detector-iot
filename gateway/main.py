#!/usr/bin/env python

#
# Entrypoint
#
# First setup serial connection to YoloBit
# Then setup a connection to AWS IoT Core MQTT broker
# Then authenticates the device with the backend through HTTP
#
import http.client
import logging
import sys
import threading
from database.sqlite import SQLDatabase, SqliteDatabase, get_sqlite_database
from network.aws import Certs, get_certs_path, get_url_endpoint
from network.mqtt import MQTTBroker
from network.mqtt import PahoMQTT
from actions.command import register_command_subscriber
from actions.sensor_data import send_sensor_data
from actions.settings import authenticate, get_token, get_setings
from repositories.command import CommandRepository, ICommandRepository
from repositories.sensor_data import ISensorDataRepository, SensorDataRepository
from repositories.settings import ILocalSettingsRepository, LocalSettingsRepository
from services.backend import IRemoteBackendService, RemoteBackendService
from services.command import CommandService, ICommandService
from services.sensor_data import ISensorDataService, SensorDataService
from services.serial_data import ISerialService, SerialService
from services.settings import ILocalSettingsService, LocalSettingsService
from utils.utils import get_backend_port, get_backend_url, get_log_level
from models.commands import *
from dto.commands import *
from dto.settings import *


def main():
    global DATA_RATE
    try:
        logging.basicConfig(level=get_log_level())
        # Get environment variables/settings here
        try:
            certs: Certs = get_certs_path()
            endpoint_url: str = get_url_endpoint()
            sqlite_db_name: str = get_sqlite_database()
            backend_url: str = get_backend_url()
            backend_port: int = get_backend_port()
        except Exception as env_err:
            # Print out exception message
            logging.fatal("Missing environment variables: {0}".format(env_err))
            # If environment variables are not found, simply exit the program
            sys.exit()

        # Establish HTTPS connection
        http_con: http.client.HTTPConnection
        port = None
        if backend_port != 0:
            port = backend_port
        conn_str = ""
        if "https" in backend_url:
            conn_str = backend_url.replace("https://", "")
            http_con = http.client.HTTPSConnection(conn_str, port, timeout=5)
        else:
            conn_str = backend_url.replace("http://", "")
            http_con = http.client.HTTPConnection(conn_str, port, timeout=5)
            logging.debug(f"Connecting to {backend_url}:{port}")
            remote_backend_service: IRemoteBackendService = RemoteBackendService(
                # Should inject HTTP client here
                http_con=http_con
            )

        # Serial data
        # Will fails when no connected device is found
        try:
            serial_con: ISerialService = SerialService("")
            serial_con.connect()
        except Exception as e:
            logging.info("Unable to find serial IOT device: {0}".format(e))
            sys.exit()

        # Connecting to the database
        db: SQLDatabase = SqliteDatabase(sqlite_db_name).connect()

        local_db_repository: ILocalSettingsRepository = LocalSettingsRepository(
            database=db
        ).initialize()

        local_db_service: ILocalSettingsService = LocalSettingsService(
            local_repository=local_db_repository,
            be_service=remote_backend_service
        )

        # The password and device_id acts as a username and password that we can use in the web app to determine
        # which device that we want to read data from and change settings for

        # Register and connect to cloud services here
        #
        # Example using AWS IoT SDK
        aws_mqtt: MQTTBroker = PahoMQTT(
            ca_path=certs.ca,
            cert_path=certs.cert,
            priv_key_path=certs.private_key,
            endpoint=endpoint_url
        )

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
            command_repository=command_repository,
            serial=serial_con
        )

        # Connect to AWS MQTT
        try:
            aws_mqtt.connect()
        except Exception as con_err:
            logging.info(
                "Unable to connect to message broker after retries: {0}".format(con_err))
            db.close()
            sys.exit()

        # Authenticating to the backend

        # This will create/retrieve a device_id, then send it to the backend
        #  will return a password associated with this device
        try:
            device_id, password = authenticate(
                service=local_db_service,
            )
        except Exception as err:
            logging.info(
                "Cannot connect with the backend: {0}".format(err))
            db.close()
            aws_mqtt.disconnect()
            sys.exit()

        logging.info(
            "Use this credentials to authenticate on web app and monitor this device")
        logging.info(f"Device ID: {device_id}")
        logging.info(f"Password: {password}")

        # Retrieve a JWT token
        try:
            token = get_token(service=remote_backend_service, deviceId=device_id, password=password)
        except Exception as err:
            logging.fatal("Unable to get data rate, stopping")
            db.close()
            aws_mqtt.disconnect()
            sys.exit()
        
        # Retrieve settings
        try:
            settings = get_setings(service=remote_backend_service)
            DATA_RATE = settings.data_rate
            logging.info(f"Data rate: {DATA_RATE}")
        except Exception as err:
            logging.fatal("Unable to get data rate, stopping")
            db.close()
            aws_mqtt.disconnect()
            sys.exit()

        # Listen to commands
        try:
            register_command_subscriber(device_id, command_service)
        except Exception as reg_err:
            logging.info(
                "Unable to register to command and settings topics: {0}".format(reg_err))
            db.close()
            aws_mqtt.disconnect()
            sys.exit()

        # Listen to sensor data from YoloBit
        # Then sends that to the service layer
        # This acts as a controller layer to complete the MVC architecture
        # Not implemented
        send_sensor_data(device_id, sensor_data_service, serial_con)
        command_service.messenger_thread.join()
        return 0
    except KeyboardInterrupt:
        sys.exit()


if __name__ == "__main__":
    main()
