
Naïve approach to write a program that prints itself.
Verifying:

```
go run self.go > /tmp/copyself.txt && diff /tmp/copyself.txt self.go && rm /tmp/copyself.txt
```
