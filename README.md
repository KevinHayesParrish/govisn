# govisn

A network visualization tool written in golang.
This project is a continuation of the work that started with vrmlNet. A second generation version was titled V15N. This was a complete re-write of vrmlNet . Unlike vrmlNet, which used the Virtual Reality Modeling Language (VRML), V15N uses the Java 3D API for rendering a 3D visualization of a network. 

govisn is the third generation of this 3D network visualization tool. It is written in the go language and uses go-sqlite3, gosnmp and g3n libraries to for their database, SNMP and 3D rendering capabilities. Many thanks to the authors of these libraries for the use of their work.

govisn is a single executable application. This supports the orginal principles of this generational series of applications, in that it is usable at no cost and easy to install and use.

Usage of govisn:
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