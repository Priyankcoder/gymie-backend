
# GitHub Branch Protection Setup Guide

## Making the Backend CI Check Required Before Merge

Follow these steps to require the "Backend CI" workflow to pass before any PR can be merged into `main`.

### Step 1: Go to Repository Settings

1. Go to your GitHub backend repository: `https://github.com/YOUR-USERNAME/Gymie-backend`
2. Click on **Settings** (top right)
3. In the left sidebar, click **Branches**

### Step 2: Add Branch Protection Rule

1. Under "Branch protection rules", click **Add rule** (or **Edit** if you already have a rule for `main`)

2. **Branch name pattern**: Enter `main`

3. **Check these boxes:**

   - ✅ **Require a pull request before merging**
     - ✅ Require approvals: `1` (or more if you have a team)
     - ✅ Dismiss stale pull request approvals when new commits are pushed
   
   - ✅ **Require status checks to pass before merging**
     - ✅ Require branches to be up to date before merging
     - In the search box, type: `Backend CI` or `build-and-test`
     - ✅ Select **Backend CI / build-and-test** from the dropdown
     - ✅ Select **Backend CI / lint** (optional but recommended)
   
   - ✅ **Require conversation resolution before merging** (optional)
   
   - ✅ **Do not allow bypassing the above settings** (recommended)

4. Click **Create** or **Save changes**

### Step 3: Verify It's Working

1. Create a new branch: `git checkout -b test-branch`
2. Make a change to any backend file
3. Push and create a PR
4. You should see:
   - ✅ "Backend CI / build-and-test" check running
   - ✅ "Backend CI / lint" check running
   - 🔒 **Merge button disabled** until checks pass

### Visual Guide

```
GitHub Repository Settings
├── Branches
│   └── Branch protection rules
│       └── Rule for: main
│           ├── ✅ Require pull request reviews
│           ├── ✅ Require status checks (SELECT THESE):
│           │   ├── ✅ Backend CI / build-and-test
│           │   └── ✅ Backend CI / lint
│           └── ✅ Do not allow bypassing
```

### What the Workflow Checks

The `Backend CI` workflow now checks:

#### **build-and-test job:**
- ✅ Go dependencies are valid
- ✅ Code formatting (`go fmt`)
- ✅ Static analysis (`go vet`)
- ✅ Code builds successfully
- ✅ API binary compiles
- ✅ All tests pass
- ✅ Race conditions detected
- ✅ Code coverage generated

#### **lint job:**
- ✅ Code quality (`golangci-lint`)
- ✅ Best practices
- ✅ Common mistakes
- ✅ Security issues

### Additional Protection Options (Recommended)

#### Protect Against Force Push
- ✅ **Restrict who can push to matching branches**
  - Add specific users/teams who can push directly

#### Require Linear History
- ✅ **Require linear history**
  - Forces rebase instead of merge commits (cleaner git history)

#### Require Signed Commits
- ✅ **Require signed commits**
  - Ensures authenticity of commits

### For Team/Organization Repos

If you have a team:

1. **CODEOWNERS file** (`.github/CODEOWNERS`):
   ```
   # Backend code requires review from backend team
   * @your-org/backend-team
   
   # All Go files require review
   *.go @your-org/backend-team
   ```

2. **Required Reviewers**: Set minimum number of approvals

3. **Review Assignments**: Auto-assign reviewers

### Testing Branch Protection

#### Test 1: PR with Passing Tests
```bash
git checkout -b feature/working-code
# Make a valid change
git commit -am "Add new feature"
git push origin feature/working-code
# Create PR → Should allow merge after CI passes
```

#### Test 2: PR with Failing Tests
```bash
git checkout -b feature/broken-code
# Introduce a bug or failing test
git commit -am "Add broken code"
git push origin feature/broken-code
# Create PR → Should BLOCK merge until fixed
```

#### Test 3: PR with Unformatted Code
```bash
git checkout -b feature/bad-formatting
# Add unformatted Go code
git commit -am "Add unformatted code"
git push origin feature/bad-formatting
# Create PR → Should FAIL go fmt check
```

### Workflow Triggers

The CI runs on:
- ✅ **Push to main** (after merge)
- ✅ **Pull Request to main** (before merge)
- ✅ **All backend changes** (since backend and frontend are separate repos)

### Troubleshooting

#### Issue: "Backend CI" doesn't appear in status checks list

**Solution:**
1. Make sure the workflow has run at least once
2. Go to Actions tab → Click on "Backend CI" workflow
3. Wait for it to complete
4. Go back to Settings → Branches → The check should now appear

#### Issue: Can't find the check in the dropdown

**Solution:**
1. Type the exact name: `build-and-test` or `lint`
2. Make sure you've pushed the `.github/workflows/go.yml` file
3. Check Actions tab to see if workflow ran

#### Issue: Workflow fails but I want to merge anyway

**Not Recommended**, but if necessary:
1. As repo admin, you can temporarily disable the rule
2. OR add yourself to bypass list
3. OR fix the code to pass the checks (recommended!)

### CI/CD Pipeline Overview

```
Developer pushes code
        ↓
GitHub detects changes
        ↓
Backend CI workflow triggers
        ↓
    ┌─────────┴─────────┐
    ↓                   ↓
build-and-test       lint
    ↓                   ↓
Both must pass ✅
        ↓
PR can be merged 🎉
        ↓
Code pushed to main
        ↓
Render auto-deploys
```

### Benefits

1. **Prevents Broken Code**: Can't merge if build fails
2. **Code Quality**: Ensures formatting and linting standards
3. **Test Coverage**: All tests must pass
4. **Security**: Catches race conditions and common bugs
5. **Team Confidence**: Everyone knows main branch works

### Cost

- **Free for public repos**
- **Free tier for private repos**: 2,000 CI minutes/month
- Typical run time: ~2-3 minutes per PR
- Can do ~600+ PRs per month on free tier

### Next Steps

1. ✅ Set up branch protection (follow steps above)
2. ✅ Test with a sample PR
3. ✅ Add CODEOWNERS if team repo
4. ✅ Configure Codecov (optional, for coverage reports)
5. ✅ Add status badge to README

### Status Badge for README

Add this to your `README.md`:

```markdown
![Backend CI](https://github.com/YOUR-USERNAME/Gymie-backend/workflows/Backend%20CI/badge.svg)
```

Replace `YOUR-USERNAME` with your GitHub username.

---

**Note:** The backend and frontend are maintained in separate repositories:
- Backend: `Gymie-backend` (this repo)
- Frontend: `Gymie-frontend` (separate repo)

Each repository should have its own CI/CD workflows and branch protection rules.

---

**Need Help?** Check GitHub's official docs:
- [About branch protection rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)
- [Require status checks](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches#require-status-checks-before-merging)
