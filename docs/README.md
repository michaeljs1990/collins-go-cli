# Generate Man Pages

We write all of the documentation in text files and convert it into man
pages using pandoc. Depending on your distro it's likely easiest to just
download the release from github.

To generate a man pages you can run the following command. Man pages should
never be directly edited.

```
pandoc query.md -f markdown -t man -s -o query.1
```
