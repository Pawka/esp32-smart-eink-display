# Smart display for ESP32: server

Smart display is build from Waveshare e-ink display and ESP32 microcontroller.
The ESP32 periodically fetch data from project's server and displays on the
screen. Data format is known for the ESP32. Adding additional information to the
display only requires changes on the server side but not the controller. This
repository contains server part of the project.

## Building & Running

It is recommended to run the server on Docker container. Container can be build
directly from the repository:

```
$ docker build -t smartdisplay https://github.com/Pawka/esp32-smart-eink-display.git
```

Running the container:

```
$ docker run -p 3000:3000 -d smartdisplay
```
