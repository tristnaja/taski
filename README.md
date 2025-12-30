<div align="center">
  <h1 style="font-size: 4rem; font-weight: bold;">TASKI</h1>
  <p>A lightweight and efficient CLI task management application built with Go.</p>
</div>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Badge">
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge" alt="License Badge">
</p>

---

## 🚀 About TASKI

TASKI is a command-line task management tool designed for simplicity and efficiency. It allows you to quickly add, view, modify, and manage tasks directly from your terminal. With features like soft delete and automatic cleanup, TASKI keeps your task list organized without losing important data prematurely.

## ✨ Features

-   ✅ **Quick Task Management:** Add, view, change, and delete tasks effortlessly.
-   🗑️ **Soft Delete:** Deleted tasks are retained for 30 days before permanent removal.
-   ♻️ **Restore Functionality:** Recover individual or all deleted tasks.
-   🧹 **Automatic Cleanup:** Old deleted tasks are automatically purged after 30 days.
-   💾 **JSON Storage:** Lightweight file-based storage using JSON.

## 🛠️ Tech Stack

| Category      | Technology                                                                 |
| :------------ | :------------------------------------------------------------------------- |
| **Language**  | [Go](https://go.dev/)                                                      |
| **Storage**   | JSON file-based storage                                                    |
| **Testing**   | Go standard testing package                                                |

## 📂 Project Structure

```
taski/
├── app/
│   ├── cmd/              # Command implementations
│   │   ├── add.go
│   │   ├── change.go
│   │   ├── delete.go
│   │   ├── restore.go
│   │   └── view.go
│   └── taski/
│       └── main.go       # Application entry point
├── internal/
│   └── io/
│       └── io.go         # Data persistence layer
├── tests/                # Test files
├── go.mod
└── README.md
```

## 🏁 Getting Started

### Prerequisites

-   [Go](https://go.dev/dl/) (v1.18 or higher)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/tristnaja/taski.git
    cd taski
    ```

2.  **Build the application:**
    ```sh
    go build -o taski ./app/taski
    ```

3.  **Run the application:**
    ```sh
    ./taski <command> [options]
    ```

### Usage

#### Add a New Task
```sh
taski add --title "Task Title" --desc "Task Description"
# or using shorthand
taski add -t "Task Title" -d "Task Description"
```

#### View All Tasks
```sh
taski view
```

#### Change an Existing Task
```sh
taski change --index <task_id> --title "New Title" --desc "New Description"
```

#### Delete a Task
```sh
taski delete --index <task_id>
```

#### Restore a Task
```sh
# Restore a specific task
taski restore --mode single --index <task_id>

# Restore all deleted tasks
taski restore --mode all
```

## 📋 Commands

| Command    | Description                                    |
| :--------- | :--------------------------------------------- |
| `add`      | Add a new task with title and description      |
| `view`     | Display all active tasks                       |
| `change`   | Modify an existing task                        |
| `delete`   | Soft delete a task (retains for 30 days)       |
| `restore`  | Restore deleted task(s)                        |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have ideas for improvements or find any bugs.

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

---

> **Note:** This README and tests are AI generated because im lazy.
