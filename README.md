# Billing-Engine Service

The **Billing-Engine Service** is a backend service responsible for handling loan management, payment processing, and calculating outstanding loans for users. It provides key business logic for the financial operations of a loan system, including querying loan details, checking overdue payments, and processing user payments.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Configuration](#configuration)
- [Setup and Installation](#setup-and-installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Postman Collection](#postman-collection)
- [Sequence Flow Diagram](#sequence-flow-diagram)
- [Tracing](#tracing)

## Overview

The **Billing-Engine Service** is designed to interact with a database to fetch user loan details, track overdue schedules, and process payments. The service uses a transaction-based approach to ensure consistency and rollback in case of errors. Additionally, it integrates with **APM (Application Performance Monitoring)** to track performance and trace the flow of requests.

### Key Components:
- **UserService**: Handles business logic for user loan management.
- **Database Interaction**: Uses GORM for database operations (e.g., querying users, loans, schedules).
- **APM Integration**: Monitors the performance and flow of requests.
- **Payment Processing**: Allows users to make payments, updating their loan balance and schedules.

## Features
- Fetch outstanding loan details and overdue schedules.
- Check if a user is delinquent based on overdue schedules.
- Process payments, update loans, and create payment records.
- Transaction management with commit/rollback support.
- Performance monitoring using APM.

## Configuration

- **Main Configuration**: `config.yaml`

### `config.yaml`

This file contains the main configuration for the service, including logging levels, APM (Application Performance Monitoring), Datastore settings, Billing engine delinquentThreshold. All config.yaml replacable by .env file if provided

## Setup and Installation

Follow these steps to set up and run the **Billing-Engine Service** locally or on your preferred server.

### Prerequisites:
- Docker
- Docker Compose (for containerized setup)
- GORM (for database interactions)
- PostgreSQL or another supported relational database
- APM tool (e.g., Elastic APM, Datadog) for performance monitoring (optional but recommended)


### Run without docker:
1. **Clone the repository**:
   ```bash
   git clone https://github.com/abialemuel/billing-engine.git
   cd billing-engine

2. Install the dependencies:

    ```sh
    go mod tidy
    ```

3. Run the service:

    ```sh
    go run main.go
    ```

### Run with docker-compose (recommended)
Highly recommend using Docker Compose for easier setup, as it will automatically handle the service dependencies, including the PostgreSQL database migrations and data seeding.

1. Run all container:

    ```sh
    docker compose up -d --build
    ```

## Usage
### API Endpoints
The Billing-Engine service exposes the following endpoints to interact with the loan and payment data:

1. GET /v1/billing-engine/user/outstanding
Get the outstanding loan details, overdue schedules and delinquent status for the authenticated user.

Example Request:

    ```sh
    curl --location 'http://localhost:8080/v1/billing-engine/user/outstanding' \
    --header 'Authorization: Basic YWJpYWxlbXVlbDphYmlhbGVtdWVs'
    ```
    
Example Response:

    ```sh
    {
        "code": 200,
        "message": "Success",
        "payload": {
            "loan_id": 1,
            "username": "abialemuel",
            "outstanding": 5500000,
            "is_delinquent": false,
            "upcoming_amount": 0,
            "missed_payment": 0
        }
    }
    ```

2. POST /v1/billing-engine/user/payment
Process a payment for the authenticated user. The payment must match the overdue amount.

Example Request:

    ```sh
    curl --location 'http://localhost:8080/v1/billing-engine/user/payment' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic YWJpYWxlbXVlbDphYmlhbGVtdWVs' \
    --data '{
        "amount": 220000
    }'
    ```

Example Response if already paid/no overdue:

    ```sh
    {
        "code": 400,
        "status": "VALIDATION_ERROR",
        "message": "no overdue schedule found"
    }
    ```

Example Response if success:

    ```sh
    {
        "code": 200,
        "status": "SUCCESS",
        "message": "Success"
    }
    ```

## Postman Collection

You can find the Postman collection for the Billing-Engine Service in the [docs/Billing Engine.postman_collection.json](docs/Billing%20Engine.postman_collection.json) file.

## Sequence Flow Diagram

You can find the sequence flow diagrams for each endpoint in the following files:
- [Get Outstanding Loan](docs/sequence/getOutstandingLoan.png)
- [Make Payment](docs/sequence/makePayment.png)

## Tracing

If you want to get trace information, you can check Jaeger APM at [http://localhost:16686](http://localhost:16686).

![Jaeger APM](docs/jaeger.png)