# **Resilience in Distributed Systems**

[![BUILD & TESTS](https://github.com/kinfinity/distributed-resilience/actions/workflows/build.yaml/badge.svg)](https://github.com/kinfinity/distributed-resilience/actions/workflows/build.yaml)

Distributed Systems need to be able to gracefully handle failures and recover from them. This is achieved through resilience, which involves designing systems while anticipating scenarios where nodes/services/resources over which the system is distributed may fail to be accessed or behave unexpected due to

- Slow Networks
- Network Timeouts
- Overcommited/Overloaded resources or Services
- Temporarily unavailable resources or service
- Partial loss of connectivity

## **Patterns**

This repository covers several patterns implemented in Golang which have been designed to handle resilience in distributed environments.

- [Retry](./retry/readme.md)
- [Circuit Breaker](./circuitbreaker/readme.md)

# **References**

- [The Resilience Patterns your Microservices Teams Should Know by Victor Rentea](https://youtu.be/IR89tmg9v3A?si=w96y4S6AbVt_CviB)
