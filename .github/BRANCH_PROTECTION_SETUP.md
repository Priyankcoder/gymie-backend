
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
     - In the search box, type: `Backend CI` or `build`
     - ✅ Select **Backend CI / Build Backend** from the dropdown
   
   - ✅ **Require conversation resolution before merging** (optional)
   
   - ✅ **Do not allow bypassing the above settings** (recommended)

4. Click **Create** or **Save changes**

### Step 3: Verify It's Working

1. Create a new branch: `git checkout -b test-branch`
2. Make a change to any backend file
3. Push and create a PR
4. You should see:
   - ✅ "Backend CI / Build Backend" check running
   - 🔒 **Merge button disabled** until check passes

### Visual Guide

```
GitHub Repository Settings
├── Branches
│   └── Branch protection rules
│       └── Rule for: main
│           ├── ✅ Require pull request reviews
│           ├── ✅ Require status checks:
│           │   └── ✅ Backend CI / Build Backend
│           └── ✅ Do not allow bypassing
```

### What the Workflow Checks

The `Backend CI` workflow checks:

✅ **Build Job:**
- Go dependencies download successfully
- All packages build without errors
- API binary compiles successfully

**That's it!** No linting, no formatting - just verifies the code builds.

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

### Testing Branch Protection

#### Test 1: PR with Valid Code
```bash
git checkout -b feature/working-code
# Make a valid change
git commit -am "Add new feature"
git push origin feature/working-code
# Create PR → Should allow merge after build passes
```

#### Test 2: PR with Broken Code
```bash
git checkout -b feature/broken-code
# Introduce a syntax error or import issue
git commit -am "Add broken code"
git push origin feature/broken-code
# Create PR → Should BLOCK merge until fixed
```

### Workflow Triggers

The CI runs on:
- ✅ **Push to main** (after merge)
- ✅ **Pull Request to main** (before merge)
- ✅ **All changes** (runs on every commit)

### Troubleshooting

#### Issue: "Backend CI" doesn't appear in status checks list

**Solution:**
1. Make sure the workflow has run at least once
2. Go to Actions tab → Click on "Backend CI" workflow
3. Wait for it to complete
4. Go back to Settings → Branches → The check should now appear

#### Issue: Can't find the check in the dropdown

**Solution:**
1. Type: `Build Backend`
2. Make sure you've pushed the `.github/workflows/go.yml` file
3. Check Actions tab to see if workflow ran

#### Issue: Build fails but I want to merge anyway

**Not Recommended**, but if necessary:
1. As repo admin, you can temporarily disable the rule
2. OR add yourself to bypass list
3. OR fix the code to build (recommended!)

### CI/CD Pipeline Overview

```
Developer pushes code
        ↓
GitHub detects changes
        ↓
Backend CI workflow triggers
        ↓
Build Backend job runs
        ↓
✅ Build succeeds
        ↓
PR can be merged 🎉
        ↓
Code pushed to main
        ↓
Render auto-deploys
```

### Benefits

1. **Prevents Broken Code**: Can't merge if build fails
2. **Fast Feedback**: Quick build-only check (~1-2 minutes)
3. **Simple**: No complex linting rules to configure
4. **Team Confidence**: Main branch always builds

### Cost

- **Free for public repos**
- **Free tier for private repos**: 2,000 CI minutes/month
- Typical run time: ~1-2 minutes per PR
- Can do ~1000+ PRs per month on free tier

### Next Steps

1. ✅ Set up branch protection (follow steps above)
2. ✅ Test with a sample PR
3. ✅ Add CODEOWNERS if team repo (optional)
4. ✅ Add status badge to README (optional)

### Status Badge for README

Add this to your `README.md`:

```markdown
![Backend CI](https://github.com/YOUR-USERNAME/Gymie-backend/workflows/Backend%20CI/badge.svg)
```

Replace `YOUR-USERNAME` with your GitHub username.

---

**Note:** This workflow only checks if the code builds. It does NOT check:
- Code formatting
- Linting/code quality
- Tests
- Code coverage

If you want these checks later, they can be added back to the workflow.

---

**Need Help?** Check GitHub's official docs:
- [About branch protection rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)
- [Require status checks](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches#require-status-checks-before-merging)
