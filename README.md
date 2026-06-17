# World of Transport

CLI tool that finds transport hubs within a given distance of a location using
the [haversine formula](https://en.wikipedia.org/wiki/Haversine_formula) and IBM
Cloudant's `airportdb`.

## Usage

```sh
./wot <latitude> <longitude> <distance_km>
```

**Examples:**

```sh
./wot 47.5 19.0 50      # hubs within 50 km of Budapest
./wot 51.5 -0.1 30      # hubs within 30 km of central London
```

Output:

```
Transport hubs within 50.0 km of (47.5000, 19.0000):

Name              Latitude    Longitude   Distance
----              --------    ---------   --------
Budapest Keteli   47.500497   19.085484   6.42 km
Ferihegy          47.436933   19.255592   20.45 km
```

### Go binary

```sh
go build -o wot ./src
./wot 47.5 19.0 50

# unit tests only (no network)
go test -short ./tests/

# all tests (requires network)
go test ./tests/
```

### Docker

```sh
docker build -t wot .
docker run --rm wot 47.5 19.0 50

# tests (uses builder stage which has Go)
docker build --target builder -t wot-test .
docker run --rm wot-test go test -short ./tests/
docker run --rm wot-test go test ./tests/
```

## Configuration

Cloudant settings are overridable via environment variables:
`CLOUDANT_URL` (default `https://mikerhodes.cloudant.com`),
`CLOUDANT_DB` (`airportdb`), `CLOUDANT_DESIGN_DOC` (`view1`),
`CLOUDANT_INDEX` (`geo`).

```sh
docker run --rm -e CLOUDANT_URL="https://other.cloudant.com" wot 47.5 19.0 50
```
