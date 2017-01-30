# temperature-app
Example temperature measuring app that can run on Raspberry-PI

### Build the application:
Building for amd64:
```
make
```

Building for arm-6:
Make sure that `xgo` is installed on the system, then run:
```
make arm
```

### Running the application:
```
./temperature-app --remote-endpoint=http://my.example.endpoint.com --temperature-threshold=33.5 --temperature-units=Celsius
```