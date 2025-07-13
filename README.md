# Home Automation with Golang!

A modern home automation system built with Go and React, designed to control smart outlets and other IoT devices. Currently focused on TP-Link Kasa smart plugs and switches.

## Features

- Smart outlet control (TP-Link Kasa)
  - Auto-discovery on local network
  - Power on/off control
  - Device state monitoring
  - System information retrieval
- Web interface with real-time updates
- RESTful API for device management

## Tech Stack

### Backend
- Go with Gin Web Framework
- Python-Kasa for device control
- Docker for containerization
- Logrus for structured logging

### Frontend
- React with CoreUI Components
- Vite for development
- Axios for API communication

## Getting Started

1. Prerequisites:
   - Docker and Docker Compose
   - Network access to your smart devices

2. Start the application
```bash
sudo docker-compose up --remove-orphans --no-deps --build api ui
```

3. Access:
   - Web User Interface: `http://localhost:3000/#/forms/outlets`
   - API Documentation: `http://localhost:6060`
   - Can also be accessed network wide via your hosts ip `hostname -I` to get ip address

Note: Your device must be on the same network as the smart devices.

## Supported Devices

Currently supports TP-Link Kasa smart devices:
- Smart Plugs (HS103, HS105)
- Smart Switches

## References
- [Python-Kasa](https://github.com/python-kasa/python-kasa)
- [TPLink Smarthome API](https://github.com/plasticrake/tplink-smarthome-api)
- [Hue Control](https://github.com/tigoe/hue-control)

## License

MIT License - See LICENSE file for details