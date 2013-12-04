pprofit
------

Libraries and utilities for pprof, written in Go.

Overview
--------

The Golang profiling tools are a pain to get set up in Windows, and they're quite dated compared
to the tools available in IDEs like Visual Studio nowadays.

The over-ambitious goal of this is to start building a cross platform toolkit for the profiling of the performance of
Go programs using Go itself and web browsers as the IDE. 
Ideally though, the core functionality will be well factored and tested such that a CLI or other forms of interaction can be added to it. 

Local testing
-------------
go build github.com/Redundancy/pprofit/sampleapp github.com/Redundancy/pprofit/integrationtest
Running integrationtest should execute sampleapp and interrogate it to generate a profile, get other symbol information and then shut it down.

This is currently being developed on Windows, but there shouldn't be anything particularly special about that. I will look into getting a CI loop set up on Travis or some other service. 

Milestones
----------

[1] ☑ Backend code is able to interrogate a sample application using net/http/pprof and obtain a profile
[2] ☑ Profile can be interpreted using http://google-perftools.googlecode.com/svn/trunk/doc/cpuprofile-fileformat.html
[3] ☑ Function pointers can be interpreted
[4] ☐ Callgraph structure can be built from the samples
[5] ☐ Simple web-app running displaying some information
[6] ☐ Bootstap / D3 / AngularJS evaluated and potentially integrated
