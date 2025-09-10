# ocmcp

## Setup

### Fetch OpenConfig YANG Models

```
curl -LO https://github.com/openconfig/public/archive/refs/tags/v5.4.0.tar.gz
tar xvf v5.4.0.tar.gz
mv public-5.4.0 public
rm v5.4.0.tar.gz
```

### Build & Local Client for Testing

```
make
./ocmcp -http_address=":8080" -yang_dirs public/ -alsologtostderr
go run github.com/bstoll/ocmcp/client -query "sid-id"
```

### Register with gemini cli

Update gemini cli config in ~/.gemini/settings.json
```
$ cat ~/.gemini/settings.json 
{
  "hasSeenIdeIntegrationNudge": true,
  "ideMode": true,
  "mcpServers": {
    "ocmcp": {
      "command": "/usr/local/google/home/bstoll/Projects/ocmcp/ocmcp",
      "args": [
        "-yang_dirs",
        "/usr/local/google/home/bstoll/Projects/ocmcp/public"
      ],
      "trust": true
    }
  },
  "theme": "GitHub"
}
```

## Examples

### Easy query

With MCP

```
$ gemini -p "What openconfig paths contain sid-id?"
E0910 22:23:34 gemini_impl: main: 65] Gemini CLI is in Dogfood. Always review commands and output before using them. See http://go/gemini-cli-dogfood for details.
Here are the OpenConfig paths that contain `sid-id`:

*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/adjacency-sids/adjacency-sid/sid-id**
    *   Description: Reference to the segment identifier to be used by the local system.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/prefix-sids/prefix-sid/config/sid-id**
    *   Description: The Segment Identifier to be used when advertising the IGP Prefix SID.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/prefix-sids/prefix-sid/state/sid-id**
    *   Description: The Segment Identifier to be used when advertising the IGP Prefix SID.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/flex-algo-prefix-sids/flex-algo-prefix-sid/config/sid-id**
    *   Description: The Segment Identifier to be used when advertising the IGP Prefix SID for the Flexible Algorithm.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/flex-algo-prefix-sids/flex-algo-prefix-sid/state/sid-id**
    *   Description: The Segment Identifier to be used when advertising the IGP Prefix SID for the Flexible Algorithm.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/prefix-sids/prefix-sid**
    *   Description: An IGP prefix that should have a segment routing IGP-Prefix SID allocated to it. The value of the SID is specified by the SID ID, as an absolute value. If the absolute value falls within the SRGB, the Global flag should be advertised by the system.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/adjacency-sids/adjacency-sid**
    *   Description: An Adjacency SID to be advertised for the specified interface. The Adj-SID's identifier (the SID ID) must be unique, with flags specified indicating the parameters that should be set for the SID. Where a SID value is specified that is allocated from the SRGB, the global flag must be set by the system.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/adjacency-sids/adjacency-sid/state/sid-id**
    *   Description: The value of the Adj-SID to be advertised. Where a static SID identifier is specified, this should be advertised directly by the system. Where the DYNAMIC value is specified, this should be treated as a dynamically allocated value. When the MPLS data plane is in use the dynamic value should not fall within a reserved-label-block.
*   **/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/adjacency-sids/adjacency-sid/config/sid-id**
    *   Description: The value of the Adj-SID to be advertised. Where a static SID identifier is specified, this should be advertised directly by the system. Where the DYNAMIC value is specified, this should be treated as a dynamically allocated value. When the MPLS data plane is in use the dynamic value should not fall within a reserved-label-block.
```

Without MCP, every path listed is hallucinated:

