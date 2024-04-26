# Contributing to DepGraph

We love your input! We want to make contributing to DepGraph as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Pull Request Process

1. Update the README.md with details of changes to the interface, if applicable.
2. Update the documentation with any new dependencies, configuration options, or features.
3. The PR will be merged once you have the sign-off of at least one other developer.

## Any contributions you make will be under the MIT Software License

In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using GitHub's [issue tracker](https://github.com/yourusername/DepGraph/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/yourusername/DepGraph/issues/new).

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can.
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Use a Consistent Coding Style

* Run `go fmt` before committing
* Use meaningful variable names
* Comment your code when necessary
* Follow Go best practices and conventions

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as contributors and maintainers pledge to making participation in our project and our community a harassment-free experience for everyone.

### Our Standards

Examples of behavior that contributes to creating a positive environment include:

* Using welcoming and inclusive language
* Being respectful of differing viewpoints and experiences
* Gracefully accepting constructive criticism
* Focusing on what is best for the community
* Showing empathy towards other community members

### Our Responsibilities

Project maintainers are responsible for clarifying the standards of acceptable behavior and are expected to take appropriate and fair corrective action in response to any instances of unacceptable behavior.

## Testing

Before submitting a pull request, make sure all tests pass:

```bash
go test ./...
```

## Documentation

We use Go's standard documentation format. Please document all exported functions, types, and packages.

Example:
```go
// FunctionName does something specific.
// It takes important parameters and returns useful results.
func FunctionName(param string) (result string, err error) {
    // Implementation
}
```

## Project Structure

When adding new features, please follow our project structure:

```
DepGraph/
├── cmd/          # Command-line tools
├── pkg/          # Public packages
│   ├── analysis/ # Dependency analysis
│   ├── github/   # GitHub API client
│   ├── graph/    # Graph data structure
│   └── web/      # Web interface
├── internal/     # Private packages
└── docs/         # Documentation
```

## License

By contributing, you agree that your contributions will be licensed under its MIT License. 