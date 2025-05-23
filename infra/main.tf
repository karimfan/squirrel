terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 2.20"
    }
  }
}

provider "docker" {}

resource "docker_network" "squirrel" {
  name = "squirrel_net"
}

resource "docker_image" "postgres" {
  name = "postgres:15"
}

resource "docker_container" "db" {
  name  = "squirrel_db"
  image = docker_image.postgres.latest

  env = [
    "POSTGRES_USER=postgres",
    "POSTGRES_PASSWORD=postgres",
    "POSTGRES_DB=squirrel"
  ]

  networks_advanced {
    name = docker_network.squirrel.name
    aliases = ["db"]
  }

  ports {
    internal = 5432
    external = 5432
  }
}

resource "docker_image" "server" {
  name = "squirrel-server"
  build {
    context    = "${path.module}/.."
    dockerfile = "${path.module}/Dockerfile"
  }
}

resource "docker_container" "server" {
  name  = "squirrel_server"
  image = docker_image.server.latest

  depends_on = [docker_container.db]

  env = [
    "DATABASE_URL=postgres://postgres:postgres@db:5432/squirrel?sslmode=disable"
  ]

  networks_advanced {
    name    = docker_network.squirrel.name
    aliases = ["server"]
  }

  ports {
    internal = 8080
    external = 8080
  }
}
