# **GoVisn**  
Copyright 2020 [Kevin Hayes Parrish](mailto:govisn@mycci.net). All rights reserved.  
Please review the LICENSE file before usinging the application.

A 3D network visualization tool written in golang.

## History of the Project
This project is a continuation of the work that started with Java application **vrmlNet**, which was developed in 1998. vrmlNet created VRML code that could be rendered by VRML Browser Plug-ins.

In the early 2000's VRML technology had sunset. A second generation of the project, titled **V15N**, was developed beginning in the Fall of 2003. This was a complete re-write of vrmlNet . Unlike vrmlNet, which used the Virtual Reality Modeling Language (VRML), V15N used the Java 3D API for rendering a 3D visualization of a network. The application took as input a seed IP Address and used SNMP to walk the routing tables of routers within a network. It then used Java3D to visualize that network in a virtual 3D space.

**GoVisn** is the third generation of this 3D network visualization tool. By the Fall of 2018, 3D visualization technology had left Java3D behind. I was looking to continue the 3D network visualization concept, plus learn a new coding language. It was time to again re-write the application. Rather than develop a multi-tiered web-based application, I chose to keep with the original pinciples of the project:  
1. Free Open Source Software with a simple implmentation.  
1. Build something a network engineer could deploy, without requiring complicated systems engineering or administation tasks.  
  
GoVisn is a single executable application. The command line executable takes startup options to scan a subnet for routers. It then queries those routers using SNMP to create a sqlite3 database. The database is used to render the 3D scene containing the router objects and the network links between them.  
  
Dynamic updates of the 3D model are still under development. This feature will periodically query router interfaces and calculate an approximate link utilization percentage. Depending on the link utilization, the 3D network link will be modified to reflect that utilization.

## Attributions

GoVisn is written in the go language and uses [go-sqlite3](https://github.com/mattn/go-sqlite3), [gosnmp](https://github.com/soniah/gosnmp) and [G3N](https://github.com/g3n/engine) libraries for their database, SNMP and 3D rendering capabilities. Many thanks to the authors of these libraries for the use of their work. Yasuhiro Matsumoto (a.k.a mattn) and G.J.R. Timmer for [go-sqlite3](https://github.com/mattn/go-sqlite3), Sonia Hamilton, sonia@snowfrog.net for [gosnmp](https://github.com/soniah/gosnmp), and Daniel Salvadori and leonsal for [G3N](https://github.com/g3n/engine). I stand on the shoulders of giants.
  
## Caveats  
1. The Apple MacOS implementation of G3N only allows a linewidth of 1. Therefore, on MacOS implementations of GoVisn, the network links will always be a linewidth of 1, regardless of the link utilization percentage.  
2. The technique to disover the Layer 3 network using a single router seed address to is still under development. I need to port the resursion algorithm from the Java V15N app to golang.  

## Usage of GoVisn:
govisn *options*
### Options
>**-a**  
>>Test opening an ArangoDB database  
>> **(DEPRECATED)**  
>
>**-co** *string*  
>> SNMP Community ReadOnly String  
>> *(default: "public")*
>
>**-cr**  
>> Create a sample database.  
>>**(DEPRECATED)**  
>
>**-de**  
>> Print Debug statements.
>
>**-di** *string*  
>> Discover a network using seed IP Address  
>> **(Under development)**
>
>**-f** *string*
>> Name of the discovered network database -or-  
>> Name of the XML input file, if combined with -l option.  
>> (default "govisnDiscoveredNet.db")
>
>**-l**  
>> Load a database from an XML document.  
>> **(DEPRECATED)**  
>
>**-m** *string*
>>Scope of discovery. Maximum number of Hops from seed. 
>>        *(default: "0")*
>
>**-s** *string*  
>> Scan the CIDR network for SNMP capable routers.  
>> CIDR format = x.x.x.x/n. ex: 192.168.1.0/24  
>> Once the network is scanned, the list of found routers
>> will be queried and their information added to the database.
>
>**-v**  
>> Print the GoVision version number.
>  
>**-vi**  
>> Visualize the Network.

### Execution Examples
1. Scan a subnet, create a database, then visualize the Layer 3 network.  
      govisn -s *192.168.1.0/24* -f *test.db* -vi -co *public*
2. Visualize the Layer 3 network, using test.db database and SNMP community public  
      govisn -vi -f test.db -co public  
3. Visualize the Layer 3 network with Debug logging enabled.  
      govisn -vi -f test.db -co public -de