# Git Sync CLI App

This is a command-line interface (CLI) application written in Go that helps synchronize Git repositories between different origins.

### Installation

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/napisani/git-sync-go.git
   ```

2. Change to the project directory:

   ```bash
   cd git-sync-go
   ```

3. Build the application:

   ```bash
   make
   ```

### Usage

To run the Git Sync CLI app, you need to provide a configuration file as a command-line argument. The configuration file should be in JSON format and specify the repositories you want to synchronize.

```bash
./git-sync <config-file>
```

For example, if your configuration file is named `config.json`, you can run the app with the following command:

```bash
./git-sync config.json
```

Make sure to adjust the configuration file according to your needs before running the application.

## Configuration

The configuration file should follow the structure defined in the `SyncConfig` struct in the `main.go` file. It consists of the following properties:

- `tempDirectory` (optional): The temporary directory where the repositories will be cloned during synchronization. If not provided, a default temporary directory will be used.

- `fromToConfigs`: An array of `FromToConfig` objects that specify the synchronization details for each repository. Each `FromToConfig` object has the following properties:

  - `fromOrigin`: The source origin of the repository.

  - `toOrigin`: The target origin where the repository should be synchronized.

  - `branches`: An array of branch names that should be synchronized.

  - `force`: A boolean indicating whether to force-push the changes to the target origin.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributions

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or create a pull request.

