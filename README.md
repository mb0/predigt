Predigten
=========

This project is to scan and publish the collected Lutheran sermons of my father Christian G.
Schnabel in German, that were lovingly prepared and copied by Erika Schöneich.

The sermons is archived in the Landeskirche Hannover and in the private family collection.
Originally written in fountain pen by my father and then recorded on tape for the people in elderly
homes and copied by type writer by Erika Schöneich. I hold both my father and Frau Schöneich dear to
my heart.

Scanning
--------

A used Brother MFC-7440N is hopefully on its way to scan the machine written sermons.
The tesseract-ocr project can be used to digitize the scanned documents.

A small program helps with automating the scanning and primary editing process.

Publishing
----------

The sermons are published as web site, without cookies and tracking, but a minimalistic view count.
Each entry should be published sequentially according to the church year.

There are about 51 sundays in a year each with their own name and theme according to different
churches. We use the https://de.wikipedia.org/wiki/Perikopenordnung to match the Lutheran context.
There are iCal files at: https://www.kirchenjahr-evangelisch.de/ical-kalender-download.php

Investigation can be made to calculate the liturgical calender using Meeus/Jones/Butcher algorithm:
https://en.wikipedia.org/wiki/Date_of_Easter#Anonymous_Gregorian_algorithm

We might want to use a dictionary with biblical names, a bible reference parser for auto linking
to bible citations, and maybe some highlight correction feature for readers.

ERF media eV https://bibleserver.com can be used to hyper-link citations and maybe as dictionary.

Since the bible translations are usually in the public domain, a SQL dump can be used.
A Luther bible translation from 1912 can be found at:
	https://ebible.org/find/details.php?id=deu1912
	https://www.biblesupersearch.com/bible-downloads/ 

There is online bible lexicon https://www.bibelwissenschaft.de/ressourcen/wibilex
we could add automatic hyper-links to keywords. Maybe they can publish or send us a csv file with
keywords and links.

Data model
----------

sundays
 * id
 * name
 * body

sermon
 * id
 * sunday - reference to the church year
 * org    - original publishing date
 * pub    - website publishing date
 * rev    - last revision
 * title? - maybe a title?
 * body   - sermon text
 * note   - editorial note for web publishing to reference the historic context

comments?
 * id
 * parent ref
 * sermon ref
 * pub
 * from
 * body

likes?
 * id
 * pub
 * sermon ref

 Vim
 ---

 :setlocal spell spellfile=/home/mb0/work/predigt/spell.utf8.add spelllang=de_de
