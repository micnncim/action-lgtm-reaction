![logo](docs/assets/logo.png)

[![actions-workflow-CI][actions-workflow-CI-badge]][actions-workflow-CI]
[![actions-marketplace][actions-marketplace-badge]][actions-marketplace]
[![release][release-badge]][release]
[![license][license-badge]][license] 

![](docs/assets/screen-record.gif)

Send LGTM reaction as image when we say `lgtm`.  

Currently supports [LGTM.app](https://www.lgtm.app) and [GIPHY](https://giphy.com).

## Usage

### Create Workflow

#### `jobs.<job_id>.steps.env`

|       Key       |             Value              |                        Required                        |
| --------------- | ------------------------------ | ------------------------------------------------------ |
| `GITHUB_TOKEN`  | `${{ secrets.GITHUB_TOKEN }}`  | `true`                                                 |
| `GIPHY_API_KEY` | `${{ secrets.GIPHY_API_KEY }}` | `true` if `jobs.<job_id>.steps.with.source` == `giphy` |

#### `jobs.<job_id>.steps.with`

|    Key     |               Default                | Required |                             Note                              |
| ---------- | ------------------------------------ | -------- | ------------------------------------------------------------- |
| `trigger`  | `'["^lgtm$", "^[gG]ood [jJ]ob!?$"]'` | `false`  | Trigger comment body. It must be JSON string array of regexp. |
| `override` | `false`                              | `false`  | Override posted comment body or not.                          |
| `source`   | `lgtmapp`                            | `false`  | `lgtmapp` or `giphy`                                          |

#### Example

For minimalists:

```yaml
name: Send LGTM reaction
on:
  issue_comment:
    types: [created]
  pull_request_review:
    types: [submitted]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@1.0.0
      - uses: micnncim/action-lgtm-reaction@
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

For nerds:

```yaml
name: Send LGTM reaction
on:
  issue_comment:
    types: [created]
  pull_request_review:
    types: [submitted]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: micnncim/action-lgtm-reaction@master # Set some version.
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GIPHY_API_KEY: ${{ secrets.GIPHY_API_KEY }}
        with:
          trigger: '[".*looks good to me.*"]'
          override: true
          source: 'giphy'
```

### Setting GIPHY

The default image source is [LGTM.app](https://www.lgtm.app) but you can also configure [GIPHY](https://giphy.com).

1. Create your app and get API key from [here](https://developers.giphy.com/dashboard).
2. Set the API key in GitHub repository (Setting > Secret) as `GIPHY_API_KEY`.
3. Configure `'giphy'` in your GitHub Actions workflow.

## Projects using `action-lgtm-reaction`

- [Cake Website](https://github.com/cake-build/website)

## Note

*Icon made by Freepik from [www.flaticon.com](https://www.flaticon.com)*

<!-- badge links -->

[actions-workflow-CI]: https://github.com/micnncim/action-lgtm-reaction/actions?query=workflow%3ACI
[actions-marketplace]: https://github.com/marketplace/actions/lgtm-reaction
[release]: https://github.com/micnncim/action-lgtm-reaction/releases
[license]: LICENSE

[actions-workflow-CI-badge]: https://img.shields.io/github/workflow/status/micnncim/action-lgtm-reaction/CI?label=CI&style=for-the-badge&logo=github
[actions-marketplace-badge]: https://img.shields.io/badge/marketplace-lgtm%20reaction-blue?style=for-the-badge&logo=github
[release-badge]: https://img.shields.io/github/v/release/micnncim/action-lgtm-reaction?style=for-the-badge&logo=github
[license-badge]: https://img.shields.io/github/license/micnncim/action-lgtm-reaction?style=for-the-badge
