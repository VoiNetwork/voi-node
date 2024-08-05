group "default" {
  targets = ["build", "test", "all"]
}

target "build" {
  platforms = ["linux/amd64", "linux/arm64"]
  tags      = ["voi-node"]
  context   = "."
  dockerfile = "Dockerfile"
}

target "test" {
  inherits = ["build"]
  args = {
    TARGET = "test"
  }
  commands = [
    "make test"
  ]
}

target "all" {
  inherits = ["build"]
  args = {
    TARGET = "all"
  }

  commands = [
    "make all"
  ]
}