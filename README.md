# gRPC-programming
A subset of *Assignment 4* of the *Distributed Systems* course in IIIT-Hyderabad (Monsoon '24)

## K Nearest Neighbours 

*More details available in the [source directory](./src/K%20Nearest%20Neighbours/).* 

![vary n](./docs/diags/knn%20vary%20n.png)

The above was run with *k=100*. 

Implemented with two separate clients, one to send data to be stored in the server, one to query each server for the $k$ nearest neighbours for their local data points. 

**Sequential Time Complexity** consists of a simple heap implementation, such that for *n* data-points, to find the *k* nearest neighbours, we would need a time complexity of:
$$
O(n \cdot log(k))
$$

**Parallel Time Complexity** consists of querying for the local nearest neighbours of each server, and merging them together on the client side efficiently. Overall, we require a time complexity of:
$$
O(\frac{n}{s}\cdot log(k) + k\cdot log(k) \cdot log(s))
$$

To see run-times and qualitative and quantitative analysis, justified with plots for varying *n*, varying *s* and so on, see [the report](./src/K%20Nearest%20Neighbours/report.pdf). 

## My Uber 

*More details available in the [source directory](./src/MyUber/).*

An over-engineered rider-driver system, consisting of:
- **load balancing** - random picking, pick-first, and round-robin policies available 
- **mutual TLS** - clients (riders and drivers) must provide a valid certificate that is issues by a trusted CA. All communication between client and server is encrypted using SSL / TLS. 
- **interceptors** - for authentication, logging and other metadata. 
- **timeouts** - the driver client must accept or reject a ride within a specified timeout period. 

The servers function completely independently. For a more detailed analysis of the implementation and functionality, see [the report](./src/MyUber/report.pdf). 
