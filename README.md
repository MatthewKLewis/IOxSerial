# IOx Application that reads from Serial Port
Simple Python and Go applications that scan from a COM port, formats the data, and posts to an API. This applications are intended to be containerized in a Docker image, and converted via ioxclient to a Cisco AP runnable program.

The activation json connects the linux container running the program from the inside-out - the package yaml from the outside-in