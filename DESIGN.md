# Design

The `multirepo` tool allows to manage several distinct git
repositories like they were inside a monorepo.


## `multirepo init`

Creates an empty multirepo in the current directory.

For example:

```bash
multirepo init [-x]
```

Flags:

- `-x`: prints executed commands.

This command implements the following steps:

1. Creates the `.multirepo` directory if it does not exist.

2. Locks the `.multirepo` directory using the `.multirepo/lock` file.

3. Creates the default configuration file `.multirepo/config.json`
if it does not exist.


## `multirepo add [-vx] {repo}`

Adds a repository to the multirepo.

Flags:

- `-v`: show the executed commands ouput.

- `-x`: prints executed commands.

For example:

```bash
multirepo add git@github.com:ooni/probe-cli
```

This command implements the following steps:

1. Locks the `.multirepo` directory using the `.multirepo/lock` file.

2. Clones the repository into the current directory.

3. Updates the configuration file `.multirepo/config.json`.


## `multirepo foreach [-kx] {command} [args...]`

Executes a command in each repository.

Flags:

- `-k`: keep running in case of failure.

- `-x`: prints executed commands.

For example:

```bash
multirepo foreach -kx -- git pull
```

This command implements the following steps:

1. Reads the configuration file `.multirepo/config.json`.

2. Sets the `MULTIREPO_ROOT` environment variable to the current directory.

3. Sets the `MULTIREPO_EXECUTABLE` environment variable to the path of
the `multirepo` executable.

4. Executes the given `command` in each repository.
