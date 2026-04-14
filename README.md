This project is licensed under the CC BY-NC-SA 4.0 License
---

# Nexa 
**A minimalist, keyboard-centric TUI task manager for power users.**

Nexa aims to solve the bloat of modern to-do software by providing a fast, tree-based architecture that lives entirely in your terminal. It focuses on automatic ranking, intuitive path-based command entry, and cross-platform persistence.

Nexa removes the friction of deciding where to start.  
Instead of scanning a to-do list, tasks are automatically ranked so you can immediately begin working.

Built for power users who prefer speed, keyboard control, and zero bloat.

> **Note to Developers:** This was built as a learning project in Go.  
The code is functional but not optimised for extensibility. Explore at your own risk!

---

<img width="1351" height="337" alt="nexa_demo_1" src="https://github.com/user-attachments/assets/7b7ffd32-ae5d-4058-852a-93c3f5f02b7d" />

---

## 🧠 How it works

1. Add tasks
2. Select a project or category
3. Nexa ranks tasks automatically
4. Start with the top task

No planning. No reordering. Just execution.

---

## ✨ Features
* **Automatic Ranking:** Tasks are automatically sorted by priority and urgency.
* **Hierarchical Logic:** Organise your life into Projects → Categories → Tasks.
* **Path-Based Commands:** Manipulate data using intuitive paths like `Work/ProjectA/Finish_Docs`.
* **Dynamic Deadlines:** Set dates using shorthand like `1_d` (tomorrow) or `2_w` (two weeks).
* **Safety First:** Integrated "Lock/Unlock" system to prevent accidental deletion of large projects.
* **Portable Data:** Your data is stored in a simple JSON file at `~/.nexa/nexa_data.json`, making it easy to sync across devices.

---

## 🏗 Architecture & Logic
Nexa follows a strict three-tier hierarchy:
1.  **Project:** The top-level container.
2.  **Category:** A subset within a project.
3.  **Task:** The actionable item (must reside within a category).

### Navigation
* **Tree View (Left):** Use `Shift + Up/Down` to navigate projects and categories.
* **Task View (Right):** Use `Shift + Left/Right` to browse tasks.
* **Highlights:** The currently active item is highlighted in **blue**.

---

## ⌨️ Usage & Commands
Nexa uses a **semicolon-delimited** syntax for commands: `command;arg1;arg2`.

### Path Syntax
You can refer to items using `/` or `-`.
* **Full Path:** `Education/Math/Homework`
* **Category Path:** `Education/Math`
* **Project Path:** `Education`
* **Shortcuts:** Use `.` or `*` to refer to the **currently highlighted** item.
    * Example: `*/English/Essay` refers to the English category in the project you are currently hovering over.

### Essential Commands

| Command | Syntax | Description |
| :--- | :--- | :--- |
| **`mk`** | `mk;path;priority;date;time` | **Make** a project, category, or task. |
| **`rm`** | `rm;path` | **Remove** an item (Recursive! Requires `ulk`). |
| **`mv`** | `mv;oldPath;newPath` | **Move** or rename an item. |
| **`dn`** | `dn;path` (or `Ctrl+D`) | Mark task as **Done** (Undo with `Ctrl+Z`). |
| **`rp`** | `rp;path;priority;Y M W D;date` | Create a **Repeating** task. |
| **`rs`** | `rs;path;date` | **Reschedule** an existing task. |
| **`ulk`** | `ulk` | **Unlock** the app to allow deletions. |
| **`lk`** | `lk` | **Lock** the app for safety. |
| **`q`** | `q` or `quit` | **Exit** Nexa. |
| **`p`** | `p;path;newPriority` | Change **priority** of a task. |

---

## 📅 Date & Priority Formatting

### Priorities
Integers from **1 to ∞**. 1 is the lowest priority; higher numbers rank higher in your list.

### Date Shorthand
When creating or rescheduling tasks, you can use absolute dates (`DD-MM-YYYY`) or dynamic shorthand:
* `x_d`: *x* days from today
* `x_w`: *x* weeks from today
* `x_m`: *x* months from today
* `x_y`: *x* years from today

### Repeating Tasks (`rp`)
The third argument defines the interval: `Years Months Weeks Days`.
* `0 0 1 0`: Repeats every week.
* `0 1 0 1`: Repeats every month and 1 day.

---

## 💾 Persistence
Nexa saves your data automatically before and after every command.
* **Location:** `~/.nexa/nexa_data.json`
* **Cross-Compatibility:** You can move this JSON file between Windows, Linux, and macOS without issues.

---

## 🛡 Deletion Prevention
To protect your data, `rm` commands are disabled by default. 
1.  Type `ulk` to unlock.
2.  Perform your `rm` command.
3.  Type `lk` to re-lock.



-----

## 🚀 Quick Start Guide

Nexa is a standalone binary. No installation or dependencies are required. Just download the version for your system and run it.

### 1\. Download the Binary

Download the appropriate file for your machine from the [Releases](https://github.com/CurryEshay/nexa/releases) page:

| OS | Architecture | Binary Name |
| :--- | :--- | :--- |
| **Linux** | AMD64 (Standard) | `nexa-linux-amd64` |
| **Windows** | AMD64 (Standard) | `nexa-windows-amd64.exe` |
| **macOS** | Apple Silicon (M1/M2/M3) | `nexa-darwin-arm64` |
| **macOS** | Intel | `nexa-darwin-amd64` |

-----

### 2\. Installation Instructions

#### 🐧 Linux

1.  Open your terminal.
2.  Give the file execution permissions:
    ```bash
    chmod +x nexa-linux-amd64
    ```
3.  Run it:
    ```bash
    ./nexa-linux-amd64
    ```
    *(Optional)* Move it to your path to run it from anywhere: `sudo mv nexa-linux-amd64 /usr/local/bin/nexa`

#### 🪟 Windows

1.  Open **PowerShell** or **Command Prompt**.
2.  Navigate to your Downloads folder:
    ```powershell
    cd ~/Downloads
    ```
3.  Run the executable:
    ```powershell
    ./nexa-windows-amd64.exe
    ```

#### 🍎 macOS

1.  Open Terminal.
2.  Grant execution permissions:
    ```bash
    chmod +x nexa-darwin-arm64
    ```
3.  **Note:** On first run, macOS may block the app because it is from an "Unidentified Developer."
      * Go to **System Settings \> Privacy & Security**.
      * Scroll down and click **"Open Anyway"** for Nexa.
4.  Run it:
    ```bash
    ./nexa-darwin-arm64
    ```

-----

### 3\. Basic Commands

Once the TUI (Terminal User Interface) is open, use these keyboard-driven commands:

  * **`mk;[Project]/[Category]/[Task];priority;date;time`**: Create a new task.
  * **`dn;[Project]/[Category]/[Task]`**: Mark a task as done.
  * **`Shift + up/down`**: Navigate the tree view.
  * * **`Shift + up/down`**: Navigate the tree view.
  * * * **`Shift + left/right`**: Navigate the task view.
  * **`q` or `Ctrl+C`**: Quit Nexa.

-----

### 🛠 Troubleshooting

  * **Data Location:** Nexa saves your tasks in `~/.nexa/nexa_data.json` on Linux/Mac, and `C:\Users\YourName\.nexa\nexa_data.json` on Windows.
  * **Nil Pointer Error:** If the app crashes on first run, ensure your user has permission to create folders in the Home directory.

-----

