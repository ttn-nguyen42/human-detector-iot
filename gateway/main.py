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
from database.sqlite import SQLDatabase, SqliteDatabase, get_sqlite_database
from network.aws import Certs, get_certs_path, get_url_endpoint
from network.mqtt import MQTTBroker
from network.mqtt import PahoMQTT
from actions.command import register_command_subscriber, register_settings_subcriber
from actions.sensor_data import send_sensor_data
from actions.settings import authenticate
from repositories.command import CommandRepository, ICommandRepository
from repositories.sensor_data import ISensorDataRepository, SensorDataRepository
from repositories.settings import ILocalSettingsRepository, LocalSettingsRepository
from services.backend import IRemoteBackendService, RemoteBackendService
from services.command import CommandService, ICommandService
from services.sensor_data import ISensorDataService, SensorDataService
from services.settings import ILocalSettingsService, LocalSettingsService
from utils.utils import get_backend_port, get_backend_url, get_log_level


def main():
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
            logging.error(env_err)
            logging.fatal("Missing environment variables")
            # If environment variables are not found, simply exit the program
            sys.exit()

        # Connecting to the database
        db: SQLDatabase = SqliteDatabase(sqlite_db_name).connect()
        local_db_repository: ILocalSettingsRepository = LocalSettingsRepository(
            database=db
        ).initialize()

        # Establish HTTPS connection
        http_con: http.client.HTTPConnection = None
        port = None
        if backend_port != 0:
            port = backend_port
        conn_str = ""
        if "https" in backend_url:
            conn_str = backend_url.replace("https://", "")
            http_con = http.client.HTTPSConnection(
                conn_str, port, timeout=5)
        else:
            conn_str = backend_url.replace("http://", "")
            http_con = http.client.HTTPConnection(conn_str, port, timeout=5)
        logging.debug(f"Connecting to {backend_url}:{port}")
        remote_backend_service: IRemoteBackendService = RemoteBackendService(
            # Should inject HTTP client here
            http_con=http_con
        )

        local_db_service: ILocalSettingsService = LocalSettingsService(
            local_repository=local_db_repository,
            be_service=remote_backend_service
        )

        # Should be authenticating to the backend

        # This will create/retrieve a device_id, then send it to the backend
        # The backend will returns a password associated with this device
        try:
            id, password = authenticate(
                service=local_db_service,
            )
        except Exception as err:
            logging.info(
                "Cannot connect with the backend, probably because it is not available")
            logging.debug(err)
            db.close()
            sys.exit()

        logging.info(
            "Use this credentials to authenticate on web app and monitor this device")
        logging.info(f"Device ID: {id}")
        logging.info(f"Password: {password}")

        # The password and device_id acts as a username and password
        # that we can use in the web app to determine which device that we want to read data from and change settings for

        # Register and connect to cloud services here
        #
        # Example using AWS IoT SDK
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
            db.close()
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
            db.close()
            aws_mqtt.close()
            sys.exit()

        # Listen to sensor data from YoloBit
        # Then sends that to the service layer
        # This acts as a controller layer to complete the MVC architecture
        # Not implemented
        send_sensor_data(id, sensor_data_service)
        return 0
    except KeyboardInterrupt:
        sys.exit()


if __name__ == "__main__":
    main()
