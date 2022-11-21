import serial
import time
import requests

# Set up serial port
ser = serial.Serial('COM5', 230400, timeout=0)

# Urls
url_location = "http://52.45.17.177:802/XpertRestApi/api/location_data"
url_alert = "http://52.45.17.177:802/XpertRestApi/api/alert_data"

# Times
postInterval = 5
lastPostTime = time.time()

print("Starting packet reading at: " + time.ctime(lastPostTime))

# Print the first element of serial-read data if the array isn't empty
while True:
    if ser.in_waiting:
            packetTime = time.time()
            packet = ser.readline()
            if (time.time() > lastPostTime + postInterval):
                lastPostTime = packetTime
                print("\n"+ packet.decode('utf'))
                res = requests.post(
                    url_location,
                    json={
                        "deviceimei": 111112222233333,
                        "altitude": 1,
                        "latitude": 38.443976,
                        "longitude": -78.874720,
                        "devicetime": 10,
                        "speed": 0,
                        "Batterylevel": "85",
                        "casefile_id": "string",
                        "address": "string",
                        "positioningmode": "string",
                        "tz": "string",
                        "alert_type": "string",
                        "alert_message": "string",
                        "alert_id": "string",
                        "offender_name": "string",
                        "offender_id": "string"
                    }
                )
                print(res.json())

ser.close()
ser.is_open