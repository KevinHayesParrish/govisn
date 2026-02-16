# Release Process for GoVisn

This document describes the process for creating and releasing new versions of GoVisn.

## Version Numbering

This project uses [Semantic Versioning](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., 0.23.1)
  - **MAJOR**: Breaking changes
  - **MINOR**: New features (backward compatible)
  - **PATCH**: Bug fixes and patches

## Release Checklist

### 1. Prepare the Code

- [X] Ensure all tests pass: `make test`
- [X] Review all changes and commits
- [X] Update code documentation if needed
- [X] Verify the application builds: `make build`
- [X] Update the version constant in `main.go` (GOVISN_VERSION)

### 2. Update Version Files

- [X] Update `version.txt` with new version number (e.g., 0.24.0)
- [X] Update `CHANGELOG.md` with changes in the new version
  - Add a new section with the version and date
  - List added features, bug fixes, breaking changes, and improvements

### 3. Create the Release

- [X] Run `make release` to commit, tag, and prepare the release
- [X] Verify the git tag was created: `git tag -l`

### 4. Build Release Binaries

- [ ] Run `make release-build` to create multi-platform binaries
- [ ] Verify binaries were created:
  - govisn-VERSION-linux-amd64
  - govisn-VERSION-darwin-amd64
  - govisn-VERSION-darwin-arm64
  - govisn-VERSION-windows-amd64.exe
- [ ] Copy data folder to each binary location (required for fonts and images):
  
  ```bash
  cp -r data govisn-VERSION-linux-amd64/
  cp -r data govisn-VERSION-darwin-amd64/
  cp -r data govisn-VERSION-darwin-arm64/
  cp -r data govisn-VERSION-windows-amd64/
  ```

### 5. Push to Repository

```bash
# Push the new version tag
git push origin v<VERSION>

# Or push with the main branch
git push origin main
git push origin --tags
```

### 6. Create GitHub Release (optional)

- [ ] Go to GitHub repository
- [ ] Navigate to "Releases"
- [ ] Click "Create a new release"
- [ ] Select the tag you just created
- [ ] Upload the built binaries with data folders
- [ ] Add release notes from CHANGELOG.md
- [ ] Publish the release

## Important Notes for GoVisn

1. **Data Folder**: GoVisn requires the `data/` folder (containing fonts and images) to be distributed with the executable. Include it in release packages.

2. **Version Constant**: Update `GOVISN_VERSION` in `main.go` to match the version in `version.txt`.

3. **Cross-Platform**: GoVisn uses G3N which has platform-specific rendering. Test binaries on target platforms when possible.

## Quick Release Commands

```bash
# Update version.txt and CHANGELOG.md manually first, then:
make release

# Build multi-platform binaries
make release-build

# Push to repository
git push origin v<VERSION>
```

## Example Release Session

```bash
# 1. Update GOVISN_VERSION constant in main.go
# 2. Edit version.txt: change to 0.24.0
# 3. Edit CHANGELOG.md: add new section for 0.24.0
# 4. Commit code changes
git add main.go CHANGELOG.md version.txt
git commit -m "Update for release 0.24.0"

# 5. Run release process
make release

# 6. Build release binaries
make release-build

# 7. Add data folders to binaries
cp -r data govisn-0.24.0-*/

# 8. Push to repository
git push origin v0.24.0

# 9. Verify
git tag -l
ls -la govisn-0.24.0-*
```

## Prerequisites for Building

- Go 1.12 or higher
- Required dependencies from `go.mod` (run `go mod download`)
- For cross-compilation: cgo for sqlite3 support may require additional tools

## Continuous Integration

Consider setting up GitHub Actions for:

- Running tests on every push
- Automatically building release binaries on tag creation
- Publishing releases to GitHub automatically

See `.github/workflows/` for examples (if added in the future).
