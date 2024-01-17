# Cadre CLI 🚀

![Golang](https://img.shields.io/badge/Go-00add8.svg?labelColor=171e21&style=for-the-badge&logo=go)

![Build](https://github.com/kmesiab/cadre/actions/workflows/go-build.yml/badge.svg)
![Lint](https://github.com/kmesiab/cadre/actions/workflows/go-lint.yml/badge.svg)
![Test](https://github.com/kmesiab/cadre/actions/workflows/go-test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/kmesiab/cadre)](https://goreportcard.com/report/github.com/kmesiab/cadre)

## Overview 🌟

Cadre CLI is a command-line application designed to automate code
reviews across various programming languages, utilizing OpenAI's
ChatGPT API. It offers intelligent insights and suggestions to
improve code quality and developer efficiency.

## Features 🛠️

- **Language-Agnostic Analysis**: Compatible with multiple programming
languages.
- **AI-Powered Insights**: Employs ChatGPT for in-depth code analysis.
- **User-Friendly CLI**: Simple and intuitive command-line interface
for easy usage.

## Installation 🔧

To install Cadre CLI, you need to have Go installed on your machine.
Follow these steps:

```bash
go install github.com/kmesiab/cadre@latest
```

Set your OpenAI API Key:

```bash
export OPENAI_API_KEY=sk-[SECRET]
```

Set your Ignore Files. These are file types that will be excluded from
code reviews. They should be a comma-separated list of file extensions.
For example:

```bash
export IGNORE_FILES=.mod,.sum
```

To run the program:

```bash
cadre
```

## Usage 💡

**Usage instructions for Cadre CLI go here. Provide examples and explain
how users can interact with it.**

## Development and Testing 🧪

### Building the Project 🏗️

```bash
make build
```

### Running Tests ✔️

```bash
make test
make test-verbose
make test-race
```

### Installing Tools 🛠️

```bash
make install-tools
```

### Linting 🧹

```bash
make lint
make lint-markdown
```

---

## Contributing 🤝

### Forking and Sending a Pull Request

1. **Fork the Repository**: Click the 'Fork' button at the top right of this
   page.
2. **Clone Your Fork**:

   ```bash
   git clone https://github.com/kmesiab/cadre
   cd cadre
   ```

3. **Create a New Branch**:

   ```bash
   git checkout -b your-branch-name
   ```

4. **Make Your Changes**: Implement your changes or fix issues.
5. **Commit and Push**:

   ```bash
   git commit -m "Add your commit message"
   git push origin your-branch-name
   ```

6. **Create a Pull Request**: Go to your fork on GitHub and click the
   'Compare & pull request' button.

## Github Guidelines

Please ensure your code adheres to the project's
[standards and guidelines](https://github.com/kmesiab/ai-code-critic/discussions/24).

### Quick Tips

Run `make lint` before committing to ensure your code is properly formatted.

1. **Always rebase, never merge commit**
2. Always use a description commit message
3. Separate your title from your description
4. Keep commit messages under 50 characters
5. Start your branch with `feat|bugfix|docs|style|refactor|perf|test`
6. Squash your commits into logical units of work

## License 📝

Information regarding the licensing of cadre will be included here.

---

*Note: This project is under active development. Additional features
and documentation will be updated in due course.* 🌈
