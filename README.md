# govisn

A 3D network visualization tool written in golang.

# History of the Project
This project is a continuation of the work that started with vrmlNet, which was developed in 1998. vrmlNet created VRML code that could be rendered by VRML Browser Plug-ins.

In the early 2000's VRML technology had sunset. A second generation of the project, titled V15N, was developed beginning in the Fall of 2003. This was a complete re-write of vrmlNet . Unlike vrmlNet, which used the Virtual Reality Modeling Language (VRML), V15N uses the Java 3D API for rendering a 3D visualization of a network.   The application took as input a seed IP Address and used SNMP to walk the routing tables of routers within a network. It then used Java3D to visualize that network in a virtual 3D space.

govisn is the third generation of this 3D network visualization tool. By the Fall of 2018, 3D visualization technology had left Java3D behind. I was looking to continue the 3D network visualization concept, plus learn a new coding language. It was time to again re-write the application. Rather than develop a multi-tiered web-based application, I chose to keep with the original pinciples of the project: Free Open Source and a simple implmentation. The original goal was to build something a network engineer could deploy, without requiring complicated systems engineering or administation tasks. 

# Attributions

govisn is written in the go language and uses go-sqlite3, gosnmp and G3N libraries for their database, SNMP and 3D rendering capabilities. Many thanks to the authors of these libraries for the use of their work. Yasuhiro Matsumoto (a.k.a mattn) and G.J.R. Timmer for go-sqlite3, Sonia Hamilton, sonia@snowfrog.net for gosnmp, and Daniel Salvadori and leonsal for G3N. I stand on the shoulders of giants.

govisn is a single executable application. This supports the orginal principles of this generational series of applications, in that it is usable at no cost and easy to install and operate. The command line executable takes startup options to scan a subnet for routers. It then queries those routers using SNMP to create a sqlite database. The database is used to render the 3D scene containing the router objects and the network links between them.

# Usage of govisn:

  -a    Test opening an ArangoDB database
  
  -co string
        SNMP Community ReadOnly String (default "public")
  -cr
        Create a sample database.
  -de
        Print Debug statements.
  -di string
        Discover a network using seed IP Address
  -f string
        Name of the discovered network database -or-
        Name of the XML input file, if combined with -l option. (default "govisnDiscoveredNet.db")
  -l    Load a database from an XML document.
  -m string
        Scope of discovery. Maximum number of Hops from seed. (Default:10) (default "0")
  -s string
        Scan the CIDR network for SNMP capable routers.
        CIDR format = x.x.x.x/n. ex: 192.168.1.0/24
        Once the network is scanned, the list of found routers
        will be queried and their information added to the database.
  -v    Print the version number.
  -vi
        Visualize the Network.