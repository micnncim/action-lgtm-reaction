# LGTM Reaction

Send LGTM reaction as image or GIF when we say `lgtm`.  
Currently supports [GIPHY](https://giphy.com).

## Usage

### Create Workflow

```yaml
name: Send LGTM reaction
on:
  issue_comment:
    types: [created]
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
        with:
          trigger: ['Looks Good To Me'] # default: lgtm, LGTM
```
