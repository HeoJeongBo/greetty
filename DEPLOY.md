# Deploying greetty to Homebrew

This is the manual release checklist. greetty ships via a **Homebrew tap**: a
second GitHub repo named `homebrew-greetty` that holds the formula. End users
then run `brew install HeoJeongBo/greetty/greetty`.

You run all of this yourself (the repos and SSH key are yours).

---

## 0. SSH note

Your `~/.ssh/config` uses a host alias:

```
Host github-personal
  HostName github.com
  User git
  IdentityFile ~/.ssh/id_ed25519_personal
```

The remote you gave (`git@github.com:...`) uses the **default** key, not the
personal one. If you want the personal key to be used, set the remote to the
alias instead:

```sh
git remote set-url origin git@github-personal:HeoJeongBo/greetty.git
```

(Both work — use whichever key you intend to push with.)

---

## 1. Push the main repo

From the project root:

```sh
git init
git add .
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:HeoJeongBo/greetty.git   # or git@github-personal:...
git push -u origin main
```

Create the repo on GitHub first as **public** (`HeoJeongBo/greetty`) if it
doesn't exist yet — via the web UI, or with `gh repo create HeoJeongBo/greetty --public --source=. --remote=origin`.

## 2. Tag a release

```sh
git tag v0.1.0
git push origin v0.1.0
```

GitHub automatically serves a source tarball at:

```
https://github.com/HeoJeongBo/greetty/archive/refs/tags/v0.1.0.tar.gz
```

(Optionally also run `gh release create v0.1.0 --generate-notes`.)

## 3. Compute the tarball sha256

```sh
curl -fsSL https://github.com/HeoJeongBo/greetty/archive/refs/tags/v0.1.0.tar.gz \
  | shasum -a 256
```

Copy the hash.

## 4. Create the tap repo

Make a second **public** repo named exactly `homebrew-greetty`
(`HeoJeongBo/homebrew-greetty`). The `homebrew-` prefix is required — that's how
`brew install HeoJeongBo/greetty/greetty` resolves.

```sh
git clone git@github.com:HeoJeongBo/homebrew-greetty.git
cd homebrew-greetty
mkdir -p Formula
# copy the formula from this repo:
cp ../greetty/Formula/greetty.rb Formula/greetty.rb
```

Open `Formula/greetty.rb` and replace `REPLACE_WITH_TARBALL_SHA256` with the
hash from step 3. Then:

```sh
git add Formula/greetty.rb
git commit -m "greetty 0.1.0"
git push
```

## 5. Test the install

```sh
brew install HeoJeongBo/greetty/greetty
greetty init
exec zsh
```

To re-test after edits without re-downloading:
`brew reinstall HeoJeongBo/greetty/greetty`, or audit with
`brew audit --new --formula greetty`.

---

## Releasing a new version later

1. Bump, commit, `git tag v0.2.0 && git push origin v0.2.0`.
2. Recompute sha256 for the new tarball (step 3).
3. In the tap repo, update `url` (tag) and `sha256` in `Formula/greetty.rb`, commit, push.

The `version` is injected at build time from the formula via ldflags
(`-X github.com/HeoJeongBo/greetty/internal/cli.version=#{version}`), so the formula's tag
drives `greetty --version` automatically.
