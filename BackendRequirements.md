| Area         | Backend method        | Purpose                                                              |
| ------------ | --------------------- | -------------------------------------------------------------------- |
| Folder       | `Select/ScanFolder`   | Receive the selected parent folder and scan its immediate subfolders |
| Projects     | `GetProjects`         | Return detected projects and their basic information                 |
| Project info | `GetProjectInfo`      | Get framework/libraries, package manager, `node_modules` size/status |
| Storage      | `GetFolderSize`       | Calculate total space used by the selected folder                    |
| Node modules | `GetNodeModulesInfo`  | Check whether `node_modules` exists and its size                     |
| Node modules | `DeleteNodeModules`   | Delete `node_modules` for one project                                |
| Installation | `InstallDependencies` | Run the appropriate package-manager install command                  |
| Queue        | `QueueInstallation`   | Add installation requests to the sequential queue                    |
| Cache        | `FindCacheFiles`      | Search projects for known framework/tool cache directories           |
| Cache        | `DeleteCacheFiles`    | Delete selected cache directories/files                              |
| Git          | `GetGitInfo`          | Get current branch, latest commit and commit age                     |
| GitHub       | `GetGitHubInfo`       | Get repository information where available                           |
| Process      | `GetProcessStatus`    | Send current operation/status to frontend                            |
| Process      | `GetProcessLogs`      | Stream terminal-like output to frontend                              |



| Frontend API                            | Go backend function     | Purpose                           |
| --------------------------------------- | ----------------------- | --------------------------------- |
| `POST /api/folders/scan`                | `ScanFolder()`          | Scan selected parent folder       |
| `GET /api/projects/:id`                 | `GetProjectInfo()`      | Get detailed project information  |
| `GET /api/folder/size`                  | `GetFolderSize()`       | Get total folder size             |
| `DELETE /api/projects/:id/node-modules` | `DeleteNodeModules()`   | Delete `node_modules`             |
| `POST /api/projects/:id/install`        | `InstallDependencies()` | Install dependencies              |
| `GET /api/cache`                        | `FindCacheFiles()`      | Find cache files                  |
| `DELETE /api/cache`                     | `DeleteCacheFiles()`    | Delete selected caches            |
| `GET /api/projects/:id/git`             | `GetGitInfo()`          | Get latest Git commit information |
