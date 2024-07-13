# Lite-API Proxy for Hotelbeds Hotels Search

Lite-API is a Go application that proxies requests to the Hotelbeds hotels availability [endpoint](https://developer.hotelbeds.com/documentation/hotels/booking-api/api-referenhttps://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#tag/Availability), providing a convenient way to interact with Hotelbeds API via a local proxy server.

## Features

- **Proxying to Hotelbeds API**: Routes requests to Hotelbeds hotels search endpoint.
- **Configuration**: Supports configuration via command line flags or environment variables.
- **Flexible Port and Host Configuration**: Customize application port and Hotelbeds API host.
- **Security**: Handles Hotelbeds API key and secret securely.

## Installation

### Prerequisites

- Go (version 1.22 or higher)
- Docker (optional, for containerized deployment)

### Clone Repository

```bash
git clone https://github.com/jeffy-pro/lite-api.git
cd lite-api
```

### Build and Run Locally
```bash
go build -o lite-api main.go
./lite-api start --port=:8080 --host=https://api.test.hotelbeds.com --apikey=<yourapikey> --secret=<yoursecret>
```

### Docker Installation (Optional)
```bash
docker build -t lite-api .
docker run -p 8080:8080 -e APP_PORT=:8080 -e HOTELBEDS_HOST=https://api.test.hotelbeds.com -e HOTELBEDS_API_KEY=yourapikey -e HOTELBEDS_SECRET=yoursecret lite-api start

```

## Usage

### Command Line Flags

You can configure the Lite-API using command line flags:
```bash
./lite-api start --port=:8080 --host=https://api.test.hotelbeds.com --apikey=yourapikey --secret=yoursecret
```

#### Flags

* -p, --port: Specify the application port (default is :8080).
* -o, --host: Specify the Hotelbeds API host URL (default is https://api.test.hotelbeds.com).
* -k, --apikey: Specify the Hotelbeds API key.
* -s, --secret: Specify the Hotelbeds API secret.

#### Environment Variables
Alternatively, you can set configuration values using environment variables:
```bash
export APP_PORT=:8080
export HOTELBEDS_HOST=https://api.test.hotelbeds.com
export HOTELBEDS_API_KEY=<yourapikey>
export HOTELBEDS_SECRET=<yoursecret>
./lite-api start
```

### Help
For more information on available commands and their options, use the help command:

## Contributing
Contributions are welcome! If you find any issues or have suggestions, please open an issue or a pull request on GitHub.

## License
This project is licensed under the MIT License. See the LICENSE file for details.