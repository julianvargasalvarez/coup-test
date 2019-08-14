COUP Test

To run the project:
```bash
sudo docker-compose up
```

This will run the application in the port `8080`

In order to retrieve some scooters

```bash
time curl -XGET -H "Content-Type: application/json" http://localhost:8080/scooters/?max=2
```


The implements an MVC pattern without an explicit Model since the data is retrieved from a URL
To retrieve as many scooters as possible, the application uses a `go routine` to retrieve the
scooter for a given ID in background and if the scooter matches the filter criteria it will
be added to a channel.

There is also another `go routine` that closes the channel after 1 second of processing, thus
if the maximum number of scooters has not been reached yet and 1 second has passed since the
begenning of the process, the application will return the scotters collected so far.

ISSUES

- Since the routines run concurrently, most of the times the channel receives more scooters
  than the `max` allowed, we should find a way to avoid this.
- The iterator that increments the IDs repeates the process up to 2 times `max`, this
  should be changed to make it iterate without upper limit and thus retrieve as many
  scooters as possible within the time bounds
