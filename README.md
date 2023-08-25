# GitHub S3 cache

Allows to save and restore files on S3

## Usage

### Save files

```yml
- name: Save node_modules
  uses: dealstack/gh-actions-s3-cache/save@v1
  if: steps.node-modules-cache.outputs.cache-hit != 'true'
  with:
    path: ./node_modules
    key: ${{ runner.os }}-node-modules-${{ hashFiles('**/pnpm-lock.yaml') }}

```

### Restore files

```yml
- name: Restore node_modules if exist
  uses: dealstack/gh-actions-s3-cache/restore@v1
  id: node-modules-cache
  with:
    path: ./node_modules
    fail-on-cache-miss: false
    key: ${{ runner.os }}-node-modules-${{ hashFiles('**/pnpm-lock.yaml') }}
```
