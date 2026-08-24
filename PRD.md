# Module Manager (MM)

## Product Requirements Document (PRD)

### 1. Product Overview

**Module Manager (MM)** is a Windows desktop application for managing Node.js project directories from a single interface.

The application will allow users to provide a parent folder containing multiple project folders. MM will detect the projects and provide quick actions to:

* View project information.
* See the frameworks/libraries used by each project.
* Remove `node_modules`.
* Install `node_modules`.
* Delete framework/cache files.
* Search and filter projects/cache files.
* View Git/GitHub history and the date of the latest commit.
* Select multiple projects for deletion.
* View total storage used by the managed folder.
* See ongoing operations through a terminal-like process panel.
* Queue multiple installation operations.

**Application name:** Module Manager
**Short name:** MM
**Platform:** Windows
**Frontend:** Vite + React
**Backend:** Go

---

# 2. Goals

The application will provide a single-page interface for quickly managing the dependencies and cache files of multiple Node.js projects stored inside a selected folder.

The primary goals are:

1. Quickly remove `node_modules` from projects.
2. Quickly reinstall project dependencies.
3. Identify the framework/libraries used by each project.
4. Identify when a project was last committed to Git.
5. Help determine whether old `node_modules` are worth keeping based on project activity.
6. Find and remove framework-generated cache files.
7. Support multiple project operations.
8. Queue installation operations rather than running all installations simultaneously.
9. Provide clear visibility into currently running and queued operations.
10. Show the total space consumed by the selected parent folder.

---

# 3. Target User

The application is intended for developers who have many local Node.js projects and want a quick way to manage the storage consumed by:

* `node_modules`
* Framework caches
* Other project cache files

---

# 4. Application Structure

MM will be a **single-page application**.

The main screen will contain:

1. Folder selection area.
2. Storage information.
3. Search/filter controls.
4. Project list.
5. Project actions.
6. Multi-selection controls.
7. Cache-file management.
8. Git activity information.
9. Terminal-like process/output panel.

---

# 5. Folder Selection

The user will be able to select a folder from Windows.

The selected folder will be treated as the **parent project directory**.

Example:

```text
Projects/
├── project-a/
├── project-b/
├── project-c/
├── project-d/
└── project-e/
```

Each relevant subfolder will be treated as a project/repository.

MM will scan the selected folder and display the projects it contains.

---

# 6. Storage Information

At the top of the application, MM will display the **total space currently occupied by the selected parent folder**.

Example:

```text
Projects
Total Space Used: 38.6 GB
```

The value should represent the storage consumed by the selected folder and its contents.

---

# 7. Project List

Each detected project will be displayed as a row.

Each project row will contain:

* Project/folder name.
* Framework/library information.
* Git activity information.
* `node_modules` information.
* Install action.
* Selection control for bulk operations.

Example conceptual layout:

```text
┌───────────────────────────────────────────────────────────────┐
│ project-a                                      [Install]      │
│ React · Vite · TypeScript                                     │
│ Last commit: 3 days ago                                       │
│ node_modules: 842 MB                                          │
└───────────────────────────────────────────────────────────────┘
```

The exact UI design can be finalized during implementation.

---

# 8. Framework and Library Detection

MM will inspect each project to determine the frameworks and libraries being used.

The information will be displayed as a label below the project/folder name.

Example:

```text
my-react-app
React · Vite · TypeScript
```

Another project could display:

```text
backend-api
Express · Node.js · MongoDB
```

The detected information should be based on the project's dependency/configuration information.

---

# 9. Node Modules Management

## 9.1 Delete `node_modules`

MM will provide an option to remove the `node_modules` directory from an individual project.

The purpose is to quickly reclaim disk space without manually navigating through project directories.

---

## 9.2 Install Dependencies

An **Install** button will be displayed directly in the project row.

Example:

```text
project-a                         [Install]
React · Vite
```

Clicking **Install** will start the dependency installation process for that project.

The installation process will use the project's existing package-management configuration.

---

# 10. Multiple Project Selection

Users will be able to select multiple projects.

The selected projects can then be used for bulk deletion of `node_modules`.

Example:

```text
☑ project-a
☑ project-b
☐ project-c
☑ project-d
```

After selecting projects, the user can choose the delete action.

---

# 11. Delete Confirmation

Deleting `node_modules` or other selected files will require confirmation.

When the user chooses delete, MM will display a popup confirmation.

Example:

```text
Delete node_modules?

3 projects are selected.

This will permanently remove their node_modules
directories.

             [Cancel]    [Delete]
```

The deletion will only begin after confirmation.

---

# 12. Git History / Last Commit

MM will inspect the Git repository associated with each project.

The application will determine the latest commit and display how long ago it occurred.

Example:

```text
Last commit: 3 days ago
```

Possible values:

```text
Last commit: 2 hours ago
Last commit: 3 days ago
Last commit: 4 months ago
Last commit: 1 year ago
```

