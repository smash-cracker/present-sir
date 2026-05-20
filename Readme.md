### How to use

`./present-linux "your_email" "your_password"`

### How to add crontab

#### in linux
`(crontab -l 2>/dev/null; echo "30 10 * * 1-6 cd $(pwd) && ./present your_email 'your_password' >> $(pwd)/present.log 2>&1") | crontab -`
