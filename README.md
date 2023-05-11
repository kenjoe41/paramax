# Paramax

Paramax is a command-line program for analyzing URL parameters. It provides both passive and active modes of operation to help identify potential security vulnerabilities and discover hidden functionality within a target domain.

## Features

- Passive mode: Analyzes the target domain by fetching and processing URLs from various aggregators.
- Active mode: Performs active analysis by modifying URL parameters and generating new URLs for testing.
- Support for subdomains: Includes subdomains when fetching URLs from aggregators (optional).
- Exclude specific file extensions from analysis.
- Output results to a file.
- Customizable placeholder string for parameter modification.
- Silent mode: Suppresses printing results to the screen when an output file is specified.

## Installation

### With GIT Install
Direct from github with go:

    ```
    go install github.com/kenjoe41/paramax/...@latest
    ```
### Manually
1. Clone the Paramax repository:

    ```shell
    git clone https://github.com/kenjoe41/paramax.git
    ```

2. Navigate to the project directory:

    ```shell
    cd paramax
    ```

3. Build the binary using the Go compiler:

    ```
    go build .
    ```

## Usage
Passive mode (default):

    ```
    paramax --domain example.com
    ```

Active mode:

    ```
    paramax active --domain example.com
    ```

For more options and flags, refer to the command-line help:

    ```
    paramax --help
    ```

## Credit
I started out and rewrote [ParamSpider](https://github.com/devanshbatham/ParamSpider) by [0xAsm0d3us](https://twitter.com/0xAsm0d3us), into Golang. All credit to him.
## Contributing
Contributions are welcome! If you find any bugs or have suggestions for improvements, please submit an issue or create a pull request.

## License
This project is licensed under the [MIT License](https://www.mit.edu/~amini/LICENSE.md).