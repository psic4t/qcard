# qcard

qcard is a CLI addressbook application for CardDAV servers written in Go. In
contrast to other tools it does not cache anything. It can fetch multiple
servers / addressbooks in parallel what makes it quite fast.

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

This simply displays all contacts from all addressbooks:

    qcard

This only shows contacts from addressbook 0:

    qcard -a 0

This displays all avaliable addressbooks with their numbers and colors:

    qcaŕd -l

This searches for contacts containing "doe" in all addressbooks:
    
    qcard -s doe

The DetailThreshold parameter in the configuration file determines when all contact details are shown for a given numer of search results. For instance, on DetailThreshold = 3 you get all details if 3 or less contacts are found for the searchword "doe".

Here's a list of all attributes:


* **M:** phoneCell
* **P:** phoneHome
* **p:** phoneWork
* **E:** emailHome
* **e:** emailWork
* **A:** addressHome
* **a:** addressWork
* **O:** Organisation
* **B:** Birthday
* **T:** Title
* **R:** Role
* **I:** Nickname
* **n:** Note

### Add new contact

This creates a contact for John Doe with a private mobile phone number and an email address in address book 1:

    qcard -a 1 -n "John Doe M:+49 172 123123 E:jdoe@data.haus"

Just combine the parameters from above like you wish.

### Edit a contact

This shows searches for "doe" in addressbook 2 and prints the corresponding filenames
("fobarxyz.vcf"):

    qcard -a 2 -s doe -f

This edits the selected vCard object in your $EDITOR (i.e. vim). When you
save-quit the modified object is automatically uploaded:

    qcard -c 2 -edit foobarxyz.vcf

## Integrations

### neomutt / other cli mail tools

To use qcard as your addressbook in neomutt, put the following in your neomuttrc:

    set query_command= "qcard -s '%s' -emailonly"
    bind editor <Tab> complete-query
    bind editor ^T complete

### External password command

Instead of putting your password in the config file you can specify an
external command to resolve your password. Put a line like this in your
addressbook config and leave the "Password" field empty:

    "PasswordCmd":"rbw get calendar-provider"

## About

Questions? Ideas? File bugs and TODOs through the [issue
tracker](https://todo.sr.ht/~psic4t/qcard) or send an email to
[~psic4t/qcard@todo.sr.ht](mailto:~psic4t/qcard@todo.sr.ht)
