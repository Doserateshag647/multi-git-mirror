# Authentication

Guide for configuring authentication for each supported Git provider.

<br/>

## Table of Contents

- [GitLab](#gitlab)
- [GitHub](#github)
- [Bitbucket](#bitbucket)
- [AWS CodeCommit](#aws-codecommit)
- [SSH Key](#ssh-key)
- [Security Best Practices](#security-best-practices)

<br/>

## GitLab

Uses OAuth2 token authentication. The token is injected into the HTTPS URL as `https://oauth2:<token>@gitlab.com/...`.

```yaml
- uses: somaz94/git-mirror-action@v1
  with:
    targets: |
      gitlab::https://gitlab.com/myorg/myrepo.git
    gitlab_token: ${{ secrets.GITLAB_TOKEN }}
```

### Required Token Scopes

- `write_repository` — Push access to the target repository
- `read_repository` — (optional) If the target repo is private

### Creating a GitLab Token

1. Go to **GitLab** > **Settings** > **Access Tokens**
2. Create a token with `write_repository` scope
3. Add the token as a GitHub secret (`GITLAB_TOKEN`)

<br/>

## GitHub

Uses x-access-token authentication. The token is injected as `https://x-access-token:<token>@github.com/...`.

```yaml
- uses: somaz94/git-mirror-action@v1
  with:
    targets: |
      github::https://github.com/myorg/myrepo.git
    github_token: ${{ secrets.MIRROR_GITHUB_TOKEN }}
```

### Required Token Scopes

- `repo` — Full control of private repositories (or `public_repo` for public repos)

> **Note**: Do not use `${{ secrets.GITHUB_TOKEN }}` for cross-repository mirroring. It only has access to the current repository. Use a Personal Access Token (PAT) instead.

<br/>

## Bitbucket

Uses username and app password authentication. Credentials are injected as `https://user:password@bitbucket.org/...`.

```yaml
- uses: somaz94/git-mirror-action@v1
  with:
    targets: |
      bitbucket::https://bitbucket.org/myorg/myrepo.git
    bitbucket_username: ${{ secrets.BITBUCKET_USERNAME }}
    bitbucket_password: ${{ secrets.BITBUCKET_APP_PASSWORD }}
```

### Creating a Bitbucket App Password

1. Go to **Bitbucket** > **Personal settings** > **App passwords**
2. Create a password with **Repositories: Write** permission
3. Add the username and app password as GitHub secrets

> **Note**: Both `bitbucket_username` and `bitbucket_password` must be provided. If either is missing, the URL is used as-is.

<br/>

## AWS CodeCommit

CodeCommit uses IAM-based authentication via the Git credential helper. The URL is used as-is without token injection.

```yaml
- uses: somaz94/git-mirror-action@v1
  with:
    targets: |
      codecommit::https://git-codecommit.us-east-1.amazonaws.com/v1/repos/myrepo
```

### IAM Configuration

Ensure the runner has AWS credentials configured with `codecommit:GitPush` permissions. For GitHub Actions, use `aws-actions/configure-aws-credentials`:

```yaml
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v4
  with:
    aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
    aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    aws-region: us-east-1

- name: Setup CodeCommit credential helper
  run: git config --global credential.helper '!aws codecommit credential-helper $@'

- uses: somaz94/git-mirror-action@v1
  with:
    targets: |
      codecommit::https://git-codecommit.us-east-1.amazonaws.com/v1/repos/myrepo
```

<br/>

## SSH Key

For providers that support SSH, use the `ssh_private_key` input with the `generic` provider.

```yaml
- uses: somaz94/git-mirror-action@v1
  with:
    targets: |
      generic::git@custom-git.example.com:org/repo.git
    ssh_private_key: ${{ secrets.SSH_PRIVATE_KEY }}
```

### How It Works

When `ssh_private_key` is provided, the action automatically:

1. Writes the key to `/root/.ssh/mirror_key` with `0600` permissions
2. Creates an SSH config with `StrictHostKeyChecking no`
3. Sets `GIT_SSH_COMMAND` to use the configured key
4. Cleans up all SSH files after mirroring completes

### Setup

1. Generate an SSH key pair: `ssh-keygen -t ed25519 -C "mirror-action"`
2. Add the public key to the target Git server
3. Add the private key as a GitHub secret (`SSH_PRIVATE_KEY`)

<br/>

## Security Best Practices

- **Always use GitHub Secrets** — Never hardcode tokens or passwords in workflow files
- **Use least-privilege tokens** — Grant only the minimum required scopes
- **Rotate tokens regularly** — Set expiration dates and rotate periodically
- **Use separate tokens per target** — Avoid reusing a single token across providers
- **Audit access** — Regularly review which tokens have access to your repositories

### Built-in Security Features

- **Credential URL encoding** — Special characters (`@`, `:`, `/`, etc.) in passwords are URL-encoded to prevent URL parsing issues
- **Debug log masking** — Tokens and passwords are replaced with `***` in debug logs
- **SSH key cleanup** — SSH key files are automatically removed after mirroring
- **Config validation** — Warns when required credentials are missing for a target provider
