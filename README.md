# Krateo MLFlow Plugin for `rest-dynamic-controller`

## Overview

The Krateo MLFlow Plugin for `rest-dynamic-controller` is a simple wrapper around the MLFlow API designed for the Krateo Operator Generator (KOG). This plugin provides a set of endpoints to interact with MLFlow, making it easier to manage experiments and other MLFlow resources.

## Summary

- [Krateo MLFlow Plugin for `rest-dynamic-controller`](#krateo-mlflow-plugin-for-rest-dynamic-controller)
  - [Overview](#overview)
  - [Summary](#summary)
  - [API](#api)
  - [Examples](#examples)
    - [Get Experiment Metadata](#get-experiment-metadata)
    - [Get Run Metadata](#get-run-metadata)
  - [Configuration](#configuration)
  - [Swagger Documentation](#swagger-documentation)

## API

1. **Get Experiment Metadata**
  - **Endpoint:** `/2.0/mlflow/experiments/get`
  - **Description:** Retrieves metadata for a specified experiment. The endpoint requires the `experiment_id` as a query parameter and returns the experiment's metadata in the response body.

2. **Get Run Metadata**
  - **Endpoint:** `/2.0/mlflow/runs/get`
  - **Description:** Retrieves metadata for a specified run. The endpoint requires the `run_id` as a query parameter and returns the run's metadata in the response body.

## Examples

### Get Experiment Metadata

To get metadata for an experiment, send a GET request to the `/2.0/mlflow/experiments/get` endpoint with the `experiment_id` query parameter.

```bash
curl -X GET "http://localhost:8080/2.0/mlflow/experiments/get?experiment_id=1" -H "accept: application/json"
```

### Get Run Metadata

To get metadata for a run, send a GET request to the `/2.0/mlflow/runs/get` endpoint with the `run_id` query parameter.

```bash
curl -X GET "http://localhost:8080/2.0/mlflow/runs/get?run_id=1" -H "accept: application/json"
```

## Configuration

To configure the plugin, set the following environment variables:

- `PLUGIN_DEBUG`: Enable or disable debug mode (default: true).
- `PLUGIN_PORT`: Port to listen on (default: 8081).
- `MLFLOW_SERVER`: URL of the MLFlow server (default: http://localhost:5000).

## Swagger Documentation

The plugin provides a Swagger UI to explore and test the available API endpoints. This can be accessed at the `/swagger` endpoint.

- **Endpoint:** `/swagger`
- **Description:** Opens the Swagger UI, which provides a web-based interface to interact with the API. This is useful for testing and understanding the available endpoints and their required parameters.

To access the Swagger UI, navigate to `http://localhost:8080/swagger` in your web browser.
````# Krateo MLFlow Plugin for `rest-dynamic-controller`

## Overview

The Krateo MLFlow Plugin for `rest-dynamic-controller` is a simple wrapper around the MLFlow API designed for the Krateo Operator Generator (KOG). This plugin provides a set of endpoints to interact with MLFlow, making it easier to manage experiments and other MLFlow resources.

## Summary

- [Krateo MLFlow Plugin for `rest-dynamic-controller`](#krateo-mlflow-plugin-for-rest-dynamic-controller)
  - [Overview](#overview)
  - [Summary](#summary)
  - [API](#api)
  - [Examples](#examples)
    - [Get Experiment Metadata](#get-experiment-metadata)
    - [Get Run Metadata](#get-run-metadata)
  - [Configuration](#configuration)
  - [Swagger Documentation](#swagger-documentation)

## API

1. **Get Experiment Metadata**
  - **Endpoint:** `/2.0/mlflow/experiments/get`
  - **Description:** Retrieves metadata for a specified experiment. The endpoint requires the `experiment_id` as a query parameter and returns the experiment's metadata in the response body.

2. **Get Run Metadata**
  - **Endpoint:** `/2.0/mlflow/runs/get`
  - **Description:** Retrieves metadata for a specified run. The endpoint requires the `run_id` as a query parameter and returns the run's metadata in the response body.

## Examples

### Get Experiment Metadata

To get metadata for an experiment, send a GET request to the `/2.0/mlflow/experiments/get` endpoint with the `experiment_id` query parameter.

```bash
curl -X GET "http://localhost:8080/2.0/mlflow/experiments/get?experiment_id=1" -H "accept: application/json"
```

### Get Run Metadata

To get metadata for a run, send a GET request to the `/2.0/mlflow/runs/get` endpoint with the `run_id` query parameter.

```bash
curl -X GET "http://localhost:8080/2.0/mlflow/runs/get?run_id=1" -H "accept: application/json"
```

## Configuration

To configure the plugin, set the following environment variables:

- `PLUGIN_DEBUG`: Enable or disable debug mode (default: true).
- `PLUGIN_PORT`: Port to listen on (default: 8081).
- `MLFLOW_SERVER`: URL of the MLFlow server (default: http://localhost:5000).

## Swagger Documentation

The plugin provides a Swagger UI to explore and test the available API endpoints. This can be accessed at the `/swagger` endpoint.

- **Endpoint:** `/swagger`
- **Description:** Opens the Swagger UI, which provides a web-based interface to interact with the API. This is useful for testing and understanding the available endpoints and their required parameters.

To access the Swagger UI, navigate to `http://localhost:8080/swagger` in your web browser.
````