# qcard

qcard is a CLI addressbook application for CardDAV servers written in Go. In
contrast to other tools it does not cache anything. It can fetch multiple
servers / addressbooks in parallel which makes it quite fast.

Its main purpose is displaying addressbook data. Nevertheless it supports basic
creation and editing of entries.

## Features

- Easy search for contacts
- Parallel fetching of multiple addressbooks 
- Easy to use filters
- Create, modify and delete contacts 
- Import VCF files
- Display VCF files
- Easy setup


## Installation / Configuration

- Have Go installed
- make && sudo make install (for MacOS: make darwin)
- copy config-sample.json to ~/.config/qcard/config.json and modify accordingly

### Arch AUR package

- Here: [AUR](https://aur.archlinux.org/packages/qcard)
- Copy config-sample.json from /usr/share/qcard/ to ~/.config/qcard/config.json and modify accordingly

## Configuration

- For additional addressbooks just add a comma and new addressbook credentials in
  curly brackets.


## Usage

Common options:

    qcard -h

### Displaying addressbooks 

This simply displays all contacts from all addressbook:

    qcard

This only shows contacts from addressbook 0:

    qcard -c 0

This shows all appointments from 01.10.2021, 00:00h to 31.10.2021, 23:59:59
(Note: This is in UTC!):

    qcal -s 20211001T000000 -e 20211031T235959

This displays all avaliable calendars with their numbers and colors:

    qcal -l

### Add new appointment

Even though the abillity to create new appointments is limited, it is easy to
create simple appointment types.

This creates an appointment on 01.12.2021 from 15:00h to 17:00h with the
summary of "Tea Time":

    qcal -n "20211201 1500 1700 Tea Time"

This creates a whole day appointment with a yearly recurrence in your second
calendar (first is 0):

    qcal -c 1 -n "20211114 Anne's Birthday" -r y

This creates a multiple day appointment:

    qcal -n "20210801 20210810 Holiday in Thailand"

### Edit an appointment

This shows the next 7 days of appointments from calendar 3 with filenames
("foobarxyz.ics"):

    qcal -c 2 -7 -f 

This edits the selected iCAL object in your $EDITOR (i.e. vim). When you
save-quit the modified object is automatically uploaded:

    qcal -c 2 -edit foobarxyz.ics


## Integrations

### neomutt / other cli mail tools

You can view received appointments in neomutt with qcal! Put this in your
mailcap (usually in .config/neomutt):

    text/calendar; qcal -p; copiousoutput

If you also want to be able to import received appointments directly from
neomutt, put the following two lines in mailcap:

    text/calendar; qcal -c 0 -u %s && notify-send "Appointment created";
    text/calendar; qcal -p; copiousoutput

The first line is only executed if you press Return. The second line just
displays the appointment as above.

### Crontab (or Statusline script, Systemd timer, etc.) 

You can get reminders of your appointments 15 mins in advance with this one
liner:

    [[ $(qcal -cron 15 2>/dev/null) ]] && notify-send "Next Appointment:" "\n $(qcal -cron 15)" || true


## About

Questions? Ideas? File bugs and TODOs through the [issue
tracker](https://todo.sr.ht/~psic4t/qcard) or send an email to
[~psic4t/qcard@todo.sr.ht](mailto:~psic4t/qcard@todo.sr.ht)
