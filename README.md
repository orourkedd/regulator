# Regulator

A simple utility for limited concurrency in go, with error handling.  If an error occurrs, `regulator` will not execute any more functions and will return the error.

## Example

Execute 10 functions, 2 at a time.

```go
r := regulator.NewRegulator(2) // concurrency of 2

for i := 0; i < 10; i++ {
    r.Execute(func() error {
        // do something here and optionally return an error
        return nil
    })
}

err := r.Wait()
if err != nil {
    panic(err) // do real error handling
}
```
