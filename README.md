ppofit
------

Overview
--------

The Golang profiling tools are a pain to get set up in windows, and they're quite dated compared
to the tools available in IDEs like Visual Studio nowadays.

The over-ambitious goal of this is to start building a cross platform tool for the profiling of
Go programs (and perhaps anything using the pprof format eventually) using Go and web browsers as the
IDE. Ideally though, the core functionality will be well factored and tested such that a CLI
or other forms of interaction can be added to it. 

My initial thinking is to consider using Dart 1.0 for the webapp, but it might just be easier to
use javascript.


Milestones
----------

[1] ☑ Backend code is able to interrogate a sample application using net/http/pprof and obtain a profile
[2] ☑ Profile can be interpreted using http://google-perftools.googlecode.com/svn/trunk/doc/cpuprofile-fileformat.html
[3] ☑ Function pointers can be interpreted
[4] ☐ Callgraph structure can be built from the samples
[5] ☐ Simple web-app running displaying some information
[6] ☐ Bootstap / D3 / AngularJS evaluated and potentially integrated
