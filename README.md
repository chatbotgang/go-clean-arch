# Go-Clean-Architecture-Template

This repo introduces a Go clean architecture template commonly used in Crescendo Lab. We are going to explain how the architecture work through a tutorial on building a sample application - **Crescendo Barter**.

# Crescendo Barter

Crescendo Barter is a second-hand goods exchange application in which people can post their old goods and exchange them with others.

### 1. User Stories

- Account management
    - As a client, I want to register a trader account.
    - As a client, I want to log in to the application through the registered trader account.
- Second-hand Goods
    - As a trader, I want to post my old goods to the application so that others can see what I have.
    - As a trader, I want to see all my posted goods.
    - As a trader, I want to see others’ posted goods.
    - As a trader, I want to remove some of my goods from the application.
- Goods Exchange
    - As a trader, I want to exchange my own goods with others.
    
### 2. Project Dependencies

* [Golang](https://go.dev): ^1.17
* [gin](https://github.com/gin-gonic/gin): ~1.7.7
* [zerolog](https://github.com/rs/zerolog): ~1.26.1
* [sqlx](https://github.com/jmoiron/sqlx): ~1.3.4
* [PostgreSQL](https://www.postgresql.org/docs/13/index.html): 13


### 3. Development Guideline

See [development guideline](./docs/development-guideline.md).
