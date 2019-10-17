POC to attempt ingesting sample points and polygons 
and attempt to match point in poly using Uber's H3 library

#### Prerequisites
* set GOPATH
* have a standalone instance of redis running on default port in localhost

#### Usage instructions
run the commands in the bin folder
* ingest point - `./ingest point`
* ingest polygon - `./ingest polygon`

#### Done Statement
- [x] ingest points
- [x] ingest polygons
- [x] run point-in-poly matches
- [ ] visualize and benchmark results