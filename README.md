# LGTM Reaction

[![CI](https://github.com/micnncim/action-lgtm-reaction/workflows/CI/badge.svg)](https://github.com/micnncim/action-lgtm-reaction/actions)
[![Release](https://img.shields.io/github/v/release/micnncim/action-lgtm-reaction.svg?logo=github)](https://github.com/micnncim/action-lgtm-reaction/releases)
[![Marketplace](https://img.shields.io/badge/marketplace-lgtm--reaction-blue?logo=github)](https://github.com/marketplace/actions/lgtm-reaction)

![](docs/assets/screen-record.gif)

Send LGTM reaction as image or GIF when we say `lgtm`.  

Currently supports [LGTM.app](https://www.lgtm.app/) and [GIPHY](https://giphy.com).

## Usage

### Create Workflow

#### `jobs.<job_id>.steps.env`

|             Key              |                   Value                   |                         Required                          |
| ---------------------------- | ----------------------------------------- | --------------------------------------------------------- |
| `GITHUB_TOKEN`               | `${{ secrets.GITHUB_TOKEN }}`             | `true`                                                    |
| `GIPHY_API_KEY`              | `${{ secrets.GIPHY_API_KEY }}`            | `true` if `jobs.<job_id>.steps.with.source` == `giphy`    |
| `GITHUB_REPOSITORY`          | `${{ github.repository }}`                | `true`                                                    |
| `GITHUB_ISSUE_NUMBER`        | `${{ github.event.issue.number }}`        | `true` if `on.issue_comment.types` == `[created]`         |
| `GITHUB_COMMENT_BODY`        | `${{ github.event.comment.body }}`        | `true` if `on.issue_comment.types` == `[created]`         |
| `GITHUB_COMMENT_ID`          | `${{ github.event.comment.id }}`          | `true` if `on.issue_comment.types` == `[created]`         |
| `GITHUB_PULL_REQUEST_NUMBER` | `${{ github.event.pull_request.number }}` | `true` if `on.pull_request_review.types` == `[submitted]` |
| `GITHUB_REVIEW_BODY`         | `${{ github.event.review.body }}`         | `true` if `on.pull_request_review.types` == `[submitted]` |
| `GITHUB_REVIEW_ID`           | `${{ github.event.review.id }}`           | `true` if `on.pull_request_review.types` == `[submitted]` |

#### `jobs.<job_id>.steps.with`

|    Key     |               Default                | Required |         Note         |
| ---------- | ------------------------------------ | -------- | -------------------- |
| `trigger`  | `'["^lgtm$", "^[gG]ood [jJ]ob!?$"]'` | `false`  |                      |
| `override` | `false`                              | `false`  |                      |
| `source`   | `lgtmapp`                            | `false`  | `lgtmapp` or `giphy` |

`jobs.<job_id>.steps.with.trigger` should be an JSON string array of regexp.

#### Example

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
      - uses: micnncim/action-lgtm-reaction@latest
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GIPHY_API_KEY: ${{ secrets.GIPHY_API_KEY }}
          GITHUB_REPOSITORY: ${{ github.repository }}
          GITHUB_ISSUE_NUMBER: ${{ github.event.issue.number }}
          GITHUB_COMMENT_BODY: ${{ github.event.comment.body }}
          GITHUB_COMMENT_ID: ${{ github.event.comment.id }}
          GITHUB_PULL_REQUEST_NUMBER: ${{ github.event.pull_request.number }}
          GITHUB_REVIEW_BODY: ${{ github.event.review.body }}
          GITHUB_REVIEW_ID: ${{ github.event.review.id }}
        with:
          trigger: '[".*looks good to me.*"]'
          override: true
          source: 'giphy'
```

## Projects using this action

- [Cake Website](https://github.com/cake-build/website)

## License

[MIT License](LICENSE)
