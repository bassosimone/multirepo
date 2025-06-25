# Design

The `multirepo` tool allows to manage several distinct git
repositories like they were inside a monorepo.

Here's the general idea:

1. `multirepo init` to create an empty multirepo in a directory.

2. `multirepo clone` to clone a repository in the multirepo directory
and track it for subsequent commands.

3. `multirepo repo add` to add an existing repository in the
multirepo directory to the multirepo and track it.

4. `multirepo repo rm` to remove a repository from the multirepo index.

5. `multirepo foreach` to execute a command in each repository.


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


## `multirepo clone [-vx] <repo>`

Clones a repository into the multirepo.

Flags:

- `-v`: show the executed commands ouput.

- `-x`: prints executed commands.

For example:

```bash
multirepo clone git@github.com:ooni/probe-cli
```

This command implements the following steps:

1. Locks the `.multirepo` directory using the `.multirepo/lock` file.

2. Clones the repository into the current directory.

3. Updates the configuration file `.multirepo/config.json`.


## `multirepo foreach [-kx] <command> [args...]`

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


## `multirepo repo add <dir>`

Adds an existing repository in the current directory to the multirepo.

For example:

```bash
multirepo repo add probe-cli
```

This command implements the following steps:

1. Locks the `.multirepo` directory using the `.multirepo/lock` file.

2. Executes `git config --get remote.origin.url` in `<dir>` to obtain the SSH URL.

3. Updates the configuration file `.multirepo/config.json`.


## `multirepo repo rm <dir>`

Removes a repository from the multirepo index without touching
the existing repository directory.

For example:

```bash
multirepo repo rm probe-cli
```

This command implements the following steps:

1. Locks the `.multirepo` directory using the `.multirepo/lock` file.

2. Updates the configuration file `.multirepo/config.json`.
