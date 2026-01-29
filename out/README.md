# backupper

A simple backup tool for backing up files and directories and sending a Telegram message when the backup is done.

## How to use

```bash
Usage of backupper:
  -backup-folders string
    	folders to backup, comma separated
  -log string
    	formats of the logs (text|json) (default "text")
  -repo-pass string
    	password to the repository
  -repo-path string
    	path to the repository
  -telegram-api-key string
    	Telegram API key
  -telegram-chat-id string
    	Telegram chat ID
```


## Example

```bash
backupper -backup-folders /home/user/folder1,/home/user/folder2 -repo-path /home/user/repo -repo-pass password -telegram-api-key api-key -telegram-chat-id chat-id
2024/06/03 17:01:29 INFO starting backup
2024/06/03 17:01:29 INFO backup started start=2024-06-03T17:01:29.331+01:00
2024/06/03 17:01:29 INFO backup finished end=2024-06-03T17:01:29.796+01:00
2024/06/03 17:01:29 INFO message sent
```