```
$ gemini -p "What openconfig paths contain sid-id?"
E0910 22:18:24 gemini_impl: main: 65] Gemini CLI is in Dogfood. Always review commands and output before using them. See http://go/gemini-cli-dogfood for details.
Based on my knowledge of the OpenConfig YANG models, `sid-id` is typically found in paths related to Segment Routing (SR) and MPLS. Here are some of the common paths where you would find a `sid-id` leaf:

*   **Segment Routing Policies:**
    *   `/network-instances/network-instance/policies/policy/sr-te-policy/candidate-paths/candidate-path/segment-lists/segment-list/next-hops/next-hop/segments/segment/sid-id`

*   **MPLS Label Switched Paths (LSPs):**
    *   `/network-instances/network-instance/mpls/lsps/constrained-path/tunnels/tunnel/p2p-tunnel-attributes/srlgs/srlg/sid-id`

*   **SRLG (Shared Risk Link Group) Configuration:**
    *   `/network-instances/network-instance/mpls/te-global-attributes/srlgs/srlg/static-srlg-members/members/member/sid-id`

*   **Aggregate SID Counters:**
    *   `/network-instances/network-instance/mpls/signaling-protocols/segment-routing/aggregate-sid-counters/aggregate-sid-counter/sid-counter/sid-id`

These paths can vary slightly depending on the specific implementation and the version of the OpenConfig models being used. If you have a specific area of configuration in mind (like SR-TE policies or MPLS), I can provide more targeted paths.
```

### More complex question...

With MCP

````
$ gemini -p "Please find the Openconfig path for the following Arista config snippet. I am using ISIS:
interface Loopback1
ip address 108.177.36.20/32
node-segment ipv4 index 833"
E0910 22:38:02 gemini_impl: main: 65] Gemini CLI is in Dogfood. Always review commands and output before using them. See http://go/gemini-cli-dogfood for details.
Based on the Arista configuration snippet you provided, the corresponding OpenConfig path for the `node-segment ipv4 index 833` under interface `Loopback1` is:

```
/network-instances/network-instance/protocols/protocol/isis/interfaces/interface/levels/level/afi-safi/af/segment-routing/prefix-sids/prefix-sid
```

To be more specific, if you were to set the SID index, the full path would look like this, including the keys for the different lists in the path:

```
/network-instances/network-instance[name=default]/protocols/protocol[identifier=ISIS][name=ISIS]/isis/interfaces/interface[interface-id=Loopback1]/levels/level[level-number=2]/afi-safi/af[afi-name=oc-types:IPV4][safi-name=oc-types:UNICAST]/segment-routing/prefix-sids/prefix-sid[prefix=108.177.36.20/32]/config/sid-id
```

And the value for `sid-id` would be `833`.

Please note that some values in the path are assumed:
*   Network instance name: `default`
*   ISIS instance name: `ISIS`
*   ISIS level: `2`
*   AFI/SAFI: `IPV4 UNICAST`
````

Without MCP
````
$ gemini -p "Please find the Openconfig path for the following Arista config snippet. I am using ISIS:
interface Loopback1
ip address 108.177.36.20/32
node-segment ipv4 index 833"
E0910 22:40:26 gemini_impl: main: 65] Gemini CLI is in Dogfood. Always review commands and output before using them. See http://go/gemini-cli-dogfood for details.
Based on the Arista config snippet provided, here are the corresponding OpenConfig paths for an ISIS configuration.

### IP Address on Loopback Interface

This is the path for the `ip address 108.177.36.20/32` configuration on `Loopback1`.

```json
/interfaces/interface[name=Loopback1]/subinterfaces/subinterface[index=0]/ipv4/addresses/address[ip=108.177.36.20]/config
```

-   **prefix-length**: `32`

### ISIS Node Segment (Prefix SID)

This is the path for the `node-segment ipv4 index 833` configuration, which maps a Prefix SID to the loopback interface in ISIS. Note that the ISIS instance name is assumed to be `DEFAULT`.

```json
/protocols/protocol[identifier=ISIS][name=DEFAULT]/isis/interfaces/interface[interface-id=Loopback1]/segment-routing/prefix-sids/prefix-sid[ip-prefix=108.177.36.20/32]/config
```

-   **sid-id**: `833`
-   **flags**: `["NODE_SID"]` (This flag is typically implicit in the `node-segment` command)
````
