# **Resilience in Distributed Systems**

<table><tbody><tr><td><a href="https://github.com/kinfinity/distributed-resilience/actions/workflows/build.yaml"><img src="https://github.com/kinfinity/distributed-resilience/actions/workflows/build.yaml/badge.svg" alt="Build &amp; Tests"></a></td><td><a href="https://github.com/kinfinity/distributed-resilience/actions/workflows/codeql.yaml"><img src="https://github.com/kinfinity/distributed-resilience/actions/workflows/codeql.yaml/badge.svg" alt="CodeQL"></a></td></tr></tbody></table>

Distributed Systems need to be able to gracefully handle failures and recover from them. This is achieved through resilience, which involves designing systems while anticipating scenarios where nodes/services/resources over which the system is distributed may fail to be accessed or behave unexpected due to

- Slow Networks
- Network Timeouts
- Overcommited/Overloaded resources or Services
- Temporarily unavailable resources or service
- Partial loss of connectivity

## **Patterns**

This repository covers several patterns implemented in Golang which have been designed to handle resilience in distributed environments.

- [Timeout](./timeout/readme.md)
- [Retry](./retry/readme.md)
- [Circuit Breaker](./circuitbreaker/readme.md)

# **References**

- [Release It! Second Edition - Stability Patterns by Michael Nygard](https://pragprog.com/titles/mnee2/release-it-second-edition/)
- [codecentric resilience-design-patterns](https://www.codecentric.de/wissens-hub/blog/resilience-design-patterns-retry-fallback-timeout-circuit-breaker)
- [The Resilience Patterns your Microservices Teams Should Know by Victor Rentea](https://youtu.be/IR89tmg9v3A?si=w96y4S6AbVt_CviB)
