# Project Modules

Based on the initial design for Bart Warn, the project is broken down into the following technical units.

## 1. API / Protobuf Contracts (`api`)
- **Responsibility:** Defines the source-of-truth `.proto` file(s) for the system.
- **Details:** Contains the gRPC service definition (e.g., `LocationPing`) and message types. This module will be used to generate both the Go server code and the Android client stubs.

## 2. Server - Entrypoint & Coordination (`server-grpc`)
- **Responsibility:** The main Go application and gRPC server.
- **Details:** Implements the gRPC handlers. When a ping is received, this module acts as the orchestrator: it logs the event, calls the routing engine to get a decision, and then passes that decision to the SMS notification client.

## 3. Server - BART Routing Engine (`routing-engine`)
- **Responsibility:** Interacts with the official BART API and applies custom routing logic.
- **Details:** Fetches real-time departures for the pinged station. It calculates the fastest route to El Cerrito Plaza and applies the preference rule: *If the optimal route is within 5 minutes of a Red line train, choose the Red line. Otherwise, choose the fastest route.*

## 4. Server - SMS Notification Client (`sms-client`)
- **Responsibility:** Handles outbound communication via SMS.
- **Details:** Integrates with the Twilio Go SDK. It takes the output from the routing engine, formats it into a concise, readable string (e.g., "Take Red line 17:05"), and dispatches it to the user's phone.

## 5. Client - Android Location Tracker (`android-tracker`)
- **Responsibility:** The background app running on the user's phone.
- **Details:** Monitors location against hardcoded geographic boundaries for BART stations. Acts as a gRPC client to send a ping to the server exactly when the user enters one of these zones.
