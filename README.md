# sensor-am2302-data-forwarder

The application sensor-am2302-data-forwarder can be used to forward data of a [sensor-am2302-history](https://github.com/fi3te/sensor-am2302-history) server to different destinations, e.g. AWS API Gateway or ntfy notification service.

## Usage

1. Update the `config.yml` file to your needs:
   - `interval-in-seconds`: Specifies the interval in seconds at which data is forwarded.
   - `first-retry-after-error-in-seconds`: Sets the period of time after which the first retry is carried out after an error. As long as the cause of the error has not been eliminated, exponential backoff is used as retry behavior.
   - `source-directory`: Directory with files of the [sensor-am2302-history](https://github.com/fi3te/sensor-am2302-history) server. 
   - `file-determination-by-date`: Specifies how to determine the file from which to forward data. If the current date is not to be used, the file is determined based on the order of the file names.
   - `retention-period-in-hours`: A time to live timestamp is calculated from this value, which can be interpreted by the destination.
   - `aws`/`http`/`ntfy`: The currently supported destinations are AWS REST APIs, HTTP endpoints and ntfy notification servers.
2. Update the docker setup if necessary.
3. Start the application using docker compose (recommended)
   ```
   docker compose up -d
   ```
   or run it without docker.
   ```
   go run .\cmd\main.go
   ```
