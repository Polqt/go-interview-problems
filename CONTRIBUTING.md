# Contributing Guide

Thank you for considering contributing to **go-interview-problems**! This repository is a collection of **brain teasers** and **real-life programming tasks** to help developers prepare for interviews. Your contributions help make this resource even more valuable.

## How to Contribute

You can contribute by:
- Adding a new interview problem.
- Improving existing problems.
- Refining solutions or providing alternative approaches.
- Fixing typos or improving documentation.

### Adding a New Problem

To add a new problem, follow these steps:

1. **Create a new folder** in the repository's root directory using the format:
```md
<problem-number>-<problem-name>/
```

Example:
```md
03-first-successful-key-lookup/
```

2. **Inside the problem folder, add the following files**:
```bash
03-first-successful-key-lookup/
├── README.md         # Problem description, tags, and source
├── task.go           # Starter code for the problem
├── solution/
│   ├── solution.go   # One or more possible solutions
```

3. **Fill out the `README.md` using this template**:
```md
# Name
<Short description of the problem>

## Tags
<tag1> <tag2>

## Source
- <Original source link, if applicable>
```

4. **Write the starter code (`task.go`)** 
The starter code should include function signatures but not the full implementation. Example:
```go
package main

import "context"

func Get(ctx context.Context, addresses []string, key string) (string, error) {
    return "", nil
}
```

5. **Add at least one solution in the `solution/` directory**
Provide a complete implementation in `solution/solution.go`.

6. **Commit your changes and submit a pull request (PR)**

* Ensure your code follows standard Go formatting (`gofmt`).
* Keep problem descriptions **clear and concise**.
* Provide **meaningful commit messages**.

## Best Practices
* Focus on **real-life tasks** and **brain teasers**, rather than purely algorithmic challenges.
* Use **clear and consistent formatting**.
* If your problem is inspired by an external source (e.g., a YouTube video), always **credit the original author** in the Source section.
* Keep **problem difficulty balanced** - it should be challenging yet solvable.
