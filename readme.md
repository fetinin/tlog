# ⏱ Tlog - time saving utility
Simple CLI that helps log time on JIRA issues.

## Usage
```bash
log 1h SCENTRE-5912     # log 1 hour into SCENTRE-5912 for today
log 1 5814              # log 1 hour into {{DefaultProject}}-5814
log 30m review          # log 30 minutes into task aliased "review"
log 1h review yesterday # log 1 hour yesterday
log 1h review monday    # log 1 hour review for monday for current week
log 1h review mon       # log 1 hour review for monday for current week
log 1h review 22        # log 1 hour review for 22nd if current month
```

## Install

### MacOS
```bash
brew install what-If-I/tap/tlog
```

### Linux
```bash
curl -sfL https://raw.githubusercontent.com/What-If-I/tlog/master/install.sh | sh
```

### From source
```bash
go install github.com/What-If-I/tlog@latest 
```
_go >= 1.18 required_

## Configuration
Upon first run, utility will create config file called `.time_logger_conf.toml` at your home directory. You can edit config to set DefaultProject and add new issues aliases.

Config example:
```bash
cat ~/.time_logger_conf.toml
```
```toml
JiraURL = "https://company.jira.ru"
JiraLogin = "user.name"
JiraPassword = "password"
DefaultProject = "SCENTRE" # if you only specify JIRA issue number, this project will be used

[ TaskAliases ]
meeting = "INT-18" # aliases "meeting" to INT-18
review = "INT-24"
```

### Things to do
- [ ] Add `config show`, `config set-alias`, `config set-project` commands
- [ ] Allow to log multiple days at once like `tlog 1h review monday-friday`
- [x] Automate releases with https://goreleaser.com/quick-start/
