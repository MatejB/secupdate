# secupdate
Execute Ubuntu unattended-upgrade and email log via SMTP.

Application expects certain ENV variables to be set:

```bash

# Email where to send report
SU_TO_EMAIL="" 

# From what email to send it
SU_FROM_NAME=""

# From name
SU_FROM_EMAIL=""

# SMTP auth username
SU_USER=""

# SMTP auth password
SU_PASS=""

# SMTP host
SU_HOST=""

# SMTP port
SU_PORT=""

# Verbose report mode
# Value can be any of the fallowing to turn it on: true, on, verbose, 1, enable, enabled
SU_VERBOSE=""

```

[Sample bash script](setenv_sample.sh) is provided for setting this variables.
Don't forget to [set ENV variables permanently](http://unix.stackexchange.com/questions/117467/how-to-permanently-set-environmental-variables) if program will be run as cron job.
