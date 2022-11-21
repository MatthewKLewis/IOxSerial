import serial

# Set up serial port
ser = serial.Serial('COM5', 230400, timeout=0)

# Print the first element of serial-read data if the array isn't empty
while True:
    data = []
    data.append(ser.readlines())
    if (len(data[0]) > 0):
        print(data)
    # further processing 
    # send the data somewhere else etc

ser.close()
ser.is_open