This information will allow the user to understand how recently a project was active.

MM will use this Git activity information when displaying project status.

---

# 13. GitHub History

Where applicable, MM will also retrieve/display the project's GitHub-related history information.

The purpose is to provide visibility into the project's latest development activity.

The relevant repository information will be associated with the local project.

---

# 14. Cache File Management

MM will also provide functionality for finding framework/library-generated cache files.

Some frameworks and development tools generate cache directories/files that can consume significant disk space.

The application will provide an option to search for these cache files.

Example:

```text
Cache Files

project-a
  .vite/
  .cache/

project-b
  .next/cache/

project-c
  node_modules/.cache/
```

The user will be able to identify these files/directories and delete them.

---

# 15. Cache Search

There will be a dedicated option to search for cache files within the selected parent folder.

The search will identify relevant framework/tool cache directories.

The results will show enough information for the user to understand what is being removed.

---

# 16. Search

MM will provide a search field for finding projects.

Example:

```text
Search projects...
```

Searching:

```text
react
```

could show:

```text
react-dashboard
react-admin
react-website
```

---

# 17. Filtering

The project list will support filtering.

Filtering will allow the user to narrow down the displayed projects based on the available project information.

The filtering functionality will remain limited to the information already provided by MM, such as project/framework/library information and project activity.

---

# 18. Process Panel

MM will contain a small **terminal-like process/output panel** showing the current work being performed by the application.

This is not intended to be a large timeline.

It should visually resemble a compact terminal/output area.

Example:

```text
> Preparing to delete node_modules...
> Project: ecommerce-dashboard
> Removing node_modules...
> Removed 1.24 GB
> Completed.
```

For installation:

```text
> Installing dependencies...
> Project: portfolio
> Running package installation...
> Installing packages...
> Completed.
```

The panel will provide visibility into the current operation.

---

# 19. Operation Queue

If the user starts installation for multiple projects, installations will be processed **sequentially**.

Example:

```text
Install requested:

1. project-a   → Running
2. project-b   → Queued
3. project-c   → Queued
4. project-d   → Queued
```

Once `project-a` finishes:

```text
1. project-a   → Completed
2. project-b   → Running
3. project-c   → Queued
4. project-d   → Queued
```

The process panel will display the progress of the currently executing operation.

---

# 20. Visual Design

The application will use a **dark gray desktop-app aesthetic inspired by macOS applications**.

The interface should be:

* Dark.
* Clean.
* Minimal.
* Developer-oriented.
* Compact.
* Easy to scan.
* Focused on project management rather than decorative UI.

The terminal-like process panel should fit naturally into this dark interface.

---

# 21. Technology Stack

### Frontend

```text
Vite
React
```

The frontend will provide the single-page user interface.

### Backend

```text
Go
```

The Go backend will handle filesystem-related and project-management operations, including:

* Folder scanning.
* Project detection.
* Storage calculation.
* `node_modules` detection.
* `node_modules` deletion.
* Dependency installation.
* Cache detection.
* Cache deletion.
* Git history/activity retrieval.
* Operation queue processing.

### Platform

```text
Windows
```

---

# 22. Core User Flow

### Initial flow

```text
Open MM
   ↓
Select parent folder
   ↓
Scan folder
   ↓
Detect projects
   ↓
Display project list
   ↓
Display framework/library labels
   ↓
Display Git activity
   ↓
Display storage information
```

### Delete flow

```text
Select projects
   ↓
Choose Delete
   ↓
Confirmation popup
   ↓
Confirm
   ↓
Deletion begins
   ↓
Terminal-like panel shows process
   ↓
Operation completes
   ↓
Project/storage information updates
```

### Install flow

```text
Click Install
   ↓
Project enters installation queue
   ↓
If no operation is running → start immediately
   ↓
Otherwise → remain queued
   ↓
Terminal-like panel displays progress
   ↓
Installation completes
   ↓
Next queued project starts
```

### Cache flow

```text
Search for cache files
   ↓
Display detected cache files
   ↓
Select cache files
   ↓
Delete
   ↓
Confirmation popup
   ↓
Delete selected cache files
   ↓
Display operation in process panel
```

---

# 23. MVP Scope

The current version of Module Manager will contain **only** the following functionality:

* Windows folder selection.
* Detection of project subfolders.
* Project listing.
* Framework/library labels.
* Project search.
* Project filtering.
* Total folder storage calculation.
* `node_modules` detection.
* Individual `node_modules` installation.
* Multiple project selection.
* Multiple `node_modules` deletion.
* Delete confirmation popup.
* Git history/latest commit information.
* GitHub-related history information where applicable.
* Cache-file search.
* Cache-file deletion.
* Terminal-like current-process display.
* Sequential installation queue.
* Dark gray macOS-inspired interface.
* Vite + React frontend.
* Go backend.

No additional product features are included in the current PRD.
