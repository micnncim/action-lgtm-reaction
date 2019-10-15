# LGTM Reaction

[![CI](https://github.com/micnncim/action-lgtm-reaction/workflows/CI/badge.svg)](https://github.com/micnncim/action-lgtm-reaction/actions)
[![Release](https://img.shields.io/github/v/release/micnncim/action-lgtm-reaction.svg?logo=github)](https://github.com/micnncim/action-lgtm-reaction/releases)
[![Marketplace](https://img.shields.io/badge/marketplace-lgtm--reaction-blue?logo=github)](https://github.com/marketplace/actions/lgtm-reaction)

![](docs/assets/screen-record.gif)

Send LGTM reaction as image or GIF when we say `lgtm`.  
Currently supports [GIPHY](https://giphy.com).

## Usage

### Create Workflow

The default trigger comment is `lgtm`, but you can specify any trigger comment with `jobs.<job_id>.steps.with`.  
The trigger comment match a comment in case-insensitive.
It means we can use the both of `lgtm` and `LGTM` as a trigger (sure, `Lgtm` and `lGtm` can be used).

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
          trigger: '[".*looks good to me.*"]' # default: '["^(lgtm|LGTM)$", "^[gG]ood [jJ]ob!?$"]'
          override: true # default: false
```

## License

[MIT License](LICENSE)
