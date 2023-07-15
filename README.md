# Quake Log Parser

Quake Log Parser is a Go-based project that provides an HTTP interface for parsing Quake 3 Arena game log files. Users can upload log files and the application will parse the games in the logs and return structured JSON data detailing each game's players, scores, and kills.

## Project Structure

Here's a brief outline of the project structure:

- `cmd/`: Contains the main application that starts the server.
- `handler/`: Contains HTTP handler functions for handling requests.
- `models/`: Contains data models for the game data.
- `parser/`: Contains logic for parsing the log files.
- `report/`: Contains logic for generating the JSON reports from parsed log files.
- `static/`: Contains static HTML files served by the server.
- `testdata/`: Contains test data and expected results for unit tests.

Each of the `cmd`, `handler`, `models`, `parser`, and `report` directories also contain respective `_test.go` files for testing their functionalities.

## Building and Running using Docker

This project uses Docker for easy setup and deployment. You can build and run the project using Docker.

Build the Docker image:

```bash
docker build -t quake-parser .
```

This will produce an executable named `quake-parser` in your project root.

## Running

You can start the server by running the executable produced by the build step:

```bash
./quake-parser
```

Once the server is running, you can access it at `http://localhost:8080`.

## Usage

To parse a log file, send a POST request to `/upload` with the log file attached as a form file under the key "qgames". The server will respond with a JSON representation of the games contained in the log.

```bash
curl -X POST -F "qgames=@/path/to/your/logfile.log" http://localhost:8080/upload
```

The server will return a JSON representation of the games contained in the log.

## Testing

To run the unit tests, navigate to the root directory of the project and use the go test command:

```bash
go test -v ./...
```

This will run all the tests in the project and display the output.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue if you have something to contribute.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